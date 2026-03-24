package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AcidicAcidity/t-mess/internal/crypto"
	"github.com/AcidicAcidity/t-mess/internal/messages"
	"github.com/AcidicAcidity/t-mess/internal/p2p"
	"github.com/AcidicAcidity/t-mess/internal/storage"
	"github.com/AcidicAcidity/t-mess/internal/tui/components"
)

type AppState int

const (
	StateNickInput AppState = iota
	StateMenu
	StateSettings
	StateChat
)

type App struct {
	splash *SplashScreen
	ready  bool
	width  int
	height int
	state  AppState

	identity *crypto.Identity
	db       *storage.Database
	p2pNode  *p2p.Node
	theme    Theme

	// Компоненты
	nickInput   *components.NickInput
	menu        *components.Menu
	settings    *components.Settings
	chatList    *components.ChatList
	messageView *components.MessageView
	inputField  *components.InputField
	rightPanel  *components.RightPanel

	// Данные
	currentChat *messages.Chat
	chats       []*messages.Chat
	messages    []*messages.Message
	nick        string
}

func NewApp(identity *crypto.Identity, db *storage.Database, p2pNode *p2p.Node) *App {
	themeColor := lipgloss.Color(Matrix.Primary)

	// Загружаем никнейм из БД (если есть)
	nick := loadNick(db)
	if nick == "" {
		nick = "Anonymous"
	}

	return &App{
		splash:      NewSplashScreen(),
		identity:    identity,
		db:          db,
		p2pNode:     p2pNode,
		theme:       Matrix,
		state:       StateNickInput,
		nick:        nick,
		nickInput:   components.NewNickInput(themeColor),
		menu:        components.NewMenu(themeColor, nick),
		settings:    components.NewSettings(themeColor, nick, "matrix", "English"),
		chatList:    components.NewChatList(themeColor),
		messageView: components.NewMessageView(themeColor),
		inputField:  components.NewInputField(),
		rightPanel:  components.NewRightPanel(themeColor),
	}
}

// loadNick загружает никнейм из БД
func loadNick(db *storage.Database) string {
	var nick string
	err := db.DB().QueryRow("SELECT value FROM settings WHERE key = 'nickname'").Scan(&nick)
	if err != nil {
		return ""
	}
	return nick
}

// saveNick сохраняет никнейм в БД
func saveNick(db *storage.Database, nick string) {
	db.DB().Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('nickname', ?)", nick)
}

// createSelfChat создаёт чат с самим собой
func (a *App) createSelfChat() {
	chat, err := a.db.CreateSelfChat()
	if err != nil {
		a.rightPanel.SetConnectionStatus(fmt.Sprintf("⚠️ Failed to create chat: %v", err))
		return
	}

	// Обновляем список чатов
	a.chats = append(a.chats, chat)
	a.chatList.SetChats(a.chats)

	// Переключаемся на новый чат
	a.switchToChat(chat)
	a.rightPanel.SetConnectionStatus("✨ Personal Notes created. Write anything!")
}

func (a *App) Init() tea.Cmd {
	return a.splash.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.handleResize(msg.Width, msg.Height)

	case components.ExitRequestMsg:
		return a, tea.Quit

	case components.MenuSelectedMsg:
		switch msg.Choice {
		case components.MenuConnectLocal:
			a.state = StateChat
			a.loadChats()
			a.loadMessages("notes")
			a.rightPanel.SetConnectionStatus("● LOCAL MODE")
			a.inputField.SetOnSend(func(text string) {
				a.sendMessage(text)
			})
			return a, nil

		case components.MenuConnectGlobal:
			a.state = StateChat
			a.loadChats()
			a.loadMessages("notes")
			a.rightPanel.SetConnectionStatus("● GLOBAL MODE (P2P active)")
			a.inputField.SetOnSend(func(text string) {
				a.sendMessage(text)
			})
			return a, nil

		case components.MenuSettings:
			a.state = StateSettings
			return a, nil
		}

	case tea.KeyMsg:
		if a.state == StateChat && msg.String() == "ctrl+q" {
			a.rightPanel.SetConnectionStatus("⚠️ Press Ctrl+C again to exit, any other key to cancel")
			return a, nil
		}
		if a.state == StateChat && msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	if !a.ready {
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		if a.splash.IsDone() {
			a.ready = true
			a.rightPanel.SetPeerID(a.identity.PeerID)
			// Если ник не задан, показываем ввод, иначе сразу меню
			if a.nick == "Anonymous" {
				a.state = StateNickInput
			} else {
				a.state = StateMenu
			}
		}
		return a, cmd
	}

	// Обработка в зависимости от состояния
	switch a.state {
	case StateNickInput:
		newInput, cmd := a.nickInput.Update(msg)
		a.nickInput = newInput

		// Устанавливаем колбек сохранения
		if a.nickInput != nil {
			a.nickInput.SetOnConfirm(func(nick string) {
				a.nick = nick
				saveNick(a.db, nick)
				a.state = StateMenu
				a.menu.SetNick(nick)
			})
		}

		if a.nickInput.IsDone() {
			a.state = StateMenu
		}
		return a, cmd

	case StateMenu:
		newMenu, cmd := a.menu.Update(msg)
		a.menu = newMenu
		return a, cmd

	case StateSettings:
		if a.settings != nil {
			a.settings.SetOnSave(func(nick, themeName, lang string) {
				a.nick = nick
				saveNick(a.db, nick)
				// TODO: применить тему
				// TODO: применить язык
				a.state = StateMenu
				a.menu.SetNick(nick)
			})
			a.settings.SetOnBack(func() {
				a.state = StateMenu
			})
		}

		newSettings, cmd := a.settings.Update(msg)
		a.settings = newSettings
		return a, cmd

	case StateChat:
		// Обработка глобальных хоткеев в чате
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "ctrl+m":
				a.state = StateMenu
				return a, nil
			case "ctrl+s":
				a.state = StateSettings
				return a, nil
			case "ctrl+q":
				return a, tea.Quit
			case "ctrl+n":
				a.createSelfChat()
				return a, nil
			case "ctrl+h": // предыдущий чат (вместо ctrl+[)
				a.prevChat()
				return a, nil
			case "ctrl+l": // следующий чат (вместо ctrl+])
				a.nextChat()
				return a, nil
			case "ctrl+1", "ctrl+2", "ctrl+3", "ctrl+4", "ctrl+5", "ctrl+6", "ctrl+7", "ctrl+8", "ctrl+9":
				idx := int(keyMsg.String()[5] - '1')
				if idx >= 0 && idx < len(a.chats) {
					a.switchToChat(a.chats[idx])
					return a, nil
				}
			}
		}

		// Обновление компонентов чата
		var cmds []tea.Cmd

		newChatList, cmd1 := a.chatList.Update(msg)
		a.chatList = newChatList
		cmds = append(cmds, cmd1)

		newMsgView, cmd2 := a.messageView.Update(msg)
		a.messageView = newMsgView
		cmds = append(cmds, cmd2)

		newInput, cmd3 := a.inputField.Update(msg)
		a.inputField = newInput
		cmds = append(cmds, cmd3)

		newRightPanel, cmd4 := a.rightPanel.Update(msg)
		a.rightPanel = newRightPanel
		cmds = append(cmds, cmd4)

		return a, tea.Batch(cmds...)
	}

	return a, nil
}

func (a *App) View() string {
	if !a.ready {
		return a.splash.View()
	}

	switch a.state {
	case StateNickInput:
		return a.nickInput.View()
	case StateMenu:
		return a.menu.View()
	case StateSettings:
		return a.settings.View()
	case StateChat:
		return a.renderChat()
	}

	return ""
}

func (a *App) renderChat() string {
	leftWidth := 30
	rightWidth := 28
	chatWidth := a.width - leftWidth - rightWidth - 6

	if chatWidth < 40 {
		chatWidth = 40
	}

	chatHeight := a.height - 8 - 3

	a.chatList.SetSize(leftWidth, a.height-3-3)
	a.messageView.SetSize(chatWidth, chatHeight)
	a.rightPanel.SetSize(rightWidth, a.height-3-3)
	a.inputField.SetSize(a.width)

	// Стили заголовков с фоном для контраста
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		Bold(true).
		Background(lipgloss.Color("#0a0a0a")).
		Padding(0, 1)

	leftHeader := headerStyle.Width(leftWidth).Render("📋 CHATS")
	centerHeader := headerStyle.Width(chatWidth).Render("💬 MESSAGES")
	rightHeader := headerStyle.Width(rightWidth).Render("ℹ️ INFO")

	// Левая панель с яркой зелёной границей
	leftBorder := lipgloss.NewStyle().
		BorderRight(true).
		BorderRightForeground(lipgloss.Color("#00ff00")).
		BorderStyle(lipgloss.ThickBorder()).
		PaddingRight(1).
		Height(a.height - 3 - 3)

	leftPanel := leftBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left, leftHeader, a.chatList.View()),
	)

	// Правая панель с яркой зелёной границей
	rightBorder := lipgloss.NewStyle().
		BorderLeft(true).
		BorderLeftForeground(lipgloss.Color("#00ff00")).
		BorderStyle(lipgloss.ThickBorder()).
		PaddingLeft(1).
		Height(a.height - 3 - 3)

	rightPanel := rightBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left, rightHeader, a.rightPanel.View()),
	)

	// Центральная панель
	centerPanel := lipgloss.JoinVertical(lipgloss.Left, centerHeader, a.messageView.View())

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		centerPanel,
		rightPanel,
	)

	topBar := a.renderTopBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		mainContent,
		a.inputField.View(),
	)
}

func (a *App) handleResize(width, height int) {
	a.width = width
	a.height = height
	a.nickInput.SetSize(width, height)
	a.menu.SetSize(width, height)
	a.settings.SetSize(width, height)
}

func (a *App) loadChats() {
	chats, err := a.db.GetChats()
	if err != nil {
		return
	}
	a.chats = chats
	a.chatList.SetChats(chats)

	a.chatList.SetOnSelect(func(chat *messages.Chat) {
		a.switchToChat(chat)
	})

	a.chatList.SetOnCreate(func() {
		a.createSelfChat()
	})

	// Если есть чаты, выбираем первый
	if len(chats) > 0 {
		a.currentChat = chats[0]
		a.loadMessages(chats[0].ID)
		a.rightPanel.SetChat(chats[0])
		a.chatList.Select(0) // выделяем первый чат в списке
	}

	a.rightPanel.SetUserNick(a.nick)

	if a.identity != nil {
		a.rightPanel.SetPeerID(a.identity.PeerID)
	}
}
func (a *App) loadMessages(chatID string) {
	msgs, err := a.db.GetMessages(chatID, 100)
	if err != nil {
		return
	}
	a.messages = msgs
	a.messageView.SetMessages(msgs)
}

func (a *App) sendMessage(text string) {
	if a.currentChat == nil {
		return
	}

	msgID := fmt.Sprintf("%d", time.Now().UnixNano())

	msg := &messages.Message{
		ID:        msgID,
		ChatID:    a.currentChat.ID,
		SenderID:  a.nick,
		Text:      text,
		Timestamp: time.Now(),
		IsOwn:     true,
		Status:    "sent",
	}

	if err := a.db.SaveMessage(msg); err != nil {
		a.rightPanel.SetConnectionStatus(fmt.Sprintf("⚠️ Save error: %v", err))
		return
	}

	a.messageView.AddMessage(msg)

	a.db.DB().Exec(
		"UPDATE chats SET last_message = ?, last_time = ? WHERE id = ?",
		text, time.Now().Unix(), a.currentChat.ID,
	)

	for _, chat := range a.chats {
		if chat.ID == a.currentChat.ID {
			chat.LastMessage = text
			chat.LastTime = time.Now()
			break
		}
	}

	a.chatList.SetChats(a.chats)
}

func (a *App) renderTopBar() string {
	title := "T-MESS"
	if a.currentChat != nil {
		title = a.currentChat.Name
	}

	// Счётчик чатов
	chatCounter := ""
	if len(a.chats) > 0 && a.currentChat != nil {
		for i, chat := range a.chats {
			if chat.ID == a.currentChat.ID {
				chatCounter = fmt.Sprintf(" [%d/%d]", i+1, len(a.chats))
				break
			}
		}
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(a.theme.Primary).
		Bold(true).
		MarginLeft(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		MarginRight(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginRight(2)

	left := lipgloss.JoinHorizontal(
		lipgloss.Left,
		titleStyle.Render(title+chatCounter),
	)

	right := lipgloss.JoinHorizontal(
		lipgloss.Right,
		hintStyle.Render("Ctrl+N:New | Ctrl+[/]:Chat | Ctrl+1-9:Switch | Ctrl+M:Menu | Ctrl+S:Settings | Ctrl+Q:Quit"),
		statusStyle.Render("● ONLINE"),
	)

	return lipgloss.NewStyle().
		Width(a.width).
		Padding(0, 1).
		BorderBottom(true).
		BorderBottomForeground(lipgloss.Color("#00ff00")).
		BorderStyle(lipgloss.ThickBorder()).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, left, right))
}

// nextChat переключается на следующий чат
func (a *App) nextChat() {
	if len(a.chats) == 0 {
		return
	}

	currentIndex := -1
	for i, chat := range a.chats {
		if chat.ID == a.currentChat.ID {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(a.chats)
	a.switchToChat(a.chats[nextIndex])
}

// prevChat переключается на предыдущий чат
func (a *App) prevChat() {
	if len(a.chats) == 0 {
		return
	}

	currentIndex := -1
	for i, chat := range a.chats {
		if chat.ID == a.currentChat.ID {
			currentIndex = i
			break
		}
	}

	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(a.chats) - 1
	}
	a.switchToChat(a.chats[prevIndex])
}

// switchToChat переключает на указанный чат
func (a *App) switchToChat(chat *messages.Chat) {
	a.currentChat = chat
	a.loadMessages(chat.ID)
	a.rightPanel.SetChat(chat)

	// Обновляем выделение в списке
	a.chatList.SetChats(a.chats)
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
