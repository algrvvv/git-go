package app

import "github.com/charmbracelet/lipgloss"

// Стили Bubble Tea интерфейса (lipgloss).
var (
	styleTitle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	styleBold     = lipgloss.NewStyle().Bold(true)
	styleDim      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	styleCyan     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
	styleGreen    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	styleSelRow   = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	styleMeta     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	styleNoDesc   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	styleSearchIc = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
)

// Ниже — ручные ANSI-стили для не-TTY fallback режима (fmt.Printf вне Bubble Tea).
// Две системы стилей оставлены намеренно: lipgloss требует терминал,
// а fallback печатает напрямую и не должен тянуть рендер-движок.
const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiItalic = "\x1b[3m"
)

func bold(s string) string        { return ansiBold + s + ansiReset }
func dim(s string) string         { return ansiDim + s + ansiReset }
func dimItalic(s string) string   { return ansiDim + ansiItalic + s + ansiReset }
func cyanBold(s string) string    { return "\x1b[1;36m" + s + ansiReset }
func greenBold(s string) string   { return "\x1b[1;32m" + s + ansiReset }
func magentaBold(s string) string { return "\x1b[1;35m" + s + ansiReset }
func yellow(s string) string      { return "\x1b[33m" + s + ansiReset }
func red(s string) string         { return "\x1b[31m" + s + ansiReset }
func ok(s string) string          { return "\x1b[32m" + s + ansiReset }
