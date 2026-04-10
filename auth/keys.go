package auth

import (
	"github.com/mcdonaldseanp/charlie/airer"
)

// UnlockKeys unlocks each key type in keyTypes in order. If updateTTY is true,
// gpg-agent's tty is updated to the current terminal before unlocking. If
// useYubikey is true, the YubiKey is mounted first using hwID.
func UnlockKeys(keyTypes []string, updateTTY bool, useYubikey bool, hwID string, sshKey string) error {
	if updateTTY {
		if err := UpdateGPGTTY(); err != nil {
			return err
		}
	}
	if useYubikey {
		err := MountYubikey(hwID)
		if err != nil {
			airr, ok := err.(*airer.Airer)
			if !ok || airr.Kind != airer.CompletedError {
				return err
			}
		}
	}
	for _, keyType := range keyTypes {
		var err error
		switch keyType {
		case "gpg":
			err = UnlockGPGKey()
		case "ssh":
			err = UnlockSSHKey(sshKey)
		default:
			err = &airer.Airer{
				Kind:    airer.InvalidInput,
				Message: "Unknown key type '" + keyType + "'. KEY TYPES should be one or more of 'gpg', 'ssh'",
				Origin:  nil,
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
