module gitlab.com/timeterm/timeterm/nats-manager

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ../backend/pkg/natspb
	gitlab.com/timeterm/timeterm/nats-manager/sdk => ./sdk
	gitlab.com/timeterm/timeterm/proto/go => ../proto/go
)

require (
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.6.8 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/nats-io/jsm.go v0.0.20-0.20201118085313-543c65ba42cd
	github.com/nats-io/jwt/v2 v2.0.0-20201030222427-057ba30017be
	github.com/nats-io/nats.go v1.10.1-0.20201111151633-9e1f4a0d80d8
	github.com/nats-io/nkeys v0.2.0
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/stretchr/testify v1.6.1
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/nats-manager/sdk v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201128102736-5b0f11963b4c
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/text v0.3.4 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gotest.tools/v3 v3.0.2 // indirect
)
