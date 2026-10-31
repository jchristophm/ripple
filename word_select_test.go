package ripple

import "testing"

func TestPrototypeSelectWordAt(t *testing.T) {
	for _, mode := range []WrapMode{WrapCharacter, WrapWord} {
		e := New()
		e.SetSize(5, 10)
		e.SetWrapMode(mode)
		text := "hello_longword, next\n\n last"
		e.SetValue(text)
		clip := &fakeClip{s: "unchanged"}
		e.SetClipboard(clip)
		e.SelectWordAt(1, 2) // middle of a word split across visual rows
		if e.SelectedText() != "hello_longword" {
			t.Fatalf("mode %v selection %q", mode, e.SelectedText())
		}
		for _, tc := range []struct {
			index int
			want  string
		}{{0, "hello_longword"}, {14, ","}, {15, ""}, {16, "next"}, {20, ""}, {21, ""}, {22, ""}, {23, "last"}, {len([]rune(text)), ""}} {
			row, col := e.visualPosition(tc.index)
			e.SelectWordAt(row, col)
			if e.SelectedText() != tc.want {
				t.Fatalf("mode %v index %d got %q want %q", mode, tc.index, e.SelectedText(), tc.want)
			}
		}
		e.SelectWordAt(0, 100) // blank area after the first soft row must not select the next word
		if e.SelectedText() != "" {
			t.Fatal("blank cell selected adjacent word")
		}
		e.SetSize(5, 2)
		e.ScrollRows(2)
		top := e.ViewportTop()
		e.SelectWordAt(0, 1)
		if e.SelectedText() != "hello_longword" || e.ViewportTop() != top {
			t.Fatal("selection after scroll")
		}
		if e.Value() != text || clip.s != "unchanged" {
			t.Fatal("word selection modified text/clipboard")
		}
		e.SetValue("")
		e.SelectWordAt(100, 100)
		if e.SelectedText() != "" || e.Value() != "" {
			t.Fatal("empty selection")
		}
	}
}
