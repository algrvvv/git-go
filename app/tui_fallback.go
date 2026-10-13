package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/algrvvv/git-go/domain"
)

// selectBranchSimple — простой fallback-режим (нет TTY): красивый список + промпт.
func selectBranchSimple(branches []*domain.GitBranch, current string) (*domain.GitBranch, error) {
	printBranchListStatic(branches, current)

	fmt.Printf("\n  %s ", dim("Select branch"))
	fmt.Printf("%s ", dim("[1-"+strconv.Itoa(len(branches))+", name, or q to cancel]"))
	fmt.Printf("%s ", cyanBold("›"))

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line == "" || line == "q" || line == "Q" || line == "quit" || line == "exit" {
		return nil, errInteractiveAbort
	}

	// 1. номер
	if n, err := strconv.Atoi(line); err == nil {
		if n >= 1 && n <= len(branches) {
			return branches[n-1], nil
		}
		fmt.Println(red("  ✘ No such number."))
		return nil, errInteractiveAbort
	}

	// 2. точное имя
	lower := strings.ToLower(line)
	for _, b := range branches {
		if strings.ToLower(b.Name) == lower {
			return b, nil
		}
	}

	// 3. подстрока по имени/описанию
	var matches []*domain.GitBranch
	for _, b := range branches {
		if strings.Contains(strings.ToLower(b.Name+" "+b.Meta), lower) {
			matches = append(matches, b)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		fmt.Printf("  %s\n", yellow("  Several branches match, be more specific:"))
		for _, m := range matches {
			fmt.Printf("    %s %s\n", cyanBold("·"), formatBranchLine(m, current, false))
		}
		return nil, errInteractiveAbort
	}

	fmt.Println(red("  ✘ No branch matches."))
	return nil, errInteractiveAbort
}

func printBranchListStatic(branches []*domain.GitBranch, current string) {
	fmt.Printf("\n  %s %s\n", magentaBold("✦"), bold("git-go"))
	fmt.Printf("  %s\n", dim(fmt.Sprintf("%d branches · current: %s", len(branches), currentOrDash(current))))
	fmt.Printf("  %s\n", dim(strings.Repeat("─", 42)))
	for i, b := range branches {
		marker := " "
		if b.Name == current {
			marker = greenBold("★")
		}
		fmt.Printf("  %s %2d  %s %s\n",
			marker, i+1,
			formatBranchLine(b, current, false),
			currentBadge(b.Name, current),
		)
	}
	fmt.Printf("  %s\n", dim(strings.Repeat("─", 42)))
}

// formatBranchLine — рендер строки ветки (fallback).
func formatBranchLine(b *domain.GitBranch, current string, selected bool) string {
	meta := b.Meta
	if meta == "" {
		meta = dimItalic("— no description —")
	} else {
		meta = dim("· " + meta)
	}
	_ = selected
	return fmt.Sprintf("%s %s", bold(b.Name), meta)
}

func currentBadge(name, current string) string {
	if name == current {
		return greenBold("● current")
	}
	return ""
}

func currentOrDash(current string) string {
	if current == "" {
		return "—"
	}
	return current
}
