package validate

import "testing"

type remoteData struct {
    S string `validate:"remote_ver"`
}

func TestRemoteVersion_Empty(t *testing.T) {
    data := remoteData{}

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestRemoteVersion_OnlyVersion(t *testing.T) {
    data := remoteData{"01.01"}

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestRemoteVersion_LessEqualVersion(t *testing.T) {
    data := remoteData{"<=01.01"}

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestRemoteVersion_GreaterEqualVersion(t *testing.T) {
    data := remoteData{">=01.01"}        

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestRemoteVersion_NoMinorVersion(t *testing.T) {
    data := remoteData{"01."}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}

func TestRemoteVersion_NoMajorVersion(t *testing.T) {
    data := remoteData{".1"}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}

func TestRemoteVersion_InvalidCharacters(t *testing.T) {
    data := remoteData{"data"}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}

type packetData struct {
    S string `validate:"pack_ver"`
}

func TestPacketVersion_Empty(t *testing.T) {
    data := remoteData{}

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestPacketVersion_OnlyVersion(t *testing.T) {
    data := remoteData{"01.01"}

    if err := Validator().Struct(data); err != nil {
        t.Fail()
    }
}

func TestPacketVersion_NoMinorVersion(t *testing.T) {
    data := packetData{"01."}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}

func TestPacketVersion_NoMajorVersion(t *testing.T) {
    data := packetData{".1"}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}

func TestPacketVersion_InvalidCharacters(t *testing.T) {
    data := packetData{"data"}

    if err := Validator().Struct(data); err == nil {
        t.Fail()
    }
}
