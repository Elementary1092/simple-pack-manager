package pmssh

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	validate "github.com/Elementary1092/pm/internal/packet/validator"
	"github.com/pkg/sftp"
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
	ErrInvalidHost     = errors.New("invalid ssh host")
	ErrInvalidPort     = errors.New("invalid port")
	ErrInvalidUser     = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoAuthData      = errors.New("no auth data")

	ErrInvalidAuthData         = errors.New("invalid auth data")
	ErrConnectionFailure       = errors.New("unable to connect to the server")
	ErrNotConnected            = errors.New("not connected to the server")
	ErrFailedToOpenSource      = errors.New("failed to open source file")
	ErrFailedToUploadFile      = errors.New("failed to upload file")
	ErrFailedToOpenDestination = errors.New("failed to open destination file")
	ErrCannotReadDirectory     = errors.New("cannot read directory")
	ErrFailedToDownloadFile    = errors.New("failed to download file")
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

// Singleton defining ssh connection
var conn *ssh.Client

// This mutex is used to make possible returning an error from Connect method
// Also, prevents concurrent download and upload
var connMutex sync.Mutex

var fs *sftp.Client
var onceSftpClient sync.Once

func fsClient() *sftp.Client {
	// assuming that NewClient never fails
	onceSftpClient.Do(func() {
		fs, _ = sftp.NewClient(conn)
	})

	return fs
}

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

	if fs != nil {
		fs.Close()
	}

	if conn == nil {
		return nil
	}

	if err := conn.Close(); err != nil {
		return ErrConnectionFailure
	}

	conn = nil

	return nil
}

func Upload(ctx context.Context, dstFilePath string, fileFullName string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn == nil {
		return ErrNotConnected
	}

	return upload(ctx, dstFilePath, fileFullName)
}

func upload(ctx context.Context, dstPath string, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return ErrFailedToOpenSource
	}
	defer src.Close()

	client := fsClient()
	if client == nil {
		return ErrFailedToUploadFile
	}

	err = client.MkdirAll(filepath.Dir(dstPath))
	if err != nil {
		return ErrFailedToUploadFile
	}

	dst, err := client.Create(dstPath)
	if err != nil {
		return ErrFailedToUploadFile
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return ErrFailedToUploadFile
	}

	return nil
}

func Download(ctx context.Context, srcFullName string, dstFullName string) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn == nil {
		return ErrNotConnected
	}

	return download(ctx, srcFullName, dstFullName)
}

func download(ctx context.Context, srcFullName string, dstFullName string) error {
	dst, err := os.Open(dstFullName)
	if err != nil {
		return ErrFailedToOpenDestination
	}
	defer dst.Close()

	client := fsClient()
	if client == nil {
		return ErrFailedToOpenSource
	}

	srcStat, err := client.Lstat(srcFullName)
	if err != nil {
		return ErrFailedToOpenSource
	}

	if srcStat.IsDir() {
		return ErrCannotReadDirectory
	}

	src, err := client.Open(srcFullName)
	if err != nil {
		return ErrFailedToOpenSource
	}
	defer src.Close()

	if _, err := src.WriteTo(dst); err != nil {
		return ErrFailedToDownloadFile
	}

	return nil
}

func DoesFileExist(ctx context.Context, remoteFilePath string) bool {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn == nil {
		return false
	}

	return doesFileExist(ctx, remoteFilePath)
}

func doesFileExist(ctx context.Context, remoteFile string) bool {
	client := fsClient()
	if client == nil {
		return false
	}

	_, err := client.Lstat(remoteFile)

	return err != nil
}
