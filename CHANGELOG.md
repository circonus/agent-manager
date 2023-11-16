# unreleased

## v0.2.16

* feat: always update agent manager version on start-up
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.24.0 to 1.25.0
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.23.0 to 1.24.0

## v0.2.15

* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.22.3 to 1.23.0
* chore(deps): bump github.com/golang-jwt/jwt/v5 from 5.0.0 to 5.1.0

## v0.2.14

* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.21.0 to 1.22.3
* chore(deps): bump golang.org/x/sync from 0.4.0 to 0.5.0
* chore(deps): bump golang.org/x/sys from 0.13.0 to 0.14.0
* chore(deps): bump github.com/spf13/cobra from 1.7.0 to 1.8.0

## v0.2.13

* chore(deps): bump github.com/shirou/gopsutil/v3 from 3.23.9 to 3.23.10
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.19.1 to 1.21.0
* chore(deps): bump github.com/aws/aws-sdk-go-v2/feature/ec2/imds from 1.13.13 to 1.14.0

## v0.2.12

* fix: ensure all config items are in config struct
* fix: add defaults for all config items to example config file

## v0.2.11

* chore(deps): bump github.com/google/uuid from 1.3.1 to 1.4.0
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.19.0 to 1.19.1
* feat: add health endpoint and config trigger endpoint
* chore: refactor environment code (platform, is docker, etc.)
* fix: return exit code and stderr when command fails

## v0.2.10

* chore(deps): bump github.com/aws/aws-sdk-go-v2/feature/ec2/imds from 1.13.11 to 1.13.13
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.18.43 to 1.19.0
* feat: add agent status [C3-1710]
* feat: add agent status poller interval
* feat: add config tracker poller
* feat: add agent status poller (for non-docker)
* chore: performance improvement
* chore: reformat error message
* feat(tracker): send status to API
* feat(tracker): send modified status once, not on every check
* feat: config tracking for local changes [C3-1697]
* chore: update agents and inventory testdata json files
* chore: keep tracker testdata dir
* chore: ignore tracker dynamic testdata
* chore: refactor to eliminate circular dependencies
* chore(deps): bump golang.org/x/sync from 0.3.0 to 0.4.0
* chore(deps): bump golang.org/x/sys from 0.12.0 to 0.13.0
* chore(deps): bump github.com/spf13/viper from 1.16.0 to 1.17.0

## v0.2.9

* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.18.39 to 1.18.43
* chore(deps): bump github.com/shirou/gopsutil/v3 from 3.23.8 to 3.23.9
* chore(deps): bump github.com/rs/zerolog from 1.30.0 to 1.31.0
* fix: is registered check to catch for inventory as well
* feat: add already registered message
* feat: add --force-register option (force a re-registration for a manager)
* feat: add IsRegistered to prevent manager attempting to re-register itself every time a container restarts
* feat: ensure an existing credential file is >0 bytes if it exists
* feat(docker): add all assets to /cam to ensure correct etc/ usage
* doc: update flags documentation
* fix: typos
* feat(rpm): sign rpm files
* chore: print env var for register
* chore(cfgbak): add mode to err msg when src not regular file
* chore(cfgbak): refactor order in cfg loop
* chore(cfgbak): update text of completion message
* chore: breakout config backup
* feat: backup configs on registration and inventory
* feat: register and run when in container
* feat: add --instance-id and --agents args for containers
* feat: add --instance-id and use ss override for hostname
* feat: add container detector

## v0.2.8

* feat: add 30s timeout to commands

## v0.2.7

* feat: break out darwin into arm64 and amd64
* feat: add check for registration and emit a more clear error message if not registered
* chore: ignore test files/dirs in etc
* chore: add default api url to example conf
* feat: update etc path if specific config supplied in non-default location
* feat(inventory): add cmd output on error getting version
* chore: remove old code
* chore(deps): update gopsutil to v3
* chore: ignore inventory.yaml in testdata
* fix: keys for tests

## v0.2.6

* feat: switch default to prod api fqdn agents-api.circonus.app
* chore: add -trimpath to flags
* build: Package trigger automation (#39)
* fix(nfpms): change license from MIT to BSD-3-Clause

## v0.2.5

* feat: promote custom reg tags to top-level
* doc: add note on tag format for env and cli

## v0.2.4

* chore: ensure manager id is loaded for refresh
* feat: add some progress debug messages to decommission
* fix: incorrect flag type for setting hidden
* fix: error message on token refresh
* feat: use machine id option
* feat: decommission the agent manager
* doc: change start to restart after registration (linux)
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.18.32 to 1.18.39
* chore(deps): bump golangci/golangci-lint-action from 3.6.0 to 3.7.0
* chore(deps): bump golang.org/x/sys from 0.11.0 to 0.12.0
* chore(deps): bump actions/checkout from 3 to 4

## v0.2.3

* feat(build): add docker images

## v0.2.2

* fix: load jwt after refresh

## v0.2.1

* feat(inventory): just return whatever the version command outputs do not verify semver

## v0.2.0

* feat: add examples for aws_ec2_tags and tags config options
* feat: tags config option
* feat: add tags cli option
* feat: single use registration token
* feat: registration token to refresh token
* feat: custom tags support
* feat: use Authorization header
* feat: access and refresh token support
* feat: refresh token on 401 from api
* feat: Authorization header
* chore: refactor reload command handling
* chore: refactor command running to simplify
* feat: add http endpoint handling option for reload command
* feat: only check for binary during agent inventory
* chore: update instructions for start/restart in brew formula
* fix: brew creating subdir in etc
* fix: brew instructions for start/restart

## v0.1.4

* fix: remove term 'collectors' from inventory option

## v0.1.3

* fix: do not run as cua with systemd

## v0.1.2

* fix: dangling cma names in service files

## v0.1.1

* fix: update default api url

## v0.1.0

* doc: add some linux specific install instructions
* feat: add darwin support w/brew
* chore: ignore vagrant stuff for testing
* feat: add mac signing script
* fix: remove unused token attribute
* feat: exit after --register or --inventory actions
* fix: update to new structure layout for agent_type endpoint
* chore: ignore .envrc
* feat: name collectors -> agents
* feat: rename agent id to manager id
* feat: rename agent to manager
* chore: update to go1.20.6 (GO-2023-1878)
* feat: rename api endpoints/json manager registration
* feat: chown configs to original gid/uid
* chore(deps): bump golang.org/x/sys from 0.10.0 to 0.11.0
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.18.31 to 1.18.32

## v0.0.4

* feat: add default api url

## v0.0.3

* fix: rename env prefix CAM

## v0.0.2

* feat: name change cma -> am (agent manager)

## v0.0.1

* initial
