package pmssh

import (
	"errors"
	"testing"
)

func TestVerifyConnData_AllDataIPv4Host(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = "password"
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestVerifyConnData_AllDataIPv6Host(t *testing.T) {
	sshHost = "0:0:0:0:0:0:0:0"
	sshPort = "22"
	sshUser = "user"
	sshPass = "password"
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestVerifyConnData_NoHost(t *testing.T) {
	sshHost = ""
	sshPort = "22"
	sshUser = "user"
	sshPass = "password"
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidHost) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidHost, err)
	}
}

func TestVerifyConnData_NoPort(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = ""
	sshUser = "user"
	sshPass = "password"
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidPort, err)
	}
}

func TestVerifyConnData_NoUser(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = ""
	sshPass = "password"
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidUser, err)
	}
}

func TestVerifyConnData_NoPassword(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = ""
	key = []byte("some valid ssh key")

	if err := verifyConnData(); err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestVerifyConnData_NoKey(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = "password"
	key = []byte("")

	if err := verifyConnData(); err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestVerifyConnData_NoAuthDataEmptyKey(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = ""
	key = []byte("")

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrNoAuthData) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrNoAuthData, err)
	}
}

func TestVerifyConnData_NoAuthDataNilKey(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = ""
	key = nil

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrNoAuthData) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrNoAuthData, err)
	}
}

func TestVerifyConnData_InvalidHost(t *testing.T) {
	sshHost = "invalidhost"
	sshPort = "22"
	sshUser = "user"
	sshPass = "password"
	key = nil

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidHost) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidHost, err)
	}
}

func TestVerifyConnData_InvalidPort(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "port"
	sshUser = "user"
	sshPass = "password"
	key = nil

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidPort, err)
	}
}

func TestVerifyConnData_InvalidUserIllegalCharacter(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user\n"
	sshPass = "password"
	key = nil

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidUser) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidUser, err)
	}
}

func TestVerifyConnData_InvalidPassword(t *testing.T) {
	sshHost = "127.0.0.1"
	sshPort = "22"
	sshUser = "user"
	sshPass = "password\xff"
	key = nil

	if err := verifyConnData(); err == nil || !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("Unexpected error: expected='%v'; got='%v'", ErrInvalidPassword, err)
	}
}
