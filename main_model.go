package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func getWindowSize() tea.Cmd {
	return func() tea.Msg {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return err
		}
		return tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		}
	}
}

type view int

const (
	menuView view = iota
	pathView
	watchView
)

type model struct {
	model       tea.Model
	currentView view
	helpModel   *helpmodel
}

func mainModel() model {
	logPath, _ := LoadPath()
	menuModel := menuModel()
	menuModel.logPath = logPath

	return model{
		model:       menuModel,
		currentView: menuView,
		helpModel:   helpModel(menuView),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.model.Init(), getWindowSize())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m model) View() string {
	return m.model.View()
}

func (m model) SwitchView(newModel tea.Model, newView view) (tea.Model, tea.Cmd) {
	m.model = newModel
	m.currentView = newView
	m.helpModel = helpModel(m.currentView)
	return m.model, tea.Batch(m.model.Init(), getWindowSize())
}
