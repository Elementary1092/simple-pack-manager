package archiver

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidRootPath              = errors.New("invalid root path")
	ErrInvalidArchiveName           = errors.New("invalid archive name")
	ErrInvalidFileName              = errors.New("invalid or empty file name")
	ErrFailedToCreateArchieve       = errors.New("failed to create archive")
	ErrFailedToOpenFile             = errors.New("failed to open file")
	ErrFailedToArchiveFile          = errors.New("failed to archive file")
	ErrFailedToCreateCompressedFile = errors.New("failed to create compressed file")
	ErrFailedToCreateDirectory      = errors.New("failed to create directory")
	ErrFailedToCreateFile           = errors.New("failed to create file")
	ErrFailedToExtractFile          = errors.New("failed to extract file")
)

// Archive creates an archive and writes contents of all files given in fileNames to it.
// root should be a common prefix of all files listed in fileNames.
// Archive path with name is returned.
// Archive file extenstion is appended by a function.
func Archive(root string, archieveName string, fileNames []string) (string, error) {
	root = strings.TrimSpace(root)
	archieveName = strings.TrimSpace(archieveName)

	if root == "" {
		return "", ErrInvalidRootPath
	}

	if archieveName == "" {
		return "", ErrInvalidArchiveName
	}

	archieveName += ".zip"

    if err := os.MkdirAll(filepath.Dir(archieveName), 0777); err != nil {
        return "", ErrFailedToCreateArchieve
    }

	archive, err := os.OpenFile(archieveName, os.O_TRUNC | os.O_CREATE | os.O_RDWR, os.ModePerm)
	// could have wrapped err from os
	if err != nil {
		return "", ErrFailedToCreateArchieve
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	for _, fileName := range fileNames {
		fileName = strings.TrimSpace(fileName)
		if fileName == "" {
			return "", ErrInvalidFileName
		}

		f, err := os.Open(fileName)
		if err != nil {
			return "", ErrFailedToOpenFile
		}
		defer f.Close()

        pathInArchive, err := filepath.Rel(root, fileName)
        if err != nil {
            return "", ErrFailedToArchiveFile
        }
		compressed, err := zipWriter.Create(pathInArchive)
		if err != nil {
			return "", ErrFailedToCreateCompressedFile
		}

		if _, err := io.Copy(compressed, f); err != nil {
			return "", ErrFailedToCreateCompressedFile
		}
	}

	return archieveName, nil
}

// ExtractFrom opens archiveFullName file and writes its contents to the extractToPath
func ExtractFrom(archiveFullName string, extractToPath string) error {
	archiveFullName = strings.TrimSpace(archiveFullName)
	extractToPath = strings.TrimSpace(extractToPath)

	if archiveFullName == "" {
		return ErrInvalidArchiveName
	}

	if extractToPath == "" {
		return ErrInvalidFileName
	}

	zipReader, err := zip.OpenReader(archiveFullName)
	if err != nil {
		return ErrFailedToOpenFile
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fileName := strings.TrimSpace(strings.TrimLeft(f.Name, "./\\"))
		if fileName == "" {
			return ErrInvalidFileName
		}
		// assuming that filepaths from a zip file are correct
		filePath := filepath.Join(extractToPath, fileName)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return ErrFailedToCreateDirectory
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return ErrFailedToCreateDirectory
		}

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return ErrFailedToCreateFile
		}
		defer file.Close()

		compressedFile, err := f.Open()
		if err != nil {
			return ErrFailedToOpenFile
		}
		defer compressedFile.Close()

		if _, err := io.Copy(file, compressedFile); err != nil {
			return ErrFailedToExtractFile
		}
	}

	return nil
}
