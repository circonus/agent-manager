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
Manager for local agents (metrics, logs, etc.)

Usage:
  circonus-am [flags]

Flags:
      --agents strings         [ENV: CAM_AGENTS] List of agents (Docker specific)
      --apiurl string          [ENV: CAM_API_URL] Circonus API URL (default "https://agents-api.circonus.app/configurations/v1")
      --aws-ec2-tags strings   [ENV: CAM_AWS_EC2_TAGS] AWS EC2 tags for registration meta data
  -c, --config string          config file (default: /Users/mgm/src/circonus/agent-manager/dist/am-macos_amd64_darwin_amd64_v1/etc/circonus-am.yaml|.json|.toml)
  -d, --debug                  [ENV: CAM_DEBUG] Enable debug messages
      --decommission           Decommission agent manager and exit
  -h, --help                   help for circonus-am
      --instance-id string     [ENV: CAM_INSTANCE_ID] Instance ID (Docker specific)
      --inventory              Inventory installed agents and exit
      --log-level string       [ENV: CAM_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty             Output formatted/colored log lines [ignored on windows]
      --poll-interval string   [ENV: CAM_POLL_INTERVAL] Polling interval for actions (default "60s")
      --register string        [ENV: CAM_REGISTER] Registration token -- register agent manager, inventory installed agents and exit
      --tags strings           [ENV: CAM_TAGS] Custom key:value tags for registration meta data
  -V, --version                Show version and exit
  ```

## Note on tags

* Environment variable format with a space separated list, e.g. `CAM_TAGS="foo:bar baz:qux"`
* CLI option format with a comma separated list, e.g. `--tags="foo:bar,baz:qux"`

## Linux installation

1. Download appropriate package from releases page
1. Install (use `sudo` to install, like all packages)
1. Run `sudo /opt/circonus/am/sbin/circonus-am --register=<registration_token>`
1. If registration successful, restart the agent manager `sudo systemctl restart circonus-am`
1. Check status of agent manager `sudo systemctl status circonus-am`
1. If an additional agent is installed AFTER the agent manager has registered
   1. Stop agent manager `sudo systemctl stop circonus-am`
   1. Run `sudo /opt/circonus/am/sbin/circonus-am --inventory`
   1. Start agent manager `sudo systemctl start circonus-am`

## Unprivileged

1. Create dedicated user and group (e.g. `cam`)
1. Change ownership of `/opt/circonus/am` and all files to the user and group (ensure user can also access `/opt` and `/opt/circonus` e.g. 755 or owned by group)
1. Edit the systemd unit to add User and Group stanzas (`sysctl edit circonus-am`)
1. Restart circonus-am if it is currently running
1. For all agent config files that the manager will manage:
   1. change the group ownership to the group created above
   1. change permissions to allow write by group
1. Add sudo configs for the commands of each installed agent
1. Change agent definitions to include `sudo` for each of the commands used for managing the agent

## Decommission (linux)

1. `sudo systemctl stop circonus-am`
1. `sudo /opt/circonus/am/sbin/circonus-am --decommission`
1. Remove package using tool used to install e.g. apt/dpkg or yum/rpm
