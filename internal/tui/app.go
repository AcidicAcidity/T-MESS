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

func loadNick(db *storage.Database) string {
	var nick string
	err := db.DB().QueryRow("SELECT value FROM settings WHERE key = 'nickname'").Scan(&nick)
	if err != nil {
		return "Anonymous"
	}
	if nick == "" {
		return "Anonymous"
	}
	return nick
}

func saveNick(db *storage.Database, nick string) {
	db.DB().Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('nickname', ?)", nick)
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
			return a, nil
		case components.MenuConnectGlobal:
			a.state = StateChat
			a.loadChats()
			a.loadMessages("notes")
			a.rightPanel.SetConnectionStatus("● GLOBAL MODE (P2P active)")
			return a, nil
		case components.MenuSettings:
			a.state = StateSettings
			return a, nil
		}

	case tea.KeyMsg:
		if a.state == StateChat && msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	if !a.ready {
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		if a.splash.IsDone() {
			a.ready = true
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
	leftWidth := 28
	rightWidth := 24
	chatWidth := a.width - leftWidth - rightWidth - 4
	chatHeight := a.height - 8

	a.chatList.SetSize(leftWidth, a.height-3)
	a.messageView.SetSize(chatWidth, chatHeight)
	a.rightPanel.SetSize(rightWidth, a.height-3)
	a.inputField.SetSize(a.width)

	leftBorder := lipgloss.NewStyle().
		BorderRight(true).
		BorderRightForeground(lipgloss.Color("#00aa00")).
		PaddingRight(1)

	rightBorder := lipgloss.NewStyle().
		BorderLeft(true).
		BorderLeftForeground(lipgloss.Color("#00aa00")).
		PaddingLeft(1)

	leftPanel := leftBorder.Render(a.chatList.View())
	centerPanel := a.messageView.View()
	rightPanel := rightBorder.Render(a.rightPanel.View())

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
		a.currentChat = chat
		a.loadMessages(chat.ID)
		a.rightPanel.SetChat(chat)
	})
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

	titleStyle := lipgloss.NewStyle().
		Foreground(a.theme.Primary).
		Bold(true).
		MarginLeft(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		MarginRight(1)

	left := lipgloss.JoinHorizontal(lipgloss.Left, titleStyle.Render(title))
	right := statusStyle.Render("● ONLINE")

	return lipgloss.NewStyle().
		Width(a.width).
		Padding(0, 1).
		BorderBottom(true).
		BorderBottomForeground(a.theme.Border).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, left, right))
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
