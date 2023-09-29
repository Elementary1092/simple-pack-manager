package parser

import (
	"bytes"
	"errors"
	"testing"
)

func TestParsePacket_EmptyDescription(t *testing.T) {
	data := ""

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePacket(reader); err == nil && !errors.Is(err, ErrInvalidPacketDescription) {
		t.Fail()
	}
}

func TestParsePacket_NoPackets(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePacket(reader); err != nil {
		t.Fatal(err)
	}
}

func TestParsePacket_WithoutTargets(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "packets": [
        {"name": "some_pack", "ver": "0.1"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	_, err := ParsePacket(reader)
	if err == nil || errors.Is(err, ErrInvalidPacketDescriptionFormat) {
		t.Fatal(err)
	}
}

func TestParsePacket_WithoutTargetsExclude(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./"}
        ],
        "packets": [
        {"name": "some_pack", "ver": "0.1"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	_, err := ParsePacket(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsePacket_WithoutPacketVer(t *testing.T) {
	data := `{
        "name": "a",
        "targets": [
            {"path": "./", "exclude":"*"}
        ],
        "packets": [
        {"name": "some_pack", "ver": "0.1"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	_, err := ParsePacket(reader)
	if err == nil {
		t.Fatal(err)
	}
}

func TestParsePacket_WithRelativeDependencyVer(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"}
        ],
        "packets": [
        {"name": "some_pack", "ver": ">=0.1"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	_, err := ParsePacket(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsePacket_WithoutDependencyVer(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"}
        ],
        "packets": [
            {"name": "some_pack"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	_, err := ParsePacket(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsePacket_WithoutName(t *testing.T) {
	data := `{
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"}
        ],
        "packets": [
            {"name": "some_pack"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePacket(reader); err == nil || !errors.Is(err, ErrInvalidPacketDescription) {
		t.Fail()
	}
}

func TestParsePacket_WithUnknownField(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"}
        ],
        "packets": [
            {"name": "some_pack"}
        ],
        "cache": "no"
    }`

	reader := bytes.NewReader([]byte(data))

	if _, err := ParsePacket(reader); err == nil || !errors.Is(err, ErrInvalidPacketDescriptionFormat) {
		t.Fail()
	}
}

func TestParsePacket_Valid(t *testing.T) {
	data := `{
        "name": "a",
        "ver": "1.0",
        "targets": [
            {"path": "./", "exclude":"*"},
            {"path": "./subdir", "exclude":"*.txt"}
        ],
        "packets": [
        {"name": "some_pack", "ver": "1.0"},
        {"name": "pack"}
        ]
    }`

	reader := bytes.NewReader([]byte(data))

	res, err := ParsePacket(reader)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Fatal("unexpected nil result")
	}
}
