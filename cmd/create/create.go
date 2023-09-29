package createcmd

import (
	"context"
	"io"
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

    return nil
}

