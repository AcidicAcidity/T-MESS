package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AcidicAcidity/t-mess/internal/messages"
)

type ChatItem struct {
	Chat *messages.Chat
}

func (c ChatItem) Title() string {
	unread := ""
	if c.Chat.UnreadCount > 0 {
		unread = fmt.Sprintf(" [%d]", c.Chat.UnreadCount)
	}

	avatar := "💬"
	if c.Chat.Type == "local" {
		avatar = "📝"
	} else if c.Chat.Type == "direct" {
		avatar = "👤"
	} else if c.Chat.Avatar != "" {
		avatar = c.Chat.Avatar
	}

	return fmt.Sprintf("%s %s%s", avatar, c.Chat.Name, unread)
}

func (c ChatItem) Description() string {
	if c.Chat.LastMessage != "" {
		return truncateString(c.Chat.LastMessage, 40)
	}
	return "No messages"
}

func (c ChatItem) FilterValue() string {
	return c.Chat.Name
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

type ChatList struct {
	list     list.Model
	width    int
	height   int
	onSelect func(*messages.Chat)
}

func NewChatList(theme lipgloss.Color) *ChatList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#ffffff")).
		PaddingLeft(2)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(theme).
		BorderLeft(false).
		PaddingLeft(2).
		Background(lipgloss.Color("#222222"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(lipgloss.Color("#666666")).
		PaddingLeft(2)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#888888")).
		PaddingLeft(2)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)

	return &ChatList{
		list: l,
	}
}

func (cl *ChatList) Update(msg tea.Msg) (*ChatList, tea.Cmd) {
	var cmd tea.Cmd
	cl.list, cmd = cl.list.Update(msg)

	if selected, ok := cl.list.SelectedItem().(ChatItem); ok {
		if cl.onSelect != nil {
			cl.onSelect(selected.Chat)
		}
	}

	return cl, cmd
}

func (cl *ChatList) View() string {
	return cl.list.View()
}

func (cl *ChatList) SetSize(width, height int) {
	cl.width = width
	cl.height = height
	cl.list.SetSize(width, height)
}

func (cl *ChatList) SetChats(chats []*messages.Chat) {
	items := make([]list.Item, len(chats))
	for i, c := range chats {
		items[i] = ChatItem{Chat: c}
	}
	cl.list.SetItems(items)
}

func (cl *ChatList) SetOnSelect(f func(*messages.Chat)) {
	cl.onSelect = f
}

func (cl *ChatList) Focus() {}

func (cl *ChatList) Blur() {}
