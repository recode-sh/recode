package entities

import "errors"

var (
	ErrRecodeNotInstalled       = errors.New("ErrRecodeNotInstalled")
	ErrUninstallExistingDevEnvs = errors.New("ErrUninstallExistingDevEnvs")
)
