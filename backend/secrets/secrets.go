package secrets

import (
	"github.com/hashicorp/vault/api"
	"github.com/go-logr/logr"
)

type Logger struct {
	log logr.Logger
}

func  (l *Logger) createClient() error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		l.log.Error(err, "could not create client")
	}

	logicalClient := client.Logical()
	logicalClient.Read()
}