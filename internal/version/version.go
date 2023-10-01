package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type versionType int

const (
    Exact versionType = iota
    LessOrEqual
    GreaterOrEqual
    Invalid
    Less
    Greater
)

var (
    ErrInvalidVersionFormat = errors.New("invalid version format")
)

func Type(v string) versionType {
    if strings.HasPrefix(v, "<=") {
        return LessOrEqual
    }

    if strings.HasPrefix(v, ">=") {
        return GreaterOrEqual
    }

    return Exact
}

func Clean(v string) string {
    return strings.TrimLeft(v, "<>=")
}

func parseVersion(v string) (uint64, uint64, error) {
    v1Splitted := strings.Split(v, ".")
    if len(v1Splitted) != 2 {
        return 0, 0, fmt.Errorf("%v: %s", ErrInvalidVersionFormat, v)
    }
    
    major, err := strconv.ParseUint(v1Splitted[0], 10, 32)
    if err != nil {
        return 0, 0, fmt.Errorf("%v: %s", ErrInvalidVersionFormat, v)
    }

    minor, err := strconv.ParseUint(v1Splitted[1], 10, 32)
    if err != nil {
        return 0, 0, fmt.Errorf("%v: %s", ErrInvalidVersionFormat, v)
    }
    
    return major, minor, nil
}

func CompareVersions(v1 string, v2 string) (versionType, error) {
    v1Major, v1Minor, err := parseVersion(v1)
    if err != nil {
        return Invalid, err
    }

    v2Major, v2Minor, err := parseVersion(v2)
    if err != nil {
        return Invalid, err
    }

    if v1Major < v2Major {
        return Less, nil
    }

    if v1Major > v2Major {
        return Greater, nil
    }

    if v1Minor < v2Minor {
        return Less, nil
    }

    if v1Minor > v2Minor {
        return Greater, nil
    }

    return Exact, nil
}

