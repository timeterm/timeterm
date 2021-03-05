module gitlab.com/timeterm/timeterm/backend

go 1.16

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ./pkg/natspb
	gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk => ../nats-manager/pkg/sdk
	gitlab.com/timeterm/timeterm/proto/go => ../proto/go
)

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/evanphx/json-patch/v5 v5.1.0
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.1.2
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/nats-io/nats.go v1.10.1-0.20201013114232-5a33ce07522f
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pquerna/cachecontrol v0.0.0-20200921180117-858c6e7e6b7e // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.6.1
	github.com/valyala/fasttemplate v1.2.1 // indirect
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk v0.0.0-20201128102736-5b0f11963b4c
	gitlab.com/timeterm/timeterm/nats-manager/sdk v0.0.0-20201128103324-08c9faae8dc3
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201128102736-5b0f11963b4c
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.25.0
)
