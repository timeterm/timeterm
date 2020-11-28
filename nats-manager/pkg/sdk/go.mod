module gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ./../../../backend/pkg/natspb
	gitlab.com/timeterm/timeterm/proto/go => ./../../../proto/go
)

require (
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.8 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/nats-io/jwt v1.2.0 // indirect
	github.com/nats-io/jwt/v2 v2.0.0-20201030222427-057ba30017be
	github.com/nats-io/nats.go v1.10.1-0.20201013114232-5a33ce07522f
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201128102736-5b0f11963b4c
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
