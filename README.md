# n8n-launcher

CLI utility to securely manage [n8n task runners](https://docs.n8n.io/PENDING).

```sh
./n8n-launcher launch -type javascript
2024/11/15 13:53:33 Starting to execute `launch` command...
2024/11/15 13:53:33 Loaded config file loaded with a single runner config
2024/11/15 13:53:33 Changed into working directory: /Users/ivov/Development/n8n-launcher/bin
2024/11/15 13:53:33 Filtered environment variables
2024/11/15 13:53:33 Authenticated with n8n main instance
2024/11/15 13:53:33 Launching runner...
2024/11/15 13:53:33 Command: /usr/local/bin/node
2024/11/15 13:53:33 Args: [/Users/ivov/Development/n8n/packages/@n8n/task-runner/dist/start.js]
2024/11/15 13:53:33 Env vars: [LANG PATH TERM N8N_RUNNERS_N8N_URI N8N_RUNNERS_GRANT_TOKEN]
```

## Setup

### Install

- Install Node.js >=18.17 
- Install n8n >= `<PENDING_VERSION>`
- Download launcher binary from [releases page](https://github.com/n8n-io/task-runner-launcher/releases)

### Config

Create a config file for the launcher:

- at `config.json` in the dir containing the launcher binary, or
- at `/etc/n8n-task-runners.json` if `SECURE_MODE=true`.

Sample config file:

```json
{
  "task-runners": [
    {
      "runner-type": "javascript",
      "workdir": "/usr/local/bin",
      "command": "/usr/local/bin/node",
      "args": [
        "/usr/local/lib/node_modules/n8n/node_modules/@n8n/task-runner/dist/start.js"
      ],
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

Task runner config fields:

- `runner-type`: Type of task runner, currently only `javascript` supported
- `workdir`: Path to directory containing the task runner binary
- `command`: Command to execute to start task runner
- `args`: Args for command to execute, currently path to task runner entrypoint
- `allowed-env`: Env vars allowed to be passed to the task runner

### Auth

Generate a secret auth token (e.g. random string) for the launcher to authenticate with the n8n main instance. During the `launch` command, the launcher will exchange this auth token for a grant token from the n8n instance, and pass the grant token to the runner.

## Usage

Once setup is complete, start the launcher:

```sh
export N8N_RUNNERS_AUTH_TOKEN=...
export N8N_RUNNERS_N8N_URI=... 
./n8n-launcher launch -type javascript
```

Or in secure mode:

```sh
export SECURE_MODE=true
export N8N_RUNNERS_AUTH_TOKEN=...
export N8N_RUNNERS_N8N_URI=... 
./n8n-launcher launch -type javascript
```

## Development

1. Install Go >=1.23

2. Clone repo and create config file:

```sh
git clone https://github.com/n8n-io/PENDING-NAME
cd PENDING_NAME
touch config.json && echo '<json-config-content>' > config.json 
```

3. Start n8n:

```sh
export N8N_RUNNERS_ENABLED=true
export N8N_RUNNERS_MODE=external 
export N8N_RUNNERS_LAUNCHER_PATH=...
export N8N_RUNNERS_AUTH_TOKEN=...
pnpm start
```

4. Make changes to launcher.

5. Build and run launcher:

```sh
go build -o bin cmd/launcher/main.go

export N8N_RUNNERS_N8N_URI=...
export N8N_RUNNERS_AUTH_TOKEN=...
./bin/main launch -type javascript
```
