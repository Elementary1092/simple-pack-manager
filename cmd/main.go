package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	createcmd "github.com/Elementary1092/pm/cmd/create"
	updatecmd "github.com/Elementary1092/pm/cmd/update"
)

const helpPrompt = `pm -create <filename> - create package from package declaration files

pm -update <filename> - update package from package description files`

const succeededPrompt = `Operation is successful.`

type Command interface {
    Execute(ctx context.Context) error
}

func main() {
    var command Command
    var file *os.File
    // Not the best method to parse commands 
    // (cobra package could be used instead of this and validator functions could be extracted), 
    // but it makes development easier
    flag.Func("create", "Upload package to the server", func(s string) error {
        // not the best way to limit number of command
        if command != nil {
            return errors.New("Expected only 1 command at a time")
        } 

        if err := validatePath(s); err != nil {
            return err
        }
        
        f, err := os.Open(s)
        if err != nil {
            return fmt.Errorf("Failed to open file %s", s)
        }
        
        file = f
        command = createcmd.NewCreateCommand(f)

        return nil
    })
    flag.Func("update", "Fetch specified packages from the server", func (s string) error {
        if command != nil {
            return errors.New("Expected only 1 command at a time")
        } 

        if err := validatePath(s); err != nil {
            return err
        }

        f, err := os.Open(s)
        if err != nil {
            return fmt.Errorf("Failed to open file %s", s)
        }
        
        file = f
        command = updatecmd.NewUpdateCommand(f)

        return nil
    })
    flag.Parse()

    if command == nil {
        fmt.Println(helpPrompt)
        return
    }
    defer file.Close()
    
    if err := command.Execute(context.Background()); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(succeededPrompt)
    }
}

// Validating files should have been performed in commands constructor,
// but logic of validation for create and update is the same.
// So, it was decided to perform this validation in flag parser.
func validatePath(s string) error {
    fileInfo, err := os.Lstat(s)
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            return fmt.Errorf("Could not find file %s", s)
        }

        return fmt.Errorf("Failed to check information about file %s", s)
    }
    
    if !fileInfo.Mode().IsRegular() {
        return errors.New("Unsupported file type")
    }

    return nil
}

