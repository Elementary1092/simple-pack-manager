package pmssh

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	validate "github.com/Elementary1092/pm/internal/packet/validator"
	"golang.org/x/crypto/ssh"
)

const (
    maxRetries = 3
)

var (
    //go:embed ssh-ip.conf
    sshHost string

    //go:embed ssh-port.conf
    sshPort string

    //go:embed ssh-user.conf
    sshUser string

    //go:embed ssh-password.conf
    sshPass string

    //go:embed ssh-key.pem
    key []byte
)

var (
    ErrInvalidHost = errors.New("invalid ssh host")
    ErrInvalidPort = errors.New("invalid port")
    ErrInvalidUser = errors.New("invalid username")
    ErrInvalidPassword = errors.New("invalid password")
    ErrNoAuthData = errors.New("no auth data")
    ErrInvalidAuthData = errors.New("invalid auth data")
    ErrConnectionFailure = errors.New("unable to connect to the server")
)

func verifyConnData() error {
    const whitespaces = "\n\t\r "
    sshHost = strings.TrimRight(sshHost, whitespaces)
    sshPort = strings.TrimRight(sshPort, whitespaces)
    sshUser = strings.TrimRight(sshUser, whitespaces)
    sshPass = strings.TrimRight(sshPass, whitespaces)

    if err := validate.Validator().Var(sshHost, "min=1,ipv4|ipv6"); err != nil {
        return ErrInvalidHost
    }

    if err := validate.Validator().Var(sshPort, "min=1,number"); err != nil {
        return ErrInvalidPort
    }

    if err := validate.Validator().Var(sshUser, "min=1,alphanum"); err != nil {
        return ErrInvalidUser
    }

    if err := validate.Validator().Var(sshPass, "printascii"); err != nil {
        return ErrInvalidPassword
    }

    if len(sshPass) == 0 && len(key) == 0 {
        return ErrNoAuthData
    }

    return nil
}

var conn ssh.Conn
// This mutex allows simplification of graceful shutdown
// Module will wait for all other tasks to be compeleted before closing connection.
var connMutex sync.Mutex

func createConnection(ctx context.Context) error {
    var authMethod ssh.AuthMethod
    if len(key) != 0 {
        signer, err := ssh.ParsePrivateKey(key)
        if err != nil {
            return ErrInvalidAuthData
        }

        authMethod = ssh.PublicKeys(signer)
    } else {
        authMethod = ssh.Password(sshPass)
    }

    cfg := ssh.ClientConfig{
        User: sshUser,
        Auth: []ssh.AuthMethod{
            ssh.RetryableAuthMethod(authMethod, maxRetries),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    address := net.JoinHostPort(sshHost, sshPort)

    sshConn, err := ssh.Dial("tcp", address, &cfg)
    if err != nil {
        fmt.Println("Failed to connect:", err)
        return ErrConnectionFailure
    }

    conn = sshConn

    return nil
}

func Connect(ctx context.Context) error {
    if err := verifyConnData(); err != nil {
        return err
    }
    connMutex.Lock()
    defer connMutex.Unlock()

    if conn != nil {
        return nil
    }

    return createConnection(ctx)
}

func Close(ctx context.Context) error {
    connMutex.Lock()
    defer connMutex.Unlock()

    if conn != nil {
        return nil
    }

    if err := conn.Close(); err != nil {
        return ErrConnectionFailure
    }

    return nil
}

