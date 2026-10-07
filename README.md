# rikami-cli

A small Go CLI (`rika`) that uploads an application's `.env` files to AWS SSM Parameter Store, so my Kubernetes clusters can pull them in through External Secrets.

While developing an app I keep its settings in local files like `.env.staging` and `.env.prod`. One command turns such a file into a single encrypted SSM parameter in the path my clusters already read from. A [rikami-operator](https://github.com/b-zago/rikami-operator) Vessel then loads it into the app as environment variables, with no secrets in Git and no manual copying.

## How it works

```
.env.staging  →  rika params put .env.staging api-keys  →  /clusters/staging/<app>/api-keys (SecureString, JSON)
```

The file's suffix decides the environment, and `rikami.yml` provides the rest of the path. The `KEY=value` lines are stored as one JSON object, so External Secrets can turn every key back into an environment variable.

## Install

Download a binary for Linux (amd64, arm64) or Windows (amd64) from [Releases](https://github.com/b-zago/rikami-cli/releases).

Or build it with Go:

```sh
go install github.com/b-zago/rikami-cli@latest   # installs the binary as rikami-cli
```

AWS credentials come from the standard AWS configuration (environment variables, `~/.aws`, SSO, and so on).

## Configuration

Create `rikami.yml` in the directory you run the CLI from, usually the app's repository:

```yaml
app: simple-api # app name, part of the SSM path
ssmPrefix: /clusters # root of the SSM path
envFilesPrefix: .env # stripped from file names to get the environment
```

## Usage

**Upload a `.env` file:**

```sh
rika params put .env.staging api-keys
# Uploaded to SSM path /clusters/staging/simple-api/api-keys
```

Existing parameters are overwritten, so re-running it updates the values.

**Read a parameter back:**

```sh
rika params get /clusters/staging/simple-api/api-keys
# {"API_KEY":"...","DB_URL":"..."}
```

**Use it in the cluster** with a rikami Vessel:

```yaml
externalSecrets:
  - name: api-keys
    extract: /clusters/staging/simple-api/api-keys
```

## Releases

Every push to `main` creates a new version with my [semver-bump action](https://github.com/b-zago/actions) and publishes binaries with GoReleaser.
