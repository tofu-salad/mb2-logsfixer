package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type initMsg string

type watchMsg struct {
	Err error
}
type processFileMsg struct {
	ModifiedLines int
	NewLines      int
	TotalLines    int
	NewPosition   int64
	Duration      time.Duration
	Err           error
}

func watchCmd(watcher *fsnotify.Watcher, filename string) tea.Cmd {
	return func() tea.Msg {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return watchMsg{Err: fmt.Errorf("error receiving event")}
			}
			if event.Op&fsnotify.Write == fsnotify.Write && event.Name == filename {
				return watchMsg{}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return watchMsg{Err: fmt.Errorf("error receiving watcher error")}
			}
			return watchMsg{Err: fmt.Errorf("file watcher error: %v", err)}
		}
		return watchMsg{}
	}
}

func processFileCmd(filename, processedDir string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(filename)
		if err != nil {
			return processFileMsg{Err: fmt.Errorf("error opening file: %v", err)}
		}
		defer file.Close()

		err = os.MkdirAll(processedDir, os.ModePerm)
		if err != nil {
			return processFileMsg{Err: fmt.Errorf("error creating subdirectory: %v", err)}
		}

		processedFileName := filepath.Join(processedDir, "WoWCombatLog.txt")
		processedFile, err := os.Create(processedFileName)
		if err != nil {
			return processFileMsg{Err: fmt.Errorf("error creating processed file: %v", err)}
		}
		defer processedFile.Close()

		reader := bufio.NewReader(file)
		writer := bufio.NewWriter(processedFile)
		pattern := regexp.MustCompile(`0x0112([0-9a-fA-F]{12})`)

		var modifiedLines int
		startTime := time.Now()

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return processFileMsg{Err: fmt.Errorf("error reading line: %v", err)}
			}
			modifiedLine := pattern.ReplaceAllStringFunc(line, func(match string) string {
				return "0x0000" + match[6:]
			})
			if modifiedLine != line {
				modifiedLines++
			}
			_, err = writer.WriteString(modifiedLine)
			if err != nil {
				return processFileMsg{Err: fmt.Errorf("error writing to processed file: %v", err)}
			}
		}

		err = writer.Flush()
		if err != nil {
			return processFileMsg{Err: fmt.Errorf("error flushing processed file writer: %v", err)}
		}

		return processFileMsg{
			ModifiedLines: modifiedLines,
			Duration:      time.Since(startTime),
		}
	}
}
