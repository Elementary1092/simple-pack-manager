package directory

import (
	"os"
	"path/filepath"
)

func MakeTempDirectoryPath() (tempDir string) {
	wd, err := os.Getwd()
	if err != nil {
		tempDir = filepath.Join(".", "tmp")
		return
	}
	tempDir, err = os.MkdirTemp(wd, "tmp")
	if err != nil {
		tempDir = filepath.Join(".", "tmp")
		return
	}

    return
}

func RemoveDirectory(path string) {
    os.RemoveAll(path)
}

func ExtractVersionFromRemoteFilePath(path string) string {
    return filepath.Base(filepath.Dir(path))
}

func MakeArchivePathName(at string, packet string, version string) string {
	return filepath.Join(at, packet, version, packet)
}

func MakeRemoteArchiveName(packet string, version string, archiveName string) string {
	return filepath.Join(".", packet, version, archiveName)
}

func MakeLatestArchiveLink(packet string) string {
	return filepath.Join(".", packet, "latest")
}

func MakeMetadataPathName(at string, packet string, version string) string {
	return filepath.Join(at, packet, version, "meta")
}

func MakeRemoteMetadataName(packet string, version string) string {
	return filepath.Join(".", "meta", packet, version, "meta")
}

func MakeLatestMetadataLink(packet string) string {
	return filepath.Join(".", "meta", packet, "latest")
}

