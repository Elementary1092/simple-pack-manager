package files

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Elementary1092/pm/internal/packet/parser"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}

func TestCollectLocalFileNames_WithEmptyTargets(t *testing.T) {
    if _, _, err := CollectLocalFileNames([]*parser.Targets{}); !errors.Is(err, ErrNoTargets) {
        t.Fatal(err)
    }
}

func TestCollectLocalFileNames_WithAbsoluteFilePathAndNoExclude(t *testing.T) {
    tmp := t.TempDir()

    fileName := filepath.Join(tmp, "file")
    target := &parser.Targets{
        Path: fileName,
    }
    
    if err := os.WriteFile(fileName, []byte("some text"), 0644); err != nil {
        t.Fatal("Failed to create test file:", err)
    }

    commonRoot, fileNames, err := CollectLocalFileNames([]*parser.Targets{target})
    if err != nil {
        t.Fatal("Failed to collect files:", err)
    }

    if commonRoot != tmp {
        t.Fatalf("Invalid common root: expected='%s'; got='%s'", tmp, commonRoot)
    }

    if len(fileNames) != 1 || fileName != fileNames[0] {
        got := "<empty>"
        if len(fileNames) == 1 {
            got = fileNames[0]
        }
        t.Fatalf("Invalid filename: expected='%s'; got='%s'", fileName, got)
    }
}

func TestCollectLocalFileNames_WithRelativePathAndNoExclude(t *testing.T) {
    currDir, err := os.Getwd()
    if err != nil {
        t.Fatal("Failed to get current working directory:", err)
    }
    tmp := t.TempDir()
    if err := os.Chdir(tmp); err != nil {
        t.Fatal("Failed to change working directory:", err)
    }
    defer os.Chdir(currDir)

    fileName := filepath.Join(".", "file")
    target := &parser.Targets{
        Path: fileName,
    }
    
    if err := os.WriteFile(fileName, []byte("some text"), 0644); err != nil {
        t.Fatal("Failed to create test file:", err)
    }

    commonRoot, fileNames, err := CollectLocalFileNames([]*parser.Targets{target})
    if err != nil {
        t.Fatal("Failed to collect files:", err)
    }

    if commonRoot != tmp {
        t.Fatalf("Invalid common root: expected='%s'; got='%s'", tmp, commonRoot)
    }

    fileName = filepath.Join(tmp, "file")
    if len(fileNames) != 1 || fileName != fileNames[0] {
        got := "<empty>"
        if len(fileNames) == 1 {
            got = fileNames[0]
        }
        t.Fatalf("Invalid filename: expected='%s'; got='%s'", fileName, got)
    }
}

func TestCollectLocalFileNames_WithAbsoluteFilePathAndWithExclude(t *testing.T) {
    tmp := t.TempDir()

    fileName := filepath.Join(tmp, "file")
    target := &parser.Targets{
        Path: fileName,
        Exclude: "file",
    }
    
    if err := os.WriteFile(fileName, []byte("some text"), 0644); err != nil {
        t.Fatal("Failed to create test file:", err)
    }

    commonRoot, fileNames, err := CollectLocalFileNames([]*parser.Targets{target})
    if err != nil {
        t.Fatal("Failed to collect files:", err)
    }

    if commonRoot != tmp {
        t.Fatalf("Invalid common root: expected='%s'; got='%s'", tmp, commonRoot)
    }

    if len(fileNames) != 0 {
        t.Fatalf("Invalid number of files found: expected=0; got=%d", len(fileNames))
    }
}

func TestCollectLocalFileNames_WithRelativePathAndWithExclude(t *testing.T) {
    currDir, err := os.Getwd()
    if err != nil {
        t.Fatal("Failed to get current working directory:", err)
    }
    tmp := t.TempDir()
    if err := os.Chdir(tmp); err != nil {
        t.Fatal("Failed to change working directory:", err)
    }
    defer os.Chdir(currDir)

    fileName := filepath.Join(".", "file")
    target := &parser.Targets{
        Path: fileName,
        Exclude: "file",
    }
    
    if err := os.WriteFile(fileName, []byte("some text"), 0644); err != nil {
        t.Fatal("Failed to create test file:", err)
    }

    commonRoot, fileNames, err := CollectLocalFileNames([]*parser.Targets{target})
    if err != nil {
        t.Fatal("Failed to collect files:", err)
    }

    if commonRoot != tmp {
        t.Fatalf("Invalid common root: expected='%s'; got='%s'", tmp, commonRoot)
    }

    if len(fileNames) != 0 {
        t.Fatalf("Invalid number of files found: expected=0; got=%d", len(fileNames))
    }
}

