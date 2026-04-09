package auth

import (
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/lookout/localexec"
)

func CheckYubikeyLocked() error {
	_, _, err := localexec.ExecReadOutput(
		strings.NewReader("test\n"),
		"gpg", "--batch", "--pinentry-mode", "error", "--clearsign",
	)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.AuthError,
			Message: "YubiKey signing key is locked. Mount YubiKey with 'charlie mount yubikey' then unlock with: charlie unlock yubikey",
			Origin:  err,
		}
	}
	return nil
}

func UnlockYubikey(hw_id string) error {
	err := MountYubikey(hw_id)
	if err != nil {
		airr, ok := err.(*airer.Airer)
		if !ok || airr.Kind != airer.CompletedError {
			return err
		}
	}
	return localexec.ExecAsShell(strings.NewReader("test\n"), "gpg", "--clearsign")
}
