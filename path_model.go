package main

import (
	"fmt"
	"mb2-logsfixer/color"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	filepickerView sessionState = iota
	inputView
)

type filePath struct {
	textInput    textinput.Model
	filePicker   filepicker.Model
	selectedFile string
	state        sessionState
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func pathModel() *filePath {
	ti := textinput.New()
	ti.Placeholder = "C:/WoW/Logs/WoWCombatLog.txt"
	ti.TextStyle = color.WhiteColor

	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.AllowedTypes = []string{".txt"}
	fp.AutoHeight = false

	return &filePath{
		filePicker: fp,
		textInput:  ti,
		state:      filepickerView,
		err:        nil,
	}
}
func (m *filePath) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.filePicker.Init())
}

func (m *filePath) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.filePicker.Height = msg.Height - 9
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return mainModel().SwitchView(menuModel(), menuView)
		case "tab":
			if m.state == filepickerView {
				m.state = inputView
				m.textInput.Focus()
				m.textInput.TextStyle = color.WhiteColor
			} else {
				m.state = filepickerView
				m.textInput.Blur()
				m.textInput.TextStyle = color.NormalColor
			}
		case "enter":
			if m.state == inputView {
				path := m.textInput.Value()
				if fileInfo, err := os.Stat(path); err == nil && !fileInfo.IsDir() {
					m.selectedFile = path
					err := SavePath(path)
					if err != nil {
						m.err = err
						return m, clearErrorAfter(2 * time.Second)
					}
					return mainModel().SwitchView(menuModel(), menuView)
				} else {
					m.err = fmt.Errorf("invalid file path selected")
					return m, clearErrorAfter(2 * time.Second)
				}
			}
		}
	case clearErrorMsg:
		m.err = nil
	}
	if m.state == filepickerView {
		m.filePicker, cmd = m.filePicker.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		err := SavePath(path)
		if err != nil {
			m.err = err
			return m, clearErrorAfter(2 * time.Second)
		}
		return mainModel().SwitchView(menuModel(), menuView)
	}

	if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
		m.err = fmt.Errorf(path + ": invalid file path selected")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(cmds...)
}

func (m *filePath) View() string {
	var s string
	filepickerText := lipgloss.JoinVertical(lipgloss.Top, color.NormalColor.Render("text input"), color.NormalColor.Render(m.textInput.View()), color.FocusedStyle.Render("file picker"), color.FocusedStyle.Render(m.filePicker.View()))
	inputText := lipgloss.JoinVertical(lipgloss.Top, color.FocusedStyle.Render("text input"), color.FocusedStyle.Render(m.textInput.View()), color.NormalColor.Render("file picker"), color.NormalColor.Render(m.filePicker.View()))

	if m.state == filepickerView {
		s += filepickerText
	} else {
		s += inputText
	}
	if m.err != nil {
		return color.ContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, s, color.RedColor.Render(m.err.Error()), helpModel(pathView).View()))
	}
	return color.ContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, s, helpModel(pathView).View()))
}

func (m *filePath) currentFocusedModel() string {
	if m.state == filepickerView {
		return "filePicker"
	}
	return "textInput"
}
