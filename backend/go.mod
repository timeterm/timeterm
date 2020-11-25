module gitlab.com/timeterm/timeterm/backend

go 1.15

replace (
	gitlab.com/timeterm/timeterm/backend/pkg/natspb => ./pkg/natspb
	gitlab.com/timeterm/timeterm/nats-manager/sdk => ../nats-manager/sdk
	gitlab.com/timeterm/timeterm/proto/go => ../proto/go
)

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/evanphx/json-patch/v5 v5.1.0
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.3.0
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/google/uuid v1.1.2
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
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
	gitlab.com/timeterm/timeterm/backend/pkg/natspb v0.0.0-20201116065525-ea9cced58338
	gitlab.com/timeterm/timeterm/nats-manager/sdk v0.0.0-00010101000000-000000000000
	gitlab.com/timeterm/timeterm/proto/go v0.0.0-20201116065525-ea9cced58338
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201116153603-4be66e5b6582
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20201116194326-cc9327a14d48 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
