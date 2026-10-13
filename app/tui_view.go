package app

import (
	"fmt"
	"strings"

	"github.com/algrvvv/git-go/domain"
	"github.com/charmbracelet/lipgloss"
)

func (m selectModel) View() string {
	// держим курсор в видимой зоне
	visible := m.visibleHeight()
	offset := m.offset
	if m.cursor < offset {
		offset = m.cursor
	}
	if m.cursor >= offset+visible {
		offset = m.cursor - visible + 1
	}
	if offset < 0 {
		offset = 0
	}

	var b strings.Builder

	// --- шапка ---
	b.WriteString("  " + styleTitle.Render("✦") + " " + styleBold.Render("git-go") + "\n")
	b.WriteString("  " + styleDim.Render(fmt.Sprintf("%d branches · current: %s · %d shown",
		len(m.branches), currentOrDash(m.current), len(m.filtered))) + "\n")
	b.WriteString("  " + styleDim.Render(strings.Repeat("─", m.dividerWidth())) + "\n")

	// --- строка поиска ---
	b.WriteString("  " + styleSearchIc.Render("⌕") + " " + m.input.View() + "\n")
	b.WriteString("  " + styleDim.Render(strings.Repeat("─", m.dividerWidth())) + "\n")

	// --- список ---
	if len(m.filtered) == 0 {
		b.WriteString("\n  " + styleDim.Render("  No matches. Keep typing or clear the filter (ctrl+u).") + "\n")
	} else {
		end := offset + visible
		if end > len(m.filtered) {
			end = len(m.filtered)
		}
		if offset > 0 {
			b.WriteString("  " + styleDim.Render(fmt.Sprintf("  ▲ %d more above", offset)) + "\n")
		}
		nameW := 0
		for i := offset; i < end; i++ {
			if w := lipgloss.Width(m.filtered[i].Branch.Name); w > nameW {
				nameW = w
			}
		}
		for i := offset; i < end; i++ {
			b.WriteString(m.renderRow(m.filtered[i].Branch, i == m.cursor, nameW) + "\n")
		}
		if end < len(m.filtered) {
			b.WriteString("  " + styleDim.Render(fmt.Sprintf("  ▼ %d more below", len(m.filtered)-end)) + "\n")
		}
	}

	// --- подвал ---
	b.WriteString("\n")
	b.WriteString("  " + styleDim.Render("↑↓ navigate · type to filter · enter select · esc cancel") + "\n")
	if len(m.filtered) > 0 {
		b.WriteString("  " + styleDim.Render(fmt.Sprintf("%d/%d", m.cursor+1, len(m.filtered))) + "\n")
	}

	return b.String()
}

func (m selectModel) renderRow(b *domain.GitBranch, isSelected bool, nameW int) string {
	meta := b.Meta
	metaStyle := styleMeta
	if meta == "" {
		meta = "— no description —"
		metaStyle = styleNoDesc
	}

	// обрезаем описание под ширину экрана
	maxMeta := m.width - nameW - 14
	if maxMeta < 10 {
		maxMeta = 10
	}
	if lipgloss.Width(meta) > maxMeta {
		meta = truncateRunes(meta, maxMeta-1) + "…"
	}

	padded := b.Name + strings.Repeat(" ", max(0, nameW-lipgloss.Width(b.Name)))

	var line string
	if isSelected {
		line = "  " + styleCyan.Render("❯") + " " +
			styleCyan.Render(padded) + " " +
			metaStyle.Render("· "+meta)
		if b.Name == m.current {
			line += " " + styleGreen.Render("● current")
		}
		return styleSelRow.Render(line)
	}

	name := padded
	if b.Name == m.current {
		name = styleGreen.Render(padded)
	}
	line = "    " + name + " " + metaStyle.Render("· "+meta)
	if b.Name == m.current {
		line += " " + styleGreen.Render("● current")
	}
	return line
}
