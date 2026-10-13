package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/algrvvv/git-go/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// errInteractiveAbort возвращается, когда пользователь отменил выбор (Esc / Ctrl-C).
var errInteractiveAbort = fmt.Errorf("selection aborted")

// runInteractive — интерактивный режим выбора ветки.
//
// Задача метода: получить все ветки (+ описания), красиво показать их
// пользователю и дать выбрать одну. Дальнейшая обработка выбора —
// за вызывающим кодом (см. пометку TODO ниже).
func (a *App) runInteractive() error {
	branches, err := a.git.ListBranches()
	if err != nil {
		return err
	}

	if len(branches) == 0 {
		fmt.Println(dim("  No branches found. Nothing to select."))
		return nil
	}

	current, _ := a.git.GetCurrentBranchName()

	selected, err := selectBranch(branches, current)
	if errors.Is(err, errInteractiveAbort) {
		fmt.Println(dim("  Cancelled."))
		return nil
	}
	if err != nil {
		return err
	}
	if selected == nil {
		return nil
	}

	// TODO: обработка выбора пользователя — дальше твоя реализация.
	// Например: a.git.Checkout(selected)
	fmt.Printf("  %s %s\n", ok("✔"), fmt.Sprintf("selected %s", bold(selected.Name)))
	if selected.Meta != "" {
		fmt.Printf("  %s\n", dim("  "+selected.Meta))
	}

	return a.handleInteractive(selected)
}

// selectBranch — Bubble Tea селектор, либо простой fallback.
func selectBranch(branches []*domain.GitBranch, current string) (*domain.GitBranch, error) {
	sel, err := selectBranchTea(branches, current)
	if err != nil {
		// Bubble Tea не смог открыть TTY (pipe, CI, нет терминала) —
		// откатываемся на простой построчный режим.
		if isNoTTYError(err) {
			return selectBranchSimple(branches, current)
		}
		return nil, err
	}
	return sel, nil
}

func isNoTTYError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "tty") ||
		strings.Contains(msg, "terminal") ||
		strings.Contains(msg, "/dev/tty") ||
		strings.Contains(msg, "inappropriate ioctl")
}

func selectBranchTea(branches []*domain.GitBranch, current string) (*domain.GitBranch, error) {
	model := newSelectModel(branches, current)
	p := tea.NewProgram(model, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return nil, err
	}

	m, ok := final.(selectModel)
	if !ok {
		return nil, fmt.Errorf("unexpected model type %T", final)
	}
	if m.aborted || m.selected == nil {
		return nil, errInteractiveAbort
	}
	return m.selected, nil
}
