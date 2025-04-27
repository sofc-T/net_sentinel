package probe

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/reiver/go-telnet"
)

// SSHConfig holds SSH connection details
type SSHConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Timeout  time.Duration
}

// RunSSHCommand executes a command over SSH
func RunSSHCommand(config SSHConfig, command string) (string, error) {
	clientConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{ssh.Password(config.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only
		Timeout:         config.Timeout,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), clientConfig)
	if err != nil {
		return "", fmt.Errorf("SSH connection error: %v", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("SSH session error: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("SSH command failed: %v", err)
	}

	return string(output), nil
}

// RunTelnetCommand connects via Telnet and executes a command
func RunTelnetCommand(host string, command string) (string, error) {
	conn, err := telnet.DialTo(host)
	if err != nil {
		return "", fmt.Errorf("Telnet connection error: %v", err)
	}
	defer conn.Close()

	// Send command
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		return "", fmt.Errorf("Telnet write failed: %v", err)
	}

	// Read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("Telnet read failed: %v", err)
	}

	response := string(buf[:n])
	return response, nil
}


// Example Usage
func main() {
	// SSH Example
	sshConfig := SSHConfig{
		Host:     "192.168.1.1",
		Port:     "22",
		Username: "admin",
		Password: "password",
		Timeout:  5 * time.Second,
	}

	sshOutput, err := RunSSHCommand(sshConfig, "show interfaces")
	if err != nil {
		log.Fatalf("SSH Error: %v", err)
	}
	fmt.Println("SSH Output:\n", sshOutput)

	// Telnet Example
	telnetOutput, err := RunTelnetCommand("192.168.1.1:23", "show ip route")
	if err != nil {
		log.Fatalf("Telnet Error: %v", err)
	}
	fmt.Println("Telnet Output:\n", telnetOutput)
}
