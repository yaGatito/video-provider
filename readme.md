# video-provider

## Quick Start (from zero)

1. Clone repo and enter folder.
2. Create `.env` in project root:

```env
POSTGRES_USER=gato
POSTGRES_PASSWORD=root
```

3. Install required tools (Go, Docker, GNU Make, `yq`, plus Go CLIs below).
4. Run:

```bash
make bootstrap
make setup
```

5. Start services in separate terminals:

```bash
make run CONFIG=video
make run CONFIG=user
```

## Cross-platform setup

### Linux (Ubuntu/Debian example)

```bash
sudo apt update
sudo apt install -y make docker.io curl
sudo snap install yq

go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Note: ensure `$(go env GOPATH)/bin` is in `PATH`.

### Windows (PowerShell + Scoop)

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
scoop install make
scoop install docker
scoop install yq

go install github.com/golang/mock/mockgen@v1.6.0
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Make targets

```bash
make help
```

Main targets:
- `make bootstrap` - verify all required tools are installed.
- `make setup` - start both DB containers, init DBs, run migrations.
- `make db-status CONFIG=video|user` - print DB settings.
- `make db-up CONFIG=video|user` - start one DB container.
- `make db-down CONFIG=video|user` - stop/remove one DB container.
- `make migrate-up CONFIG=video|user` - run migrations for one service.
- `make run CONFIG=video|user` - run selected service.
- `make test PKG=<pkg> TEST=<regex>` - focused tests.
- `make tests` - full test suite.
