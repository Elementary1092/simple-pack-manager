package updatecmd

import (
	"context"
	"io"
)

type updateCommand struct {
    data io.Reader
}

func NewUpdateCommand(data io.Reader) *updateCommand {
    if data == nil {
        return nil
    }

    return &updateCommand{
        data: data,
    }
}

func (up *updateCommand) Execute(ctx context.Context) error {

    return nil
}

