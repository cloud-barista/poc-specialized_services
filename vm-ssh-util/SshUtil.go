package sshutil

import (
	"fmt"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
)

func SshRunWithKey(vmIP string, userName string, privateKey string, cmd string) (string, error) {

	// Remote VM SSH connection information
	serverEndpoint := fmt.Sprintf("%s:22", vmIP)
	sshInfo := sshrun.SSHInfo{
		ServerPort: serverEndpoint,
		UserName:   userName,
		PrivateKey: []byte(privateKey),
	}

	// Run remote SSH commnad
	if result, err := sshrun.SSHRun(sshInfo, cmd); err != nil {
		return "Error : ", err
	} else {
		return result, nil
	}
}


func SshCopyWithKey(vmIP string, userName string, privateKey string, sourceFilePath string, targetFilePath string) (error) {

	// Remote VM SSH connection information
	serverEndpoint := fmt.Sprintf("%s:22", vmIP)
	sshInfo := sshrun.SSHInfo{
		ServerPort: serverEndpoint,
		UserName:   userName,
		PrivateKey: []byte(privateKey),
	}

	// File info to Copy
	sourceFile := sourceFilePath
	targetFile := targetFilePath

	// Run remote SSH commnad
	if err := sshrun.SSHCopy(sshInfo, sourceFile, targetFile); err != nil {
		return err
	} else {
		return nil
	}
}


func SshRunWithKeyPath(vmIP string, userName string, privateKeyPath string, cmd string) (string, error) {

	// Remote VM SSH connection information
	serverEndpoint := fmt.Sprintf("%s:22", vmIP)
	sshKeypathInfo := sshrun.SSHKeyPathInfo{
		ServerPort: serverEndpoint,
		UserName:   userName,
		KeyPath: privateKeyPath,
	}

	// Run remote SSH commnad
	if result, err := sshrun.SSHRunByKeyPath(sshKeypathInfo, cmd); err != nil {
		return "Error : ", err
	} else {
		return result, nil
	}
}


func SshCopyWithKeyPath(vmIP string, userName string, privateKeyPath string, sourceFilePath string, targetFilePath string) (error) {

	// Remote VM SSH connection information
	serverEndpoint := fmt.Sprintf("%s:22", vmIP)
	sshKeypathInfo := sshrun.SSHKeyPathInfo{
		ServerPort: serverEndpoint,
		UserName:   userName,
		KeyPath: privateKeyPath,
	}

	// File info to Copy
	sourceFile := sourceFilePath
	targetFile := targetFilePath

	// Run remote SSH Copy commnad
	if err := sshrun.SSHCopyByKeyPath(sshKeypathInfo, sourceFile, targetFile); err != nil {
		return err
	} else {
		return nil
	}
}
