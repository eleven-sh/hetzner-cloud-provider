package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/eleven-sh/eleven/entities"
	"golang.org/x/crypto/ssh"
)

const (
	ServerSSHPort  = 22
	ServerRootUser = "root"
)

type RawInitServerScriptResults struct {
	ExitCode      string `json:"exit_code"`
	SSHHostKeys   string `json:"ssh_host_keys"`
	CloudInitLogs string `json:"cloud_init_logs"`
}

type InitServerScriptResults struct {
	ExitCode    string                   `json:"exit_code"`
	SSHHostKeys []entities.EnvSSHHostKey `json:"ssh_host_keys"`
}

func LookupServerInitScriptResults(
	serverPublicIPAddress string,
	serverSSHPort string,
	serverLoginUser string,
	sshPrivateKeyContent string,
) (returnedInitScriptResults *InitServerScriptResults, returnedError error) {

	pollTimeoutChan := time.After(4 * time.Minute)
	pollSleepDuration := 4 * time.Second

	for {
		select {
		case <-pollTimeoutChan:
			cloudInitLogs, err := runCmdOnServerViaSSH(
				serverPublicIPAddress,
				serverSSHPort,
				serverLoginUser,
				sshPrivateKeyContent,
				"cat /var/log/cloud-init-output.log",
			)

			if err != nil {
				cloudInitLogs = fmt.Sprintf(
					"<cannot_retrieve_cloud_init_logs> (\"%v\")",
					err,
				)
			}

			returnedError = entities.ErrEnvCloudInitError{
				Logs: cloudInitLogs,
				ErrorMessage: fmt.Sprintf(
					"Server cloud init script timed out (\"%v\")",
					returnedError,
				),
			}
			return
		default:
			initScriptOutput, err := runCmdOnServerViaSSH(
				serverPublicIPAddress,
				serverSSHPort,
				serverLoginUser,
				sshPrivateKeyContent,
				"cat /tmp/eleven-init-results",
			)

			// Make sure timeout returns last error
			returnedError = err

			if err != nil {
				break // wait "pollSleepDuration" and retry until timeout
			}

			var initScriptResults *RawInitServerScriptResults
			err = json.Unmarshal([]byte(initScriptOutput), &initScriptResults)

			if err != nil {
				returnedError = entities.ErrEnvCloudInitError{
					Logs: initScriptOutput,
					ErrorMessage: fmt.Sprintf(
						"Server cloud init script exited with invalid JSON (\"%v\").",
						err,
					),
				}
				return
			}

			if initScriptResults.ExitCode != "0" {
				returnedError = entities.ErrEnvCloudInitError{
					Logs: initScriptResults.CloudInitLogs,
					ErrorMessage: fmt.Sprintf(
						"Server cloud init script exited with code \"%s\".",
						initScriptResults.ExitCode,
					),
				}
				return
			}

			parsedSSHHostKeys, err := entities.ParseSSHHostKeysForEnv(
				initScriptResults.SSHHostKeys,
			)

			if err != nil {
				returnedError = entities.ErrEnvCloudInitError{
					Logs: initScriptResults.SSHHostKeys,
					ErrorMessage: fmt.Sprintf(
						"Server cloud init script exited with invalid SSH host keys (\"%v\").",
						err,
					),
				}
				return
			}

			returnedInitScriptResults = &InitServerScriptResults{
				ExitCode:    initScriptResults.ExitCode,
				SSHHostKeys: parsedSSHHostKeys,
			}
			return
		} // <- end of select

		time.Sleep(pollSleepDuration)
	} // <- end of for
}

func WaitForSSHAvailableInServer(
	serverPublicIPAddress string,
	serverSSHPort string,
) (returnedError error) {

	pollTimeoutChan := time.After(4 * time.Minute)
	pollSleepDuration := 4 * time.Second

	SSHConnTimeout := 8 * time.Second

	for {
		select {
		case <-pollTimeoutChan:
			return
		default:
			conn, err := net.DialTimeout(
				"tcp",
				net.JoinHostPort(
					serverPublicIPAddress,
					serverSSHPort,
				),
				SSHConnTimeout,
			)

			// Make sure timeout returns last error
			returnedError = err

			if err != nil {
				break // wait "pollSleepDuration" and retry until timeout
			}

			conn.Close()
			return
		}

		time.Sleep(pollSleepDuration)
	}
}

func runCmdOnServerViaSSH(
	serverPublicIPAddress string,
	serverSSHPort string,
	loginUser string,
	privateKeyContent string,
	cmd string,
) (string, error) {

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyContent))

	if err != nil {
		return "", err
	}

	SSHConnTimeout := time.Second * 8

	config := &ssh.ClientConfig{
		User: loginUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         SSHConnTimeout,
	}

	client, err := ssh.Dial(
		"tcp",
		net.JoinHostPort(
			serverPublicIPAddress,
			serverSSHPort,
		),
		config,
	)

	if err != nil {
		return "", err
	}

	session, err := client.NewSession()

	if err != nil {
		return "", err
	}

	defer session.Close()

	var output bytes.Buffer
	session.Stdout = &output

	err = session.Run(cmd)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
