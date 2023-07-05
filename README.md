# Circonus Agent Manager (circonus-am)

Requires:

* [go](https://go.dev/dl/)
* [goreleaser](https://goreleaser.com/install/)
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

1. Build (for testing): `goreleaser build --clean --snapshot`
1. Release:
   1. Ensure all commits/PRs are merged
   1. Ensure repo is up-to-date
   1. Tag with semver from CHANGELOG.md
   1. Run `goreleaser release --clean`

Configuration:

```text
Manage locally installed agents (metrics, logs, etc.)

Usage:
  circonus-cma [flags]

Flags:
      --apiurl string              [ENV: CAM_API_URL] Circonus API URL (default "https://something.circonus.com")
      --aws-ec2-tags stringArray   [ENV: CAM_AWS_EC2_TAGS] AWS EC2 tags for registration meta data
  -c, --config string              config file (default: /opt/circonus/cma/etc/circonus-cma.yaml|.json|.toml)
  -d, --debug                      [ENV: CAM_DEBUG] Enable debug messages
  -h, --help                       help for circonus-cma
      --inventory                  [ENV: CAM_INVENTORY] Inventory installed collectors
      --log-level string           [ENV: CAM_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty                 Output formatted/colored log lines [ignored on windows]
      --poll-interval string       [ENV: CAM_POLL_INTERVAL] Polling interval for actions (default "60s")
      --register string            [ENV: CAM_REGISTER] Registration token
  -V, --version                    Show version and exit
  ```
