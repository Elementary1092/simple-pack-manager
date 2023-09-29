package validate

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate
var once sync.Once

var (
    regexRemoteVersion = `^((<=)|(>=))?((0|[1-9])\d*)\.((0|[1-9])\d*)$`
    regexPacketVersion = `^((0|[1-9])\d*)\.((0|[1-9])\d*)$`
)

var (
    remoteVersionMatcher = regexp.MustCompile(regexRemoteVersion)
    packetVersionMatcher = regexp.MustCompile(regexPacketVersion)
)

func init() {
    once.Do(initValidator)
}

func initValidator() {
    v = validator.New()
    v.RegisterValidation("remote_ver", validateRemoteVersion)
    v.RegisterValidation("pack_ver", validatePacketVersion)
}

func Validator() *validator.Validate {
    return v
}

func validateRemoteVersion(fl validator.FieldLevel) bool {
    if fl.Field().Kind() != reflect.String {
        return false
    }

    str := fl.Field().String()
    if str == "" {
        return true
    }

    return remoteVersionMatcher.MatchString(str)
}

func validatePacketVersion(fl validator.FieldLevel) bool {
    if fl.Field().Kind() != reflect.String {
        return false
    }

    str := fl.Field().String()
    if str == "" {
        return true
    }

    return packetVersionMatcher.MatchString(str)
}

