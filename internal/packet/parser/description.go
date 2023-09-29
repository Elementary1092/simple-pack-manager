package parser

import (
	"encoding/json"
	"errors"
	"io"

	validate "github.com/Elementary1092/pm/internal/packet/validator"
)

var (
	ErrInvalidPackageDescriptionFormat = errors.New("package description format is invalid")
	ErrInvalidPackageDescription       = errors.New("invalid package description")
)

type PackageDescription struct {
	Name    string `json:"name" validate:"required"`
	Version string `json:"ver,omitempty" validate:"remote_ver"`
}

type Package struct {
	Packages []PackageDescription `json:"packages" validate:"required,min=1,dive"`
}

func ParsePackage(data io.Reader) (*Package, error) {
	var pack Package
	decoder := json.NewDecoder(data)
	decoder.DisallowUnknownFields()

	// Could have wrapped json error into a struct implementing error interface
	if err := decoder.Decode(&pack); err != nil {
		return nil, ErrInvalidPackageDescriptionFormat
	}

	if err := validate.Validator().Struct(pack); err != nil {
		return nil, ErrInvalidPackageDescription
	}

	return &pack, nil
}
