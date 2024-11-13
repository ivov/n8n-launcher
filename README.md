# n8n-launcher

CLI utility to securely launch and kill [n8n task runners](https://docs.n8n.io/pending).

Task runners are nodes that execute tasks in n8n workflows, specifically Code node execution tasks.

n8n may use a runner:

- in `internal` mode, where n8n spawns the task runner as a child process (npm users), or
- in `external` mode, where an orchestrator starts the runner via a secure launcher (Docker).

## Usage

```sh
# launch a task runner
./task-runner-launcher launch -type javascript

# kill a task runner
./task-runner-launcher kill -type javascript -pid <process_id>
```

## Setup

- Download the launcher binary for your target platform from the [releases page](https://github.com/n8n-io/task-runner-launcher/releases)
- Create a [config file](#config-file) and make it accessible to the launcher binary
- Install n8n@<pending_version>, which contains the `@n8n/task-runner` package
- Install Node.js >=20.15

### Config file

The launcher expects a JSON config file

- at `config.json` in the directory containing the binary, or
- at `/etc/n8n-task-runners.json` if `SECURE_MODE=true`.

Sample config file:

```jsonc
{
  "task-runners": [
    {
      // type of task runner, only "javascript" currently supported
      "runner-type": "javascript",

      // path to directory containing launcher (Go binary)
      "workdir": "/usr/local/bin",

      // command to start runner
      "command": "/usr/local/bin/node",

      // arguments containing path to task runner entrypoint
      "args": [
        "/usr/local/lib/node_modules/n8n/node_modules/@n8n/task-runner/dist/start.js"
      ],

      // env vars allowed to be passed to task runner
      // allowed by default: LANG, PATH, TZ, TERM
      "allowed-env": [
        "PATH",
        "N8N_RUNNERS_GRANT_TOKEN",
        "N8N_RUNNERS_N8N_URI",
        "N8N_RUNNERS_MAX_PAYLOAD",
        "N8N_RUNNERS_MAX_CONCURRENCY",
        "NODE_FUNCTION_ALLOW_BUILTIN",
        "NODE_FUNCTION_ALLOW_EXTERNAL",
        "NODE_OPTIONS"
      ]
    }
  ]
}
```

### Auth

The launcher optionally supports authenticating with n8n through a grant token.

To enable auth, set `N8N_RUNNERS_AUTH_TOKEN` to the authentication token for n8n, and set `N8N_RUNNERS_N8N_URI` to the URI of the n8n instance. When these set, the launcher will fetch a grant token from the n8n instance, and pass the grant token to the runner through the `N8N_RUNNERS_GRANT_TOKEN` env vars.

## Development

### Setup

1. Install Go >=1.23

2. Clone repository and create a config file:

```sh
git clone https://github.com/n8n-io/task-runner-launcher
cd task-runner-launcher

touch config.json && echo '<your-config>' > config.json 
mv config.json /etc/n8n-task-runners.json # if secure
```

3. Install n8n and start it with `N8N_RUNNERS_ENABLED=true` and `N8N_RUNNERS_MODE=external`

4. Build and run launcher:

```sh
go build -o bin cmd/task-runner-launcher/main.go

N8N_RUNNERS_N8N_URI=127.0.0.1:5679 ./bin/main launch -type javascript # or
N8N_RUNNERS_N8N_URI=127.0.0.1:5679 SECURE_MODE=true ./bin/main launch -type javascript
```
