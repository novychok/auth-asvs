package vault

import (
	"net/http"

	vault "github.com/hashicorp/vault/api"
)

type Config struct {
	Address          string `mapstructure:"VAULT_ADDRESS"`
	Token            string `mapstructure:"VAULT_TOKEN"`
	PartnerKey       string `mapstructure:"VAULT_PARTNER_KEY"`
	TransitPath      string `mapstructure:"VAULT_TRANSIT_PATH"`
	PartnersDKIMPath string `mapstructure:"VAULT_PARTNERS_DKIM_PATH"`
}

func New(cfg *Config, httpClient *http.Client) (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.HttpClient = httpClient

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	err = client.SetAddress(cfg.Address)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Token)

	return client, nil
}
