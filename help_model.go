package main

import (
	"mb2-logsfixer/color"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Left        key.Binding
	Right       key.Binding
	Help        key.Binding
	Quit        key.Binding
	Tab         key.Binding
	OneTwoThree key.Binding
	R           key.Binding
	currentView view
}

func (k keyMap) ShortHelp() []key.Binding {
	switch k.currentView {
	default:
		return []key.Binding{k.Help, k.Quit}
	}

}

func (k keyMap) FullHelp() [][]key.Binding {
	switch k.currentView {

	case pathView:
		return [][]key.Binding{
			{k.Tab, k.Quit},
		}
	case watchView:
		return [][]key.Binding{
			{k.R, k.Quit},
		}
	default:
		return [][]key.Binding{
			{k.OneTwoThree, k.Up, k.Down, k.Left, k.Right},
			{k.Enter, k.Quit},
		}

	}
}

type helpmodel struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func helpModel(s view) *helpmodel {
	var keys keyMap

	switch s {
	case menuView:
		keys = keyMap{
			OneTwoThree: key.NewBinding(key.WithKeys("1", "2", "3"), key.WithHelp("1/2/3", "select")),
			Enter:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q/esc/ctrl+c", "quit"),
			),
			currentView: menuView,
		}
	case pathView:
		keys = keyMap{
			Tab: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "switch between explorer and input"),
			),
			R: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "manually run log fixer"),
			),
			Quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc/ctrl+c", "go back to menu"),
			),
			currentView: pathView,
		}
	case watchView:
		keys = keyMap{
			R: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "manually run log fixer"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q/esc/ctrl+c", "go back to menu"),
			),
			currentView: watchView,
		}
	}
	return &helpmodel{
		keys:       keys,
		help:       help.New(),
		inputStyle: color.RedColor,
	}
}

func (m *helpmodel) Init() tea.Cmd {
	return nil
}

func (m *helpmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	}

	return m, nil
}

func (m *helpmodel) View() string {
	helpView := m.help.FullHelpView(m.keys.FullHelp())

	return helpView
}
