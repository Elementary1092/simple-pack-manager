package archiver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestArchive_EmptyRoot(t *testing.T) {
	if _, err := Archive("", "archive", []string{"file"}); !errors.Is(err, ErrInvalidRootPath) {
		t.Fail()
	}
}

func TestArchive_SpacedRoot(t *testing.T) {
	if _, err := Archive("   ", "archive", []string{"file"}); !errors.Is(err, ErrInvalidRootPath) {
		t.Fail()
	}
}

func TestArchive_EmptyArchiveName(t *testing.T) {
	if _, err := Archive("./", "", []string{"file"}); !errors.Is(err, ErrInvalidArchiveName) {
		t.Fail()
	}
}

func TestArchive_SpacedArchiveName(t *testing.T) {
	if _, err := Archive("./", "   ", []string{"file"}); !errors.Is(err, ErrInvalidArchiveName) {
		t.Fail()
	}
}

func TestArchive_EmptyFileName(t *testing.T) {
    tmp := t.TempDir()
    
    if _, err := Archive(tmp, "archive", []string{""}); !errors.Is(err, ErrInvalidFileName) {
        t.Fail()
    }
}

func TestArchieve_SpacedFileName(t *testing.T) {
    tmp := t.TempDir()

    if _, err := Archive(tmp, "archieve", []string{"   "}); !errors.Is(err, ErrInvalidFileName) {
        t.Fail()
    }
}

func TestArchieve_CreateArchieve(t *testing.T) {
    tmp := t.TempDir()
    tmpFS := os.DirFS(tmp)

    if err := os.WriteFile(filepath.Join(tmp, "file"), []byte("some text\n"), 0644); err != nil {
        t.Fatal("Failed on creating temporary file:", err)
    }

    archivePath, err := Archive(tmp, "packet", []string{"file"})
    if err != nil {
        t.Fatal("Failed during archivation:", err)
    }

    if err := fstest.TestFS(tmpFS, filepath.Base(archivePath), "file"); err != nil {
        t.Fatal("Could not find created files:", err)
    }
}

func TestArchieve_CreateArchieveWithDir(t *testing.T) {
    tmp := t.TempDir()
    tmpFS := os.DirFS(tmp)

    dirPath := filepath.Join(tmp, "data")
    if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
        t.Fatal("Failed to create name directory", err)
    }

    if err := os.WriteFile(filepath.Join(dirPath, "file"), []byte("some text\n"), 0644); err != nil {
        t.Fatal("Failed on creating temporary file:", err)
    }

    archivePath, err := Archive(tmp, "packet", []string{filepath.Join("data", "file")})
    if err != nil {
        t.Fatal("Failed during archivation:", err)
    }

    if err := fstest.TestFS(tmpFS, filepath.Base(archivePath), filepath.Join("data", "file")); err != nil {
        t.Fatal("Could not find created files:", err)
    }
}

func TestExtractFrom_EmptyArchiveName(t *testing.T) {
	if err := ExtractFrom("", "archive"); !errors.Is(err, ErrInvalidArchiveName) {
		t.Fail()
	}
}

func TestExtractFrom_SpacedArchiveName(t *testing.T) {
	if err := ExtractFrom("   ", "archive"); !errors.Is(err, ErrInvalidArchiveName) {
		t.Fail()
	}
}

func TestExtractFrom_EmptyExtractPath(t *testing.T) {
	if err := ExtractFrom("./", ""); !errors.Is(err, ErrInvalidFileName) {
		t.Fail()
	}
}

func TestExtractFrom_SpacedExtractPath(t *testing.T) {
	if err := ExtractFrom("./", "   "); !errors.Is(err, ErrInvalidFileName) {
		t.Fail()
	}
}

func TestExtractFrom_ExtractFromAchive(t *testing.T) {
    tmp := t.TempDir()
    tmp = filepath.Join(tmp, "extract")
    extractPath := filepath.Join(tmp, "archive")
    tmpFS := os.DirFS(tmp)

    if err := os.MkdirAll(tmp, os.ModePerm); err != nil {
        t.Fatal("Failed to create temporary directory", err)
    }

    filePath := filepath.Join(tmp, "data")
    if err := os.WriteFile(filePath, []byte("some text"), 0644); err != nil {
        t.Fatal("Failed to create archive:", err)
    }
    
    archivePath, err := Archive(tmp, "archive", []string{filepath.Base(filePath)})
    if err != nil {
        t.Fatal("Failed to create archive:", err)
    }

    err = ExtractFrom(archivePath, extractPath)
    if err != nil {
        t.Fatal("Failed to extract files from the archive:", err)
    }

    if err := fstest.TestFS(tmpFS, filepath.Base(archivePath), filepath.Join("archive", "data")); err != nil {
        t.Fatal("Couldn't find needed files:", err)
    }
}

