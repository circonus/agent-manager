# Example Agent

Requires:

* [go](https://go.dev/dl/)
* [goreleaser](https://goreleaser.com/install/)
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

To build the example: `goreleaser build --clean --snapshot`

```shell
dist/ea_darwin_amd64_v1/sbin/example-agentd
{"level":"info","name":"example-agent","version":"dev","time":1682611949,"message":"starting"}
{"level":"info","example_arg":"default","time":1682611949,"message":"example argument"}

dist/ea_darwin_amd64_v1/sbin/example-agentd --example-arg="example cli"
{"level":"info","name":"example-agent","version":"dev","time":1682612993,"message":"starting"}
{"level":"info","example_arg":"example cli","time":1682612993,"message":"example argument"}

dist/ea_darwin_amd64_v1/sbin/example-agentd --config=etc/example-agent.yaml
{"level":"info","name":"example-agent","version":"dev","time":1682611978,"message":"starting"}
{"level":"info","example_arg":"example cfg yaml","time":1682611978,"message":"example argument"}

EA_EXAMPLE_ARG="example env" dist/ea_darwin_amd64_v1/sbin/example-agentd
{"level":"info","name":"example-agent","version":"dev","time":1682611996,"message":"starting"}
{"level":"info","example_arg":"example env","time":1682611996,"message":"example argument"}
```

## To use as a template

1. Click [Use this template], select _Create a new repository_ from the drop-down list
2. Clone the new repository
    1. Rename the `example-agent` subdirectory under `cmd/` to refelect the agent being built
    1. Edit `go.mod` to update the `module` path to refelct the new repo
    1. Edit all files which refrence imports for this repo to reflect the correct module path for the new repo
    1. Add additional package in `internal/` to flesh out the agent
    1. Update `internal/agent/agent.go` to use the new packages and go into a "running" state
    1. Update `goreleaser.yml` to reflect the binary being built and the new repository
