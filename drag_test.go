package ripple

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPrototypeExtendSelectionTo(t *testing.T) {
	for _, mode := range []WrapMode{WrapCharacter, WrapWord} {
		e := New()
		e.SetSize(10, 8)
		e.SetWrapMode(mode)
		text := "hello world again\nsecond hard line\nlast"
		e.SetValue(text)
		clip := &fakeClip{s: "unchanged"}
		e.SetClipboard(clip)
		e.PlaceCursor(0, 2)
		e.ExtendSelectionTo(1, 3)
		end := 13
		if mode == WrapWord {
			end = 9
		}
		if e.SelectedText() != text[2:end] {
			t.Fatalf("mode %v: forward selection %q", mode, e.SelectedText())
		}
		e.ExtendSelectionTo(0, 0)
		if e.SelectedText() != "he" || e.selAnchor != 2 {
			t.Fatal("backward drag lost anchor")
		}
		e.ExtendSelectionTo(0, 2)
		e.ExtendSelectionTo(3, 2)
		if e.selAnchor != 2 || !strings.Contains(e.SelectedText(), "\n") || e.Value() != text || clip.s != "unchanged" {
			t.Fatal("cross-row selection changed anchor/text/clipboard")
		}
		selected := e.SelectedText()
		cmd := apply(&e, tea.KeyMsg{Type: tea.KeyCtrlC})
		if cmd == nil {
			t.Fatal("missing copy command")
		}
		cmd()
		if clip.s != selected {
			t.Fatal("copy did not use drag selection")
		}
		e.PlaceCursor(0, 1)
		if e.SelectedText() != "" || e.CursorIndex() != 1 {
			t.Fatal("subsequent click")
		}
		e.SetSize(10, 2)
		e.ScrollRows(2)
		top := e.ViewportTop()
		e.PlaceCursor(0, 1)
		anchor := e.CursorIndex()
		end = e.indexAtVisual(top+1, 3)
		e.ExtendSelectionTo(1, 3)
		if e.selAnchor != anchor || e.CursorIndex() != end || e.ViewportTop() != top {
			t.Fatal("scrolled selection")
		}
		e.ExtendSelectionTo(999, 999)
		if e.ViewportTop() != top {
			t.Fatal("out-of-bounds selection auto-scrolled")
		}
		if e.Value() != text {
			t.Fatal("drag modified document")
		}
	}
}
