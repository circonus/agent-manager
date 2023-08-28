# unreleased

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
* feat: cutom tags support
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

* doc: add some linux spcefic install instructions
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
