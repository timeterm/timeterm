package secrets

import (
	"github.com/hashicorp/vault/api"
)

func createClient() error {
	client, err := api.NewClient()
	if err != nil {
		je moeder
	}

	logicalClient := client.Logical()
}