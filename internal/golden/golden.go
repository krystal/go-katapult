package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Update() bool {
	return os.Getenv("UPDATE_GOLDEN") == "1"
}

func Path(t *testing.T) string {
	return filepath.Join("testdata", filepath.FromSlash(t.Name())+".golden")
}

func Get(t *testing.T) []byte {
	gp := Path(t)
	g, err := ioutil.ReadFile(gp)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}

	return g
}

func Set(t *testing.T, got []byte) {
	gp := Path(t)
	dir := filepath.Dir(gp)

	t.Logf("updating .golden file: %s", gp)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to update .golden directory: %s", err)
	}

	if err := ioutil.WriteFile(gp, got, 0o644); err != nil { //nolint:gosec
		t.Fatalf("failed to update .golden file: %s", err)
	}
}
