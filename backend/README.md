# Questionnaire Back-end

Back-end application for an API first platform for creating and running questionnaires.

## Requirements

- Go 1.27+
- [Task](https://taskfile.dev) (optional, for the commands below)

## Getting Started

```bash
task build        # generate server stubs and build ./bin/server
./bin/server      # or: go run ./cmd/server
```

The server listens on `:8080`. Data is kept in memory and is lost on restart.


## Design Decisions

- Follow Domain Driven Principles
- OpenAPI first, server stubs are generated from `../api/openapi.yaml`
- Keep dependencies to a minimum
