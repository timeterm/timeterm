module gitlab.com/timeterm/timeterm/nats-manager

go 1.15

replace gitlab.com/timeterm/timeterm/backend/pkg/natspb => ../backend/pkg/natspb

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/nats-io/jwt v1.2.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	github.com/nats-io/nats.go v1.10.0
	gitlab.com/timeterm/timeterm/backend v0.0.0-20201107204214-3a92b404c989
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
)
