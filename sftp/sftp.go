package sftp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/sftp"
	"go.elastic.co/apm"
	"golang.org/x/crypto/ssh"
)

// SecureFileTransfer represent function exist in the SFTP BRIMO's class
type SecureFileTransfer interface {
	SendLocalFileToRemote(ctx context.Context, localpath, remotepath string) (int, error)
	SendLocalFileToRemoteWithDelete(ctx context.Context, localpath, remotepath string) (int, error)
	GetConnection(ctx context.Context) *ssh.Client
}

// SftpBrimo represet object of sftpClient
type SftpBrimo struct {
	sftpClient *sftp.Client
	Conn       *ssh.Client
}

// New implement SecureFileTransfer interface by returning object SftpBrimo,
// reconnection of SFTP for idle connection included in this constructor
func New(user, pass, host, port string) (SecureFileTransfer, error) {
	// Create Connection sftp
	conf := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", host+":"+port, conf)
	if err != nil {
		fmt.Println("dial tcp ssh fail : ", err)
		return nil, err
	}

	sftpClientNew, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Println("creation of object sftp client fail : ", err)
		return nil, err
	}

	sftpCurrent := &SftpBrimo{
		sftpClient: sftpClientNew,
		Conn:       conn,
	}

	go handleReconnect(sftpCurrent, user, pass, host, port)

	return sftpCurrent, nil
}

// SendLocalFileToRemote return length of byte written in remote FTP server
// this function move file from local path to FTP server by remote path
func (s *SftpBrimo) SendLocalFileToRemote(ctx context.Context, localpath, remotepath string) (int, error) {
	apmSpan, _ := apm.StartSpan(ctx, "Send Local File to Remote", "SFTP")
	defer apmSpan.End()

	written, err := s.sendFile(remotepath, localpath)
	if err != nil {
		return 0, err
	}

	return written, nil
}

// SendLocalFileToRemoteWithDelete move file from local to server,
// plus delete file after file sucessfully sent to FTP server
func (s *SftpBrimo) SendLocalFileToRemoteWithDelete(ctx context.Context, localpath, remotepath string) (int, error) {
	apmSpan, _ := apm.StartSpan(ctx, "Send&Delete Local File to Remote", "SFTP")
	defer apmSpan.End()

	written, err := s.sendFile(remotepath, localpath)
	if err != nil {
		return 0, err
	}

	// Delete local file
	err = os.Remove(localpath)
	if err != nil {
		fmt.Println("error on deletion local file : ", err)
		return 0, err
	}

	return written, nil
}

func (s *SftpBrimo) sendFile(remotepath, localpath string) (int, error) {
	remoteFile, err := s.sftpClient.Create(remotepath)
	if err != nil {
		fmt.Println("error on creating pipeline to remote hhost : ", err)
		return 0, err
	}
	defer remoteFile.Close()

	// Open the local file to be sent
	localFile, err := os.Open(localpath)
	if err != nil {
		fmt.Println("error on open local file : ", err)
		return 0, err
	}
	defer localFile.Close()

	// Read the local file and write it to the remote file
	bytes, err := ioutil.ReadAll(localFile)
	if err != nil {
		fmt.Println("error on read local file : ", err)
		return 0, err
	}

	written, err := remoteFile.Write(bytes)
	if err != nil {
		fmt.Println("error on write to remot file : ", err)
		return 0, err
	}

	return written, nil
}

// GetConnection return pointer of ssh client
func (s *SftpBrimo) GetConnection(ctx context.Context) *ssh.Client {
	return s.Conn
}

func handleReconnect(sftpCurrent *SftpBrimo, user, pass, host, port string) {
	closed := make(chan string)

	go func() {
		closed <- sftpCurrent.Conn.Wait().Error()
		fmt.Println("closed message : ", closed)
	}()

	select {
	case <-closed:
		// Create Connection sftp
		conf := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		fmt.Println("SFTP reconnection attempt ... ")

		conn, err := ssh.Dial("tcp", host+":"+port, conf)
		if err != nil {
			fmt.Println("error on dial to server : ", err)
			return
		}

		sftpClientNew, err := sftp.NewClient(conn)
		if err != nil {
			fmt.Println("error on creating object sftp clien : ", err)
			return
		}

		sftpCurrent.Conn = conn
		sftpCurrent.sftpClient = sftpClientNew

		fmt.Println("SFTP reconnection success on : ", time.Now().String())

		handleReconnect(sftpCurrent, user, pass, host, port)
	}

}
