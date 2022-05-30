package entities

import (
	"reflect"
	"testing"
)

func TestParseValidSSHHostKeys(t *testing.T) {
	givenSSHHostKeysContent := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzd root@ip-10-0-0-179
	ssh-ed25519 AAAAC3NzaC1lZD root@ip-10-0-0-179
	ssh-rsa AAAAB3NzaC root@ip-10-0-0-179
`
	expectedSSHHostKeys := []DevEnvSSHHostKey{
		{
			Algorithm:   "ecdsa-sha2-nistp256",
			Fingerprint: "AAAAE2VjZHNhLXNoYTItbmlzd",
		},

		{
			Algorithm:   "ssh-ed25519",
			Fingerprint: "AAAAC3NzaC1lZD",
		},

		{
			Algorithm:   "ssh-rsa",
			Fingerprint: "AAAAB3NzaC",
		},
	}

	returnedSSHHostKeys, err := ParseSSHHostKeys(givenSSHHostKeysContent)

	if err != nil {
		t.Fatalf(
			"expected no error, got %s",
			err,
		)
	}

	if !reflect.DeepEqual(expectedSSHHostKeys, returnedSSHHostKeys) {
		t.Fatalf(
			"expected SSH host keys to equal '%+v', got '%+v'",
			expectedSSHHostKeys,
			returnedSSHHostKeys,
		)
	}
}

func TestParseInvalidSSHHostKeys(t *testing.T) {
	giventSSHHostKeysContent := "host_keys_content"
	returnedSSHHostKeys, err := ParseSSHHostKeys(giventSSHHostKeysContent)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	if returnedSSHHostKeys != nil {
		t.Fatalf(
			"expected no SSH host keys, got '%+v'",
			returnedSSHHostKeys,
		)
	}
}
