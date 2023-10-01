package updatecmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Elementary1092/pm/internal/adapter/pmssh"
	"github.com/Elementary1092/pm/internal/directory"
	"github.com/Elementary1092/pm/internal/packet/archiver"
	"github.com/Elementary1092/pm/internal/packet/parser"
	"github.com/Elementary1092/pm/internal/version"
)

var (
    ErrFailedToCreateDestinationDir = errors.New("failed to create destination directory")
)

type updateCommand struct {
    data io.Reader
    name string
}

func NewUpdateCommand(data io.Reader, name string) *updateCommand {
    if data == nil {
        return nil
    }

    return &updateCommand{
        data: data,
        name: name,
    }
}

// Execute assumes that all dependencies are listed in a file
func (up *updateCommand) Execute(ctx context.Context) error {
    fmt.Println("Parsing package description.")
    description, err := parser.ParsePackage(up.data)
    if err != nil {
        return err
    }

    err = pmssh.Connect(ctx)
    if err != nil {
        return err
    }
    defer pmssh.Close(ctx)

    tempPath := directory.MakeTempDirectoryPath()
    defer directory.RemoveDirectory(tempPath)

    wd, err := os.Getwd()
    if err != nil {
        wd = "."
    }
    packDestination := filepath.Join(wd, up.name)
    if err := os.MkdirAll(packDestination, os.ModePerm); err != nil {
        return ErrFailedToCreateDestinationDir
    }

    for _, pack := range description.Packages {
        versionToGet := pack.Version

        var err error
        if versionToGet == "" {
            versionToGet, err = up.getLatest(ctx, pack.Name)
            if err != nil {
                return fmt.Errorf("failed to find satisfying version ('lastest') of a package '%s'", pack.Name)
            }
        } else {
            versionToGet, err = up.getSpecificVersion(ctx, pack.Name, versionToGet)
            if err != nil {
                return fmt.Errorf("failed to find satisfying version ('%s') of a package '%s'", pack.Version, pack.Name)
            }
        }

        fmt.Printf("Fetching package '%s' of version '%s'\n", pack.Name, versionToGet)
        archPath := directory.MakeArchivePathName(tempPath, pack.Name, versionToGet)
        archNamePath := filepath.Join(archPath, pack.Name+".zip")
        if err := os.MkdirAll(archPath, os.ModePerm); err != nil {
            return ErrFailedToCreateDestinationDir
        }

        remoteArchName := directory.MakeRemoteArchiveName(pack.Name, versionToGet, pack.Name)

        err = pmssh.Download(ctx, remoteArchName, archNamePath)
        if err != nil {
            return err
        }

        fmt.Println("Extracting package", pack.Name)
        err = archiver.ExtractFrom(archNamePath, filepath.Join(packDestination, pack.Name))
        if err != nil {
            return err
        }
    }


    return nil
}

func (up *updateCommand) getLatest(ctx context.Context, packName string) (string, error) {
    lastestLink := directory.MakeLatestArchiveLink(packName)

    filePath, err := pmssh.ResolveLink(ctx, lastestLink)
    if err != nil {
        return "", err
    }

    return directory.ExtractVersionFromRemoteFilePath(filePath), nil
}

// Assumes that all versions are correct
func (up *updateCommand) getSpecificVersion(ctx context.Context, packName string, ver string) (string, error) {
    versionType := version.Type(ver)
    ver = version.Clean(ver)
    if versionType == version.LessOrEqual {
        ver = up.findLessOrEqualVersion(ctx, packName, ver)
    } else if versionType == version.GreaterOrEqual {
        ver = up.findGreaterOrEqualVersion(ctx, packName, ver)
    }
        
    if ver == "" {
        return "", errors.New("no satisfying version")
    }

        
    return ver, nil
}

// Assumes that user provided version is present in the server
func (up *updateCommand) findLessOrEqualVersion(ctx context.Context, packName string, maxVer string) string {
    lastestVersion, err := up.getLatest(ctx, packName)
    if err != nil {
        return ""
    }

    verType, err := version.CompareVersions(maxVer, lastestVersion)
    if err != nil {
        return ""
    }

    if verType == version.Less || verType == version.Exact {
        return maxVer
    }

    return lastestVersion
}

func (up *updateCommand) findGreaterOrEqualVersion(ctx context.Context, packName string, minVer string) string {
    latestVersion, err := up.getLatest(ctx, packName)
    if err != nil {
        return ""
    }

    verType, err := version.CompareVersions(minVer, latestVersion)
    if verType == version.Greater || err != nil {
        return ""
    }

    return latestVersion
}

