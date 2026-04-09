package auth

import (
	"os"

	"github.com/mcdonaldseanp/lookout/localexec"
)

func UnlockSSHKey(ssh_key string) error {
	// ssh-keygen -Y sign requires a file path for the public key, so write it to a temp file
	key_file, err := os.CreateTemp("", "charlie_ssh_key")
	if err != nil {
		return err
	}
	key_filename := key_file.Name()
	defer os.Remove(key_filename)
	_, err = key_file.WriteString(ssh_key + "\n")
	key_file.Close()
	if err != nil {
		return err
	}

	// Create the file to be signed
	sign_file, err := os.CreateTemp("", "charlie_ssh_unlock")
	if err != nil {
		return err
	}
	sign_file.Close()
	sign_filename := sign_file.Name()
	defer os.Remove(sign_filename)
	defer os.Remove(sign_filename + ".sig")

	// Sign the file interactively — this prompts for the YubiKey PIN and caches the auth subkey
	return localexec.ExecAsShell(nil, "ssh-keygen", "-Y", "sign", "-f", key_filename, "-n", "test", sign_filename)
}
