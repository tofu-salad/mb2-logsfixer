package color

import "github.com/charmbracelet/lipgloss"

var (
	ContainerStyle = lipgloss.NewStyle().Margin(1, 2)
	FocusedStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ee5396"))
	NormalColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#525252"))
	RedColor       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ee5396"))
	WhiteColor     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	OrangeColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffab91"))
	VioletColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#673ab7"))
	GreenColor     = lipgloss.NewStyle().Foreground(lipgloss.Color("#42be65"))
	BlueColor      = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f62fe"))
)
