package app

import (
	"github.com/algrvvv/git-go/compare"
	"github.com/algrvvv/git-go/domain"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// selectModel — Bubble Tea модель селектора веток.
type selectModel struct {
	branches []*domain.GitBranch
	current  string
	filtered []compare.Item
	cursor   int
	offset   int
	input    textinput.Model
	width    int
	height   int

	selected *domain.GitBranch
	aborted  bool
}

func newSelectModel(branches []*domain.GitBranch, current string) selectModel {
	ti := textinput.New()
	ti.Placeholder = "type to filter…"
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 120

	return selectModel{
		branches: branches,
		current:  current,
		filtered: compare.FilterAndRank(branches, ""),
		input:    ti,
		width:    80,
		height:   24,
	}
}

func (m selectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit

		case tea.KeyEnter:
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor].Branch
			} else {
				m.aborted = true
			}
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlP:
			m.moveCursor(-1)
			return m, nil

		case tea.KeyDown, tea.KeyCtrlN:
			m.moveCursor(1)
			return m, nil

		case tea.KeyHome:
			m.cursor = 0
			return m, nil

		case tea.KeyEnd:
			if len(m.filtered) > 0 {
				m.cursor = len(m.filtered) - 1
			}
			return m, nil

		case tea.KeyPgUp:
			m.moveCursor(-10)
			return m, nil

		case tea.KeyPgDown:
			m.moveCursor(10)
			return m, nil

		case tea.KeyCtrlU:
			m.input.SetValue("")
			m.refilter()
			return m, nil

		case tea.KeyCtrlW:
			m.input.SetValue(trimLastWord(m.input.Value()))
			m.refilter()
			return m, nil
		}

		// j/k — навигация, пока фильтр пуст; иначе это текст.
		if m.input.Value() == "" && (msg.String() == "j" || msg.String() == "k") {
			if msg.String() == "j" {
				m.moveCursor(1)
			} else {
				m.moveCursor(-1)
			}
			return m, nil
		}

		// всё остальное — в поле фильтра
		before := m.input.Value()
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		if m.input.Value() != before {
			m.refilter()
		}
		return m, cmd
	}

	// прочие сообщения (blink курсора и т.п.) — в input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *selectModel) moveCursor(delta int) {
	n := len(m.filtered)
	if n == 0 {
		m.cursor = 0
		return
	}
	m.cursor += delta
	// зацикливаем — приятно на коротких списках
	if m.cursor < 0 {
		m.cursor = n - 1
	} else if m.cursor >= n {
		m.cursor = 0
	}
}

func (m *selectModel) refilter() {
	m.filtered = compare.FilterAndRank(m.branches, m.input.Value())
	m.cursor = 0
	m.offset = 0
}

func (m selectModel) visibleHeight() int {
	v := m.height - 10
	if v < 3 {
		v = 3
	}
	if v > 20 {
		v = 20
	}
	return v
}

// dividerWidth — ширина разделителя с защитой от нулевого/маленького
// размера окна (первый кадр до WindowSizeMsg).
func (m selectModel) dividerWidth() int {
	return max(10, min(m.width-4, 52))
}
