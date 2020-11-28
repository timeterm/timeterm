module gitlab.com/timeterm/timeterm/nats-manager

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ../backend/pkg/natspb
	gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch => ./pkg/jwtpatch
	gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk => ./pkg/sdk
	gitlab.com/timeterm/timeterm/proto/go => ../proto/go
)

require (
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/nats-io/jsm.go v0.0.20-0.20201118085313-543c65ba42cd
	github.com/nats-io/jwt/v2 v2.0.0-20201030222427-057ba30017be
	github.com/nats-io/nats.go v1.10.1-0.20201111151633-9e1f4a0d80d8
	github.com/nats-io/nkeys v0.2.0
	github.com/stretchr/testify v1.6.1
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch v0.0.0-00010101000000-000000000000
	gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201128102736-5b0f11963b4c
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
)
