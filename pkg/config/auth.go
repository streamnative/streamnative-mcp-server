package config

import (
	"path/filepath"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth/store"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
)

const (
	ServiceName  = "StreamNativeMCP"
	KeychainName = "snmcp"
)

type AuthOptions struct {
	BackendOverride string
	storage         Storage

	// AuthOptions is a facade for the token store
	// note: call Complete before using the token store methods
	store.Store
}

func NewDefaultAuthOptions() AuthOptions {
	return AuthOptions{}
}

func (o *AuthOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.BackendOverride, "keyring-backend", "",
		"If present, the backend to use")
}

func (o *AuthOptions) Complete(storage Storage) error {
	o.storage = storage
	kr, err := o.makeKeyring()
	if err != nil {
		return err
	}
	o.Store, err = store.NewKeyringStore(kr)
	if err != nil {
		return err
	}
	return nil
}

func (o *AuthOptions) makeKeyring() (keyring.Keyring, error) {
	var backends []keyring.BackendType
	if o.BackendOverride != "" {
		backends = append(backends, keyring.BackendType(o.BackendOverride))
	}

	return keyring.Open(keyring.Config{
		ServiceName:              ServiceName,
		KeychainName:             KeychainName,
		KeychainTrustApplication: true,
		AllowedBackends:          backends,
		FileDir:                  filepath.Join(o.storage.GetConfigDirectory(), "credentials"),
		FilePasswordFunc:         keyringPrompt,
	})
}

func keyringPrompt(_ string) (string, error) {
	return "", nil
}
