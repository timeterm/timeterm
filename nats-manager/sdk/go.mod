module gitlab.com/timeterm/timeterm/nats-manager/sdk

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ../../backend/pkg/natspb
	gitlab.com/timeterm/timeterm/proto/go => ../../proto/go
)

require (
	github.com/nats-io/jwt v1.2.0 // indirect
	github.com/nats-io/nats.go v1.10.1-0.20201013114232-5a33ce07522f
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201110122546-fd086d39b6a5
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201110122546-fd086d39b6a5
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
)
