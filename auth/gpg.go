package auth

import (
	"strings"

	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/lookout/localexec"
)

func CheckGPGKeyLocked() error {
	_, _, err := localexec.ExecReadOutput(
		strings.NewReader("test\n"),
		"gpg", "--batch", "--pinentry-mode", "error", "--clearsign",
	)
	if err != nil {
		return &airer.Airer{
			Kind:    airer.AuthError,
			Message: "YubiKey signing key is locked. Mount YubiKey with 'charlie mount yubikey' then unlock with: charlie unlock keys",
			Origin:  err,
		}
	}
	return nil
}

func UnlockGPGKey() error {
	return localexec.ExecAsShell(strings.NewReader("test\n"), "gpg", "--clearsign")
}
