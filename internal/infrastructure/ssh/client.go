package ssh

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client handles SSH connections
type Client struct{}

// NewClient creates a new SSH client
func NewClient() *Client {
	return &Client{}
}

// Connect establishes an SSH connection with the given credentials
func (c *Client) Connect(ip, username, port, password string) (*ssh.Client, error) {
	// Configure SSH client
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Not recommended for production
		Timeout:         5 * time.Second,
	}

	// Format connection string
	addr := fmt.Sprintf("%s:%s", ip, port)

	// Establish connection
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	return client, nil
}

// RunCommand executes a command on an established SSH session
func RunCommand(client *ssh.Client, command string) (string, error) {
	// Create a new session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture output
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	// Run the command
	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("failed to run command '%s': %w", command, err)
	}

	return stdoutBuf.String(), nil
}
