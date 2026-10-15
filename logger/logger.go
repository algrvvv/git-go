package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	colorGray  = "\033[90m"
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorBlue  = "\033[34m"
)

func NewLogger(debug bool) *slog.Logger {
	if debug {
		return slog.New(newDebugHandler(os.Stdout))
	}

	return slog.New(slog.DiscardHandler)
}

type PrettyHandler struct {
	w io.Writer
}

func newDebugHandler(w io.Writer) *PrettyHandler {
	return &PrettyHandler{w: w}
}

// Enabled определяет, какие уровни логов обрабатывать (в данном случае все от INFO)
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelDebug
}

// Handle парсит запись и выводит её в нашем кастомном формате
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// 1. Форматируем время в нормальный человеческий вид и красим в серый
	timeStr := r.Time.Format("02.01.2006 15:04:05")
	fmt.Fprintf(h.w, "%s%s%s ", colorGray, timeStr, colorReset)

	// 2. Раскрашиваем уровень лога для наглядности
	var levelStr string
	switch r.Level {
	case slog.LevelDebug:
		levelStr = colorBlue + "DEBUG" + colorReset
	case slog.LevelInfo:
		levelStr = colorGreen + "INFO " + colorReset
	case slog.LevelWarn:
		levelStr = colorRed + "WARN " + colorReset
	case slog.LevelError:
		levelStr = colorRed + "ERROR" + colorReset
	default:
		levelStr = r.Level.String()
	}
	fmt.Fprintf(h.w, "[%s] ", levelStr)

	// 3. Выводим сообщение "как есть" (без экранирования \n)
	fmt.Fprint(h.w, r.Message)

	// 4. Собираем и выводим дополнительные атрибуты (key=value), если они есть
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		// Чтобы внутри значений атрибутов тоже работали переносы строк
		valStr := a.Value.String()
		attrs = append(attrs, fmt.Sprintf("%s%s%s=%s", colorGray, a.Key, colorReset, valStr))
		return true
	})

	if len(attrs) > 0 {
		// Выводим атрибуты с новой строки с отступом, чтобы не кашеварить
		fmt.Fprintf(h.w, " %s", strings.Join(attrs, " "))
	}

	// Обязательно добавляем перенос строки в самом конце лога
	fmt.Fprintln(h.w)
	return nil
}

// WithAttrs и WithGroup нужны для интерфейса slog.Handler (в данном примере оставим базовыми)
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *PrettyHandler) WithGroup(name string) slog.Handler       { return h }
