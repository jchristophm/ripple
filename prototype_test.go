package ripple

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPrototypeWrapModes(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		width      int
		want       []string
	}{
		{"words", "hello world again", 10, []string{"hello ", "world ", "again"}},
		{"long word", "abcdefghijkl", 5, []string{"abcde", "fghij", "kl"}},
		{"url", "https://example.org/path", 10, []string{"https://ex", "ample.org/", "path"}},
		{"whitespace and newline", "ab  cd\nef", 5, []string{"ab  ", "cd", "ef"}},
		{"wide", "ab 界界 x", 6, []string{"ab ", "界界 x", ""}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := New()
			e.SetSize(tc.width, 10)
			e.SetValue(tc.text)
			if e.wrapMode != WrapCharacter {
				t.Fatal("default changed")
			}
			character := e.visualLines()
			e.SetWrapMode(WrapWord)
			var got []string
			for _, line := range e.visualLines() {
				got = append(got, string(e.text[line.start:line.end]))
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("rows = %#v, want %#v", got, tc.want)
			}
			if e.Value() != tc.text {
				t.Fatal("wrap edited buffer")
			}
			e.SetWrapMode(WrapCharacter)
			if !reflect.DeepEqual(character, e.visualLines()) {
				t.Fatal("character wrapping changed")
			}
		})
	}
}

func TestPrototypeCursorSelectionAndScroll(t *testing.T) {
	for _, mode := range []WrapMode{WrapCharacter, WrapWord} {
		e := New()
		e.SetSize(10, 2)
		e.SetWrapMode(mode)
		text := "hello world again\n" + strings.Repeat("another row here\n", 5)
		e.SetValue(text)
		e.ScrollRows(-1000)
		e.PlaceCursor(0, 2)
		if e.CursorIndex() != 2 {
			t.Fatal("placement")
		}
		apply(&e, tea.KeyMsg{Type: tea.KeyShiftDown})
		cursor, selection := e.CursorIndex(), e.SelectedText()
		if selection == "" {
			t.Fatal("wrapped vertical selection is empty")
		}
		e.ScrollRows(3)
		if e.ViewportTop() != 3 || e.CursorIndex() != cursor || e.SelectedText() != selection {
			t.Fatal("scroll changed caret/selection")
		}
		_ = e.View(Options{})
		if e.ViewportTop() != 3 {
			t.Fatal("render undid scroll")
		}
		want := e.indexAtVisual(e.ViewportTop()+1, 3)
		e.PlaceCursor(1, 3)
		if e.CursorIndex() != want || e.SelectedText() != "" || e.Value() != text {
			t.Fatal("placement after scroll")
		}
		e.ScrollRows(1000)
		if e.ViewportTop() != max(0, len(e.visualLines())-2) {
			t.Fatal("scroll lower clamp")
		}
		e.ScrollRows(-1000)
		want = e.indexAtVisual(0, 100)
		e.PlaceCursor(0, 100)
		if e.CursorIndex() != want {
			t.Fatal("placement beyond visual end")
		}
		apply(&e, tea.KeyMsg{Type: tea.KeyShiftRight})
		cursor, selection = e.CursorIndex(), e.SelectedText()
		e.SetWrapMode(WrapWord)
		if e.CursorIndex() != cursor || e.SelectedText() != selection || e.Value() != text {
			t.Fatal("wrap state preservation")
		}
	}
}
