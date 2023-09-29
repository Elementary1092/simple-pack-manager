package parser

import (
	"bytes"
	"errors"
	"testing"
)

func TestParsePackage_EmptyDescription(t *testing.T) {
	data := ""

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePackage(reader); err == nil && !errors.Is(err, ErrInvalidPackageDescription) {
		t.Fail()
	}
}

func TestParsePackage_NoPackages(t *testing.T) {
	data := `{
        "packages": []
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePackage(reader); err == nil || !errors.Is(err, ErrInvalidPackageDescription) {
		t.Fail()
	}
}

func TestParsePackage_WithoutVer(t *testing.T) {
	data := `{
        "packages": [
            {"name": "some_pack"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	res, err := ParsePackage(reader)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Packages) != 1 || res.Packages[0].Name != "some_pack" {
		t.Fatal("invalid decoding")
	}
}

func TestParsePackage_WithoutName(t *testing.T) {
	data := `{
        "packages": [
            {}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePackage(reader); err == nil || !errors.Is(err, ErrInvalidPackageDescription) {
		t.Fatal(err)
	}
}

func TestParsePackage_WithUnknownField(t *testing.T) {
	data := `{
        "packages": [
            {"name": "some_pack"}
        ],
        "cache": "no"
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePackage(reader); err == nil || !errors.Is(err, ErrInvalidPackageDescriptionFormat) {
		t.Fail()
	}
}

func TestParsePackage_Valid(t *testing.T) {
	data := `{
        "packages": [
        {"name": "some_pack", "ver": "1.0"},
        {"name": "pack"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	res, err := ParsePackage(reader)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Fatal("unexpected nil result")
	}
}
