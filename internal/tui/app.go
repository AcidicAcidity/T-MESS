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

type App struct {
	splash *SplashScreen
	ready  bool
	width  int
	height int

	identity *crypto.Identity
	db       *storage.Database
	p2pNode  *p2p.Node
	theme    Theme

	// Компоненты
	chatList    *components.ChatList
	messageView *components.MessageView
	inputField  *components.InputField
	rightPanel  *components.RightPanel

	// Данные
	currentChat *messages.Chat
	chats       []*messages.Chat
	messages    []*messages.Message
}

func NewApp(identity *crypto.Identity, db *storage.Database, p2pNode *p2p.Node) *App {
	themeColor := lipgloss.Color(Matrix.Primary)

	return &App{
		splash:      NewSplashScreen(),
		identity:    identity,
		db:          db,
		p2pNode:     p2pNode,
		theme:       Matrix,
		chatList:    components.NewChatList(themeColor),
		messageView: components.NewMessageView(themeColor),
		inputField:  components.NewInputField(),
		rightPanel:  components.NewRightPanel(themeColor),
	}
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

	case tea.KeyMsg:
		if !a.ready {
			return a, nil
		}
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	if !a.ready {
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		if a.splash.IsDone() {
			a.ready = true
			a.loadChats()
			a.loadMessages("notes")
			a.rightPanel.SetConnectionStatus("● LOCAL ONLY (mock)")

			// Устанавливаем коллбек для отправки
			a.inputField.SetOnSend(func(text string) {
				a.sendMessage(text)
			})
		}
		return a, cmd
	}

	// Обновляем компоненты
	var cmd tea.Cmd

	newChatList, cmd1 := a.chatList.Update(msg)
	a.chatList = newChatList
	cmd = tea.Batch(cmd, cmd1)

	newMsgView, cmd2 := a.messageView.Update(msg)
	a.messageView = newMsgView
	cmd = tea.Batch(cmd, cmd2)

	newInput, cmd3 := a.inputField.Update(msg)
	a.inputField = newInput
	cmd = tea.Batch(cmd, cmd3)

	newRightPanel, cmd4 := a.rightPanel.Update(msg)
	a.rightPanel = newRightPanel
	cmd = tea.Batch(cmd, cmd4)

	return a, cmd
}

func (a *App) sendMessage(text string) {
	if a.currentChat == nil {
		return
	}

	// Генерируем ID сообщения
	msgID := fmt.Sprintf("%d", time.Now().UnixNano())

	msg := &messages.Message{
		ID:        msgID,
		ChatID:    a.currentChat.ID,
		SenderID:  a.identity.PeerID,
		Text:      text,
		Timestamp: time.Now(),
		IsOwn:     true,
		Status:    "sent",
	}

	// Сохраняем в БД
	if err := a.db.SaveMessage(msg); err != nil {
		a.rightPanel.SetConnectionStatus(fmt.Sprintf("⚠️ Save error: %v", err))
		return
	}

	// Добавляем в UI (мгновенно)
	a.messageView.AddMessage(msg)

	// Обновляем последнее сообщение в чате в БД
	_, err := a.db.DB().Exec(
		"UPDATE chats SET last_message = ?, last_time = ? WHERE id = ?",
		text, time.Now().Unix(), a.currentChat.ID,
	)
	if err != nil {
		a.rightPanel.SetConnectionStatus("⚠️ Failed to update chat")
	}

	// Обновляем локальный чат в списке
	for _, chat := range a.chats {
		if chat.ID == a.currentChat.ID {
			chat.LastMessage = text
			chat.LastTime = time.Now()
			break
		}
	}

	// Перерисовываем список чатов (чтобы обновилось превью)
	a.chatList.SetChats(a.chats)

	// TODO: отправка через P2P (пока только локально)
	a.rightPanel.SetConnectionStatus("● LOCAL MODE (P2P coming soon)")
}

func (a *App) View() string {
	if !a.ready {
		return a.splash.View()
	}

	// Ширины панелей
	leftWidth := 28
	rightWidth := 24
	chatWidth := a.width - leftWidth - rightWidth - 4 // -4 на разделители
	chatHeight := a.height - 8

	a.chatList.SetSize(leftWidth, a.height-3)
	a.messageView.SetSize(chatWidth, chatHeight)
	a.rightPanel.SetSize(rightWidth, a.height-3)
	a.inputField.SetSize(a.width)

	// Стили для разделителей
	leftBorder := lipgloss.NewStyle().
		BorderRight(true).
		BorderRightForeground(lipgloss.Color("#00aa00")).
		PaddingRight(1)

	rightBorder := lipgloss.NewStyle().
		BorderLeft(true).
		BorderLeftForeground(lipgloss.Color("#00aa00")).
		PaddingLeft(1)

	// Левая панель с рамкой
	leftPanel := leftBorder.Render(a.chatList.View())

	// Центральная панель (без рамки)
	centerPanel := a.messageView.View()

	// Правая панель с рамкой
	rightPanel := rightBorder.Render(a.rightPanel.View())

	// Основное содержимое
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		centerPanel,
		rightPanel,
	)

	// Верхний бар
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
	// Размеры компонентов будут установлены в View
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}

func (a *App) loadChats() {
	chats, err := a.db.GetChats()
	if err != nil {
		return
	}
	a.chats = chats
	a.chatList.SetChats(chats)

	// Устанавливаем коллбек выбора чата
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
	a.messageView.SetMessages(msgs) // уже в правильном порядке (старые → новые)
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

	left := lipgloss.JoinHorizontal(
		lipgloss.Left,
		titleStyle.Render(title),
	)

	right := statusStyle.Render("● ONLINE")

	return lipgloss.NewStyle().
		Width(a.width).
		Padding(0, 1).
		BorderBottom(true).
		BorderBottomForeground(a.theme.Border).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, left, right))
}
