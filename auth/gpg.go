package auth

import (
	"github.com/mcdonaldseanp/charlie/airer"
	"github.com/mcdonaldseanp/lookout/localexec"
)

func CheckYubikeyLocked() error {
	_, _, err := localexec.ExecReadOutput("sh", "-c", "echo test | gpg --batch --pinentry-mode error --clearsign")
	if err != nil {
		return &airer.Airer{
			Kind:    airer.AuthError,
			Message: "YubiKey signing key is locked. Mount YubiKey with 'charlie mount yubikey' then unlock with: echo test | gpg --clearsign",
			Origin:  err,
		}
	}
	return nil
}
