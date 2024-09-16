package main

import (
	"fmt"
	"mb2-logsfixer/color"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
)

type watch struct {
	filename     string
	processedDir string
	lines        int
	status       string
	spinner      spinner.Model
	err          error
	watcher      *fsnotify.Watcher
	debug        []string
}

func watchModel() *watch {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = color.RedColor
	filename, err := LoadPath()
	if err != nil {
		return &watch{
			err:     fmt.Errorf("error loading saved path: %v", err),
			spinner: s,
		}

	}

	if filename == "" {
		return &watch{
			err:     fmt.Errorf("no saved path found, please set a path first"),
			spinner: s,
		}
	}

	return &watch{
		filename:     filename,
		processedDir: filepath.Join(filepath.Dir("C:/Users/tofu/Games/tauri_mop/Logs/WoWCombatLog.txt"), "MB2"),
		status:       "initializing...",
		spinner:      s,
		debug:        []string{},
	}
}

func (m *watch) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			return initMsg(m.filename)
		},
	)
}

func (m *watch) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.watcher != nil {
				m.watcher.Close()
			}
			return mainModel().SwitchView(menuModel(), menuView)
		case "r":
			return m, processFileCmd(m.filename, m.processedDir)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case initMsg:
		m.status = "watching file: " + m.filename
		var err error
		m.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			m.err = fmt.Errorf("error creating file watcher: %v", err)
			return m, nil
		}
		err = m.watcher.Add(filepath.Dir(m.filename))
		if err != nil {
			m.err = fmt.Errorf("error adding directory to watcher: %v", err)
			return m, nil
		}
		return m, watchCmd(m.watcher, m.filename)

	case watchMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}
		m.status = fmt.Sprintf("change detected. Processing file: %s", m.filename)
		return m, processFileCmd(m.filename, m.processedDir)

	case processFileMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.status = fmt.Sprintf("file processed. modified %d lines in %s", msg.ModifiedLines, msg.Duration)
		return m, watchCmd(m.watcher, m.filename)
	}

	return m, nil
}

func (m *watch) View() string {
	content := fmt.Sprintf(
		"\n %s\n\n%s %s\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("live log"),
		m.spinner.View(),
		m.status,
	)

	if m.err != nil {
		content = fmt.Sprintf("Error: %v\n%s", m.err, content)
	}

	return color.ContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, content,
		helpModel(watchView).View(),
	))
}
