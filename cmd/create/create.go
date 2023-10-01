package createcmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Elementary1092/pm/internal/adapter/pmssh"
	"github.com/Elementary1092/pm/internal/directory"
	"github.com/Elementary1092/pm/internal/packet/archiver"
	"github.com/Elementary1092/pm/internal/packet/files"
	"github.com/Elementary1092/pm/internal/packet/parser"
)

var (
	ErrInternalError = errors.New("internal command error")
)

type createCommand struct {
	data io.Reader
}

func NewCreateCommand(data io.Reader) *createCommand {
	if data == nil {
		return nil
	}

	return &createCommand{
		data: data,
	}
}

func (cr *createCommand) Execute(ctx context.Context) error {
    fmt.Println("Parsing packet description.")
	description, err := parser.ParsePacket(cr.data)
	if err != nil {
		return err
	}

    fmt.Println("Collecting local files.")
	base, filenames, err := files.CollectLocalFileNames(description.Targets)
	if err != nil {
		return err
	}

    tempPath := directory.MakeTempDirectoryPath()
	defer directory.RemoveDirectory(tempPath)

    fmt.Println("Creating archive.")
	archiveName := directory.MakeArchivePathName(tempPath, description.Name, description.Version)
	archiveName, err = archiver.Archive(base, archiveName, filenames)
	if err != nil {
		return err
	}
	defer os.Remove(archiveName)

	metaFilePath := directory.MakeMetadataPathName(tempPath, description.Name, description.Version)
	if err := makeMetadataFile(metaFilePath, description.Packets); err != nil {
		return ErrInternalError
	}
	defer os.Remove(metaFilePath)

	if err := pmssh.Connect(ctx); err != nil {
		return err
	}
	defer pmssh.Close(ctx)

    fmt.Println("Uploading files...")
	remoteArchive := directory.MakeRemoteArchiveName(description.Name, description.Version, description.Name)
	if err := pmssh.Upload(ctx, remoteArchive, archiveName); err != nil {
		return err
	}

	linkName := directory.MakeLatestArchiveLink(description.Name)
	if err := pmssh.CreateSymbolicLink(ctx, linkName, archiveName); err != nil {
		return err
	}

	remoteMetadata := directory.MakeRemoteMetadataName(description.Name, description.Version)
	if err := pmssh.Upload(ctx, remoteMetadata, metaFilePath); err != nil {
		return err
	}

	metaLinkName := directory.MakeLatestMetadataLink(description.Name)
	if err := pmssh.CreateSymbolicLink(ctx, metaLinkName, metaFilePath); err != nil {
		return err
	}

	return nil
}

func makeMetadataFile(filePath string, data any) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	return encoder.Encode(data)
}

