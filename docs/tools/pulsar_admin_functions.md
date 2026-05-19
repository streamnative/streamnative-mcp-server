#### pulsar_admin_functions

**Claude connector safety:** Actual MCP tools are split into `pulsar_admin_functions_read` and `pulsar_admin_functions_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

Pulsar Functions are lightweight compute processes that consume messages from Pulsar topics, apply user-defined logic, and optionally produce results to another topic.

### `pulsar_admin_functions_read`

Read function lists, configuration, runtime status, statistics, state, and package data.

Common identity parameters:

- `fqfn` (string, optional): Fully qualified function name in `tenant/namespace/name` form
- `tenant` (string, optional): Tenant name; default `public`
- `namespace` (string, optional): Namespace name; default `default`
- `name` (string, required for operations targeting one function unless `fqfn` is provided): Function name

Operations:

- **list**: List functions in a namespace
  - `tenant` (string, optional)
  - `namespace` (string, optional)
- **get**: Get function configuration
  - Common identity parameters
- **status**: Get runtime status
  - Common identity parameters
  - `instanceId` (number, optional): Function instance ID for per-instance status
- **stats**: Get runtime statistics
  - Common identity parameters
  - `instanceId` (number, optional): Function instance ID for per-instance stats
- **querystate**: Query state for a key
  - Common identity parameters
  - `key` (string, required): State key
- **download**: Download function package data
  - `destinationFile` (string, required): Local destination path
  - `path` (string, optional): Direct Pulsar package storage path
  - Common identity parameters when downloading by function identity

### `pulsar_admin_functions_write`

Manage function lifecycle, runtime state, manual triggers, and package uploads.

Common identity parameters:

- `fqfn` (string, optional): Fully qualified function name in `tenant/namespace/name` form
- `tenant` (string, optional): Tenant name; default `public`
- `namespace` (string, optional): Namespace name; default `default`
- `name` (string, required for operations targeting one function unless `fqfn` is provided): Function name

Operations:

- **create**: Deploy a function
  - Common identity parameters
  - Deployment parameters include `classname`, `functionType`, `inputs`, `topicsPattern`, `inputSpecs`, `output`, `jar`, `py`, `go`, `parallelism`, resource settings (`cpu`, `ram`, `disk`), schemas/serde settings, processing guarantees, subscription settings, retry/dead-letter settings, secrets, window settings, and `functionConfigFile`
- **update**: Update function configuration
  - Common identity parameters
  - Same deployment parameters as `create`, plus `updateAuthData`
- **delete**: Delete a function
  - Common identity parameters
- **start**: Start a stopped function
  - Common identity parameters
- **stop**: Stop a running function
  - Common identity parameters
- **restart**: Restart a function
  - Common identity parameters
- **putstate**: Store a value in a function state store
  - Common identity parameters
  - `key` (string, required): State key
  - `value` (string, required): State value
- **trigger**: Trigger a function manually
  - Common identity parameters
  - `topic` (string, optional): Input topic to trigger on
  - exactly one of `triggerValue` or `triggerFile`
- **upload**: Upload a local file into Pulsar function package storage
  - `sourceFile` (string, required): Local source path
  - `path` (string, required): Pulsar package storage destination path
