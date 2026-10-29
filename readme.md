<div align="center">

# 📚 onelogin-go-sdk

**Go client for the OneLogin REST API: apps, users, roles, auth servers and more.**

A Glovo-maintained fork of OneLogin's Go SDK at v1.1.20. It is a client library only: it runs no service, and it is not the official OneLogin SDK.

[![Go](https://img.shields.io/badge/Go-1.12_(go.mod)-00ADD8?logo=go&logoColor=white)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/glovo/onelogin-go-sdk)](https://goreportcard.com/report/github.com/glovo/onelogin-go-sdk)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-95%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Tests](https://img.shields.io/badge/test_files-24-blue)](.github/workflows/go.yml)

</div>

---

## Table of contents

- [Overview](#overview)
- [Usage](#usage)
- [Build & local development](#build--local-development)
- [Testing](#testing)
- [Code quality](#code-quality)
- [CI/CD](#cicd)
- [Repository layout](#repository-layout)
- [Contributing & conventions](#contributing--conventions)
- [Ownership & support](#ownership--support)

---

## Overview

The SDK wraps the OneLogin API behind one `APIClient`. You give it a client ID, client secret and region, and it exposes one typed service per API: apps, users, roles, auth servers and so on. Each service sends authenticated HTTP calls to `https://api.<region>.onelogin.com` and returns typed models.

| | |
|---|---|
| Language / build | Go, module `github.com/glovo/onelogin-go-sdk` ([`go.mod`](go.mod) declares `go 1.12`) |
| Regions | `us`, `eu` ([`client_config.go`](pkg/client/client_config.go)); or set `Url` yourself |
| Default timeout | 5 seconds |
| Releases | Git tags `v*` (latest tag `v1.1.23`) create a GitHub release ([`release.yml`](.github/workflows/release.yml)) |
| Origin | Fork of OneLogin's Go SDK at v1.1.20 (first commit of this repo) |

### Supported APIs

Each API is a field on `client.Services` ([`client.go`](pkg/client/client.go)).

| API | Field | Package |
|---|---|---|
| [Apps](https://developers.onelogin.com/api-docs/2/apps/overview) (v2) | `AppsV2` | [`pkg/services/apps`](pkg/services/apps) |
| [App Rules](https://developers.onelogin.com/api-docs/2/app-rules/overview) | `AppRulesV2` | [`pkg/services/apps/app_rules`](pkg/services/apps/app_rules) |
| [Auth Servers](https://developers.onelogin.com/api-docs/2/api-authorization/overview) (v2) | `AuthServersV2` | [`pkg/services/auth_servers`](pkg/services/auth_servers) |
| Access token claims (auth servers, v2) | `AccessTokenClaimsV2` | [`auth_servers/access_token_claims`](pkg/services/auth_servers/access_token_claims) |
| Scopes (auth servers, v2) | `ScopesV2` | [`auth_servers/scopes`](pkg/services/auth_servers/scopes) |
| Client apps (auth servers, v2) | `ClientAppsV2` | [`auth_servers/client_apps`](pkg/services/auth_servers/client_apps) |
| [Roles](https://developers.onelogin.com/api-docs/2/roles/list-roles) (v1 service) | `RolesV1` | [`pkg/services/roles`](pkg/services/roles) |
| [Privileges](https://developers.onelogin.com/api-docs/1/privileges/list-privileges) (v1) | `PrivilegesV1` | [`pkg/services/privileges`](pkg/services/privileges) |
| [Session Login Tokens](https://developers.onelogin.com/api-docs/1/login-page/create-session-login-token) (v1) | `SessionLoginTokensV1` | [`pkg/services/session_login_tokens`](pkg/services/session_login_tokens) |
| [Smart Hooks](https://developers.onelogin.com/api-docs/2/smart-hooks/overview) (v1 service) | `SmartHooksV1` | [`pkg/services/smarthooks`](pkg/services/smarthooks) |
| Smart Hooks environment variables | `SmartHooksEnvVarsV1` | [`smarthooks/envs`](pkg/services/smarthooks/envs) |
| [User Mappings](https://developers.onelogin.com/api-docs/2/user-mappings/overview) (v2) | `UserMappingsV2` | [`pkg/services/user_mappings`](pkg/services/user_mappings) |
| [Users](https://developers.onelogin.com/api-docs/2/users/list-users) (v2) | `UsersV2` | [`pkg/services/users`](pkg/services/users) |

> [!NOTE]
> The old README labelled Roles and Smart Hooks as "(v2)". The code names those services `V1` (`RolesV1`, `SmartHooksV1`), so the table follows the code.

---

## Usage

Install and import (needs [Go](https://golang.org/doc/install)):

```bash
go get -u "github.com/glovo/onelogin-go-sdk/pkg/client"
```

```go
import "github.com/glovo/onelogin-go-sdk/pkg/client"
```

Start the client. `ClientID` and `ClientSecret` are required; `Region` is required unless you set `Url`:

```go
sdkClient, err := client.NewClient(&client.APIClientConfig{
  Timeout:      5,
  ClientID:     "your_onelogin_developer_client_id",
  ClientSecret: "your_onelogin_developer_client_secret",
  Region:       "us",
})
if err != nil {
  // handle error
}

app, err := sdkClient.Services.AppsV2.GetOne(12345)
if err != nil {
  // handle error
}
```

Models use pointer fields; build them with the helpers in [`pkg/oltypes`](pkg/oltypes) (`oltypes.String`, `oltypes.Int32`, `oltypes.Bool`, `oltypes.Int64`).

- **Query apps.** Query values are strings:
  ```go
  awsApps, err := sdkClient.Services.AppsV2.Query(&apps.AppsQuery{ConnectorID: "9", Limit: "10"}) // get 10 aws apps
  ```
- **Create an app.** The app is updated in place with the state OneLogin returns:
  ```go
  newApp := &apps.App{Name: oltypes.String("new app"), ConnectorID: oltypes.Int32(9)}
  err := sdkClient.Services.AppsV2.Create(newApp)
  if err != nil {
    return err
  }
  fmt.Println("my app", *newApp.ID)
  ```
- **Update an app:**
  ```go
  someApp, err := sdkClient.Services.AppsV2.GetOne(123)
  someApp.Name = oltypes.String("updated name")
  updated, err := sdkClient.Services.AppsV2.Update(someApp) // saves the new state to OneLogin
  ```
- **Destroy an app:**
  ```go
  err := sdkClient.Services.AppsV2.Destroy(*someApp.ID)
  ```

---

## Build & local development

> [!NOTE]
> Prerequisite: [Go](https://golang.org/doc/install). CI uses the default toolchain of `ubuntu-latest`; `go.mod` declares `go 1.12`.

```bash
go build ./...   # same as: make build
go test ./...    # unit tests, no network needed
```

| Make target | What it does |
|---|---|
| `make build` | `go build ./...` |
| `make test` | Installs `gopherbadger` v2.4.0, runs the tests and **rewrites the coverage badge in `readme.md`** |
| `make secure` | Downloads and runs `gosec` (excluding rule `G104`) |
| `make link` | Symlinks the repo into `$GOPATH/src/github.com/glovo/` |
| `make run` | Broken: runs `./cmd/main.go`, which does not exist |

> [!WARNING]
> `make test` edits `readme.md` in place: it replaces the `gopherbadger-tag-do-not-edit` badge. Keep that badge line intact, and don't commit the badge change by accident.

---

## Testing

```bash
go test ./...
```

- 24 test files, next to the code they cover, in every `pkg/` package and in [`internal/customerrors`](internal/customerrors).
- Services are tested against a mock repository ([`internal/test/test_repository.go`](internal/test/test_repository.go)), so no OneLogin account is needed.
- The coverage badge (95%) is generated by `make test`, not measured by CI.

---

## Code quality

| Tool | Role |
|---|---|
| `go fmt`, `go vet` | Run in CI ([`go.yml`](.github/workflows/go.yml)) |
| [`gosec`](Makefile) | Security scan, run on demand with `make secure` |
| [Go Report Card](https://goreportcard.com/report/github.com/glovo/onelogin-go-sdk) | Public grade, linked from the badge |

No git hooks are configured.

---

## CI/CD

| Workflow | Trigger | Effect |
|---|---|---|
| [`go.yml`](.github/workflows/go.yml) (Test) | Pull requests, and pushes to `main`, that touch `pkg/**`, `internal/**`, `go.mod` or `go.sum` | `go fmt`, `go build`, `go vet`, `go test` |
| [`release.yml`](.github/workflows/release.yml) (Release) | Push of a tag matching `v*` | Creates a published GitHub release (non-draft, non-pre-release) with the `antonyurchenko/git-release` action |

> [!NOTE]
> `go.yml` lists its own path as `.github/worflows/go.yml` (typo), so editing the workflow file does not trigger it. README-only changes do not run CI.

To release, push a tag such as `v1.1.24`. Consumers pick the version up with `go get`.

---

## Repository layout

```text
.
├── pkg/
│   ├── client/            # ⭐ NewClient, APIClientConfig, Services (entry point)
│   ├── oltypes/           # pointer helpers: String, Int32, Int64, Bool
│   ├── services/          # one package per OneLogin API (apps, users, roles, ...)
│   │   ├── olhttp/        # authenticated HTTP repository used by every service
│   │   └── interfaces.go  # Repository and HTTPClient interfaces
│   └── utils/             # shared helpers
├── internal/
│   ├── customerrors/      # error types
│   └── test/              # mock repository for tests
├── .github/workflows/     # go.yml (tests), release.yml (tag releases)
├── Makefile               # build, test, secure, link, run (broken)
├── CHANGELOG.md           # changelog (empty file)
├── config.yml             # contains only: version "0.2.1"
└── readme.md              # this file (badge at the top is auto-updated by make test)
```

---

## Contributing & conventions

- Each API lives in its own package under [`pkg/services`](pkg/services) with a `model.go`, a versioned service file (`v1.go` or `v2.go`) and a `*_test.go`. Follow that shape for a new API, then register it on `Services` in [`pkg/client/client.go`](pkg/client/client.go).
- Keep models' fields as pointers with `omitempty` JSON tags, as the existing models do.
- Run `go fmt`, `go vet` and `go test ./...` before opening a PR; CI runs the same.
- Tag releases as `v<major>.<minor>.<patch>`; every `v*` tag publishes a GitHub release.

---

## Ownership & support

- **Team:** no `CODEOWNERS` file and no catalog entry exist in this repo, so no owning team is recorded.
- **Upstream:** this is a fork of OneLogin's Go SDK (v1.1.20); the license is OneLogin's ([`LICENSE`](LICENSE)).
