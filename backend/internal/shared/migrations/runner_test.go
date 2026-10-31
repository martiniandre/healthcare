package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldRunOnStart(t *testing.T) {
	cases := []struct {
		env, flag string
		want      bool
	}{
		{"development", "", true},
		{"test", "", true},
		{"production", "", false},
		{"production", "true", true},
		{"development", "false", false},
		{"production", "garbage", false},
		{"development", "garbage", true},
		{"production", " 1 ", true},
	}
	for _, migrationCase := range cases {
		if got := ShouldRunOnStart(migrationCase.env, migrationCase.flag); got != migrationCase.want {
			t.Errorf("ShouldRunOnStart(%q,%q)=%v want %v", migrationCase.env, migrationCase.flag, got, migrationCase.want)
		}
	}
}

func TestEvaluateSchemaVersion(t *testing.T) {
	if err := EvaluateSchemaVersion(5, false, 5); err != nil {
		t.Errorf("up to date should pass: %v", err)
	}
	if err := EvaluateSchemaVersion(6, false, 5); err != nil {
		t.Errorf("newer schema should pass: %v", err)
	}
	if err := EvaluateSchemaVersion(4, false, 5); err == nil {
		t.Error("behind schema should fail")
	}
	if err := EvaluateSchemaVersion(5, true, 5); err == nil {
		t.Error("dirty schema should fail")
	}
}

func TestLatestSourceVersion(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"001_a.up.sql", "001_a.down.sql", "002_b.up.sql", "002_b.down.sql", "010_c.up.sql", "010_c.down.sql"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("SELECT 1;"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := LatestSourceVersion("file://" + filepath.ToSlash(dir))
	if err != nil {
		t.Fatal(err)
	}
	if got != 10 {
		t.Errorf("got %d want 10", got)
	}
	if _, err := LatestSourceVersion("file://" + filepath.ToSlash(t.TempDir())); err == nil {
		t.Error("empty dir should error")
	}
}
