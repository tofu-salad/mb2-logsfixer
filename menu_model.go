package main

import (
	"mb2-logsfixer/color"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type menu struct {
	err     error
	list    list.Model
	logPath string
}

func menuModel() *menu {
	logPath, _ := LoadPath()
	desc := "where your combat log file is"
	if logPath != "" {
		desc = logPath
	}

	items := []list.Item{
		item{title: "[1] WoWCombatLog.txt path", desc: desc},
		item{title: "[2] Live Logs", desc: "file watching for WCL Live Logs."},
		item{title: "[3] Fix Logs", desc: "manually fix the logs to upload them."},
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.Title = "Mistblade 2 WCL Fixer"
	return &menu{
		list:    list,
		err:     nil,
		logPath: logPath,
	}
}

func (m *menu) Init() tea.Cmd {
	return nil
}

func (m *menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := color.ContainerStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			return mainModel().SwitchView(pathModel(), pathView)
		case "2":

			return mainModel().SwitchView(watchModel(), watchView)
			// TODO: manual model
		// case "3"
		// 			return mainModel().SwitchView(manualModel())
		case "enter":
			selectedIndex := m.list.Index()
			switch selectedIndex {
			case 0:
				return mainModel().SwitchView(pathModel(), pathView)
			case 1:
				return mainModel().SwitchView(watchModel(), watchView)
				// TODO: Manual model
				// case 2:
				// 	return mainModel().SwitchView(manualModel())
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *menu) View() string {
	return color.ContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, m.list.View(),
		helpModel(menuView).View(),
	))
}
