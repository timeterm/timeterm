module gitlab.com/timeterm/timeterm/nats-manager

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ../backend/pkg/natspb
	gitlab.com/timeterm/timeterm/nats-manager/sdk => ./sdk
	gitlab.com/timeterm/timeterm/proto/go => ../proto/go
)

require (
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/nats-io/jsm.go v0.0.19
	github.com/nats-io/jwt/v2 v2.0.0-20201006231922-e00ffcea7738
	github.com/nats-io/nats.go v1.10.1-0.20201013114232-5a33ce07522f
	github.com/nats-io/nkeys v0.2.0
	github.com/stretchr/testify v1.6.1
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201110122546-fd086d39b6a5
	gitlab.com/timeterm/timeterm/nats-manager/sdk v0.0.0-00010101000000-000000000000
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201110122546-fd086d39b6a5
	go.uber.org/zap v1.16.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)
