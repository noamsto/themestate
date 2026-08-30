package themestate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPath(t *testing.T) {
	t.Run("XDG_STATE_HOME", func(t *testing.T) {
		stateHome := t.TempDir()
		t.Setenv("XDG_STATE_HOME", stateHome)
		t.Setenv("HOME", t.TempDir())

		want := filepath.Join(stateHome, "theme-state.json")
		if got := Path(); got != want {
			t.Fatalf("Path() = %q, want %q", got, want)
		}
	})

	t.Run("HOME fallback", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("XDG_STATE_HOME", "")
		t.Setenv("HOME", home)

		want := filepath.Join(home, ".local", "state", "theme-state.json")
		if got := Path(); got != want {
			t.Fatalf("Path() = %q, want %q", got, want)
		}
	})
}

func TestDetect(t *testing.T) {
	t.Run("read error", func(t *testing.T) {
		stateHome := filepath.Join(t.TempDir(), "state-home")
		if err := os.WriteFile(stateHome, []byte("not a directory"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("XDG_STATE_HOME", stateHome)
		t.Setenv("HOME", t.TempDir())

		if got := Detect(); got != "dark" {
			t.Fatalf("Detect() = %q, want %q", got, "dark")
		}
	})

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{name: "light", content: `{"theme":"light"}`, want: "light"},
		{name: "dark", content: `{"theme":"dark"}`, want: "dark"},
		{name: "missing", want: "dark"},
		{name: "malformed JSON", content: `{"theme":`, want: "dark"},
		{name: "empty theme", content: `{"theme":""}`, want: "dark"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateHome := t.TempDir()
			t.Setenv("XDG_STATE_HOME", stateHome)
			t.Setenv("HOME", t.TempDir())
			if tt.content != "" {
				if err := os.WriteFile(filepath.Join(stateHome, "theme-state.json"), []byte(tt.content), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			if got := Detect(); got != tt.want {
				t.Fatalf("Detect() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestModTime(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "theme-state.json")
		if err := os.WriteFile(path, []byte(`{"theme":"dark"}`), 0o600); err != nil {
			t.Fatal(err)
		}
		want := time.Now().Add(-time.Hour).Truncate(time.Second)
		if err := os.Chtimes(path, want, want); err != nil {
			t.Fatal(err)
		}

		got, err := ModTime(path)
		if err != nil {
			t.Fatalf("ModTime() error = %v", err)
		}
		if !got.Equal(want) {
			t.Fatalf("ModTime() = %v, want %v", got, want)
		}
	})

	t.Run("absent", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing")
		got, err := ModTime(path)
		if err == nil {
			t.Fatal("ModTime() error = nil, want an error")
		}
		if !got.IsZero() {
			t.Fatalf("ModTime() time = %v, want zero", got)
		}
	})
}
