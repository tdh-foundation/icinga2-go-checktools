package icinga2_go_checktools

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os/user"
	"regexp"
	"time"
)

//noinspection GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst
const (
	// icinga2 Constant

	// Ok Status
	OkMsg  = "OK"
	OkExit = 0

	// Warning Status
	WarMsg  = "WARNING"
	WarExit = 1

	// Critial Status
	CriMsg  = "CRITICAL"
	CriExit = 2

	// Unknown Status
	UnkMsg  = "UNKNOWN"
	UnkExit = 3
)

//
type SSHTools struct {
	sshClient *ssh.Client
	Stdout    string
	Msg       string
}

// NewSSHTools establish SSH connection to a server and return SSHTools object or nil if an error occurs
func NewSSHTools(host string, username string, password string, identity string, port int) (*SSHTools, error) {
	var signer ssh.Signer

	var ict = SSHTools{}

	// replacing tilde char by real home directory
	home, _ := user.Current()
	re := regexp.MustCompile(`^~(.*)$`)
	identity = re.ReplaceAllString(identity, home.HomeDir+"${1}")

	// Reading and parsing identity file (private key)
	key, err := ioutil.ReadFile(identity)
	if err == nil {
		// Create the Signer for this private key.
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			signer = nil
		}
	} else {
		signer = nil
	}

	// Configure authentication methods
	var auths []ssh.AuthMethod
	if signer != nil {
		auths = append(auths, ssh.PublicKeys(signer))
	}
	if password != "" {
		auths = append(auths, ssh.Password(password))
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // No key validation in known host
	}

	// Connecting SSH Server (Backup Exec Server)
	ict.sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("NewSSHTools, error establishing SSH connection: %s", err)
	}
	return &ict, err
}

// func SendSSH
// Send command to remote SSH server and save Output in Stdout
func (ict *SSHTools) SendSSH(command string) error {

	session, err := ict.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("SendSSH, error creating session on SSH server: %s", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("SendSSH, unable to setup stdout for session: %s", err)
	}

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("SendSSH, unable to setup stdout for session: %s", err)
	}

	buff := make([]byte, 1024)
	for {
		n, err := stdout.Read(buff)
		if err != nil && err != io.EOF {
			return fmt.Errorf("SendSSH, error reading output buffer: %s", err)
		}
		if n > 0 {
			ict.Stdout += string(buff[:n])
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}

// SendSSHhasPTY send multiple command to remote SSH session simulating a "human" ssh connection
// SendSSH must be preferred as SendSSHhasPTY, if prompt is misconfigured function will block
func (ict *SSHTools) SendSSHhasPTY(commands []string, prompt string) error {

	rePrompt := regexp.MustCompile(prompt)

	session, err := ict.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("SendSSHhasPTY, error creating session on SSH server: %s", err)
	}

	if err := session.RequestPty("xterm", 80, 40, nil); err != nil {
		session.Close()
		return fmt.Errorf("SendSSHhasPTY, request for pseudo terminal failed: %s", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("SendSSHhasPTY, unable to setup stdin for session: %s", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("SendSSHhasPTY, unable to setup stdout for session: %s", err)
	}

	err = session.Shell()
	ict.Stdout = ""
	for _, c := range commands {
		buff := make([]byte, 1024)
		stdin.Write([]byte(c))

		time.Sleep(1 * time.Second)

		for {
			n, err := stdout.Read(buff)
			if err != nil && err != io.EOF {
				return fmt.Errorf("SendSSHhasPTY, error reading output buffer: %s", err)
			}
			if n > 0 {
				ict.Stdout += string(buff[:n])
			}
			if rePrompt.MatchString(ict.Stdout) || err == io.EOF {
				break
			}
		}
	}
	return nil
}

// func CleanStdOutBuffer
// Remove consecutive blank char and newline of Stdout
func (ict *SSHTools) CleanStdOutBuffer() {
	reBlank := regexp.MustCompile(`(?m)[\s]{2,}`)
	reNewLine := regexp.MustCompile(`(?m)[\r|\n]+`)

	ict.Stdout = reBlank.ReplaceAllString(ict.Stdout, " ")
	ict.Stdout = reNewLine.ReplaceAllString(ict.Stdout, "")
}
