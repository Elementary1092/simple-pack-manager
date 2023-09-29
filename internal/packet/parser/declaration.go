package parser

import (
	"encoding/json"
	"errors"
	"io"

	validate "github.com/Elementary1092/pm/internal/packet/validator"
)

var (
	ErrInvalidPacketDescriptionFormat = errors.New("invalid packet description format")
	ErrInvalidPacketDescription       = errors.New("invalid packet description")
)

type Targets struct {
	Path    string `json:"path" validate:"min=1"`
	Exclude string `json:"exclude,omitempty" validate:"omitempty,min=1"`
}

type Packet struct {
	Name    string               `json:"name" validate:"min=1"`
	Version string               `json:"ver" validate:"min=1,pack_ver"`
	Targets []Targets            `json:"targets" validate:"min=1,dive"`
	Packets []PackageDescription `json:"packets,omitempty" validate:"omitempty,dive"`
}

func ParsePacket(data io.Reader) (*Packet, error) {
	var pack Packet
	decoder := json.NewDecoder(data)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pack); err != nil {
		return nil, ErrInvalidPacketDescriptionFormat
	}

	if err := validate.Validator().Struct(pack); err != nil {
		return nil, ErrInvalidPacketDescription
	}

	return &pack, nil
}
