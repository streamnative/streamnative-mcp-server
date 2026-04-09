#### pulsar_admin_functions

Manage Apache Pulsar Functions for stream processing. Pulsar Functions are lightweight compute processes that can consume messages from one or more Pulsar topics, apply user-defined processing logic, and produce results to another topic. Functions support Java, Python, and Go runtimes, enabling complex event processing, data transformations, filtering, and integration with external systems.

This tool provides a comprehensive set of operations to manage the entire function lifecycle:

- **list**: List all functions in a namespace
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)

- **get**: Get function configuration
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)

- **status**: Get runtime status of a function (instances, metrics)
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - `instanceId` (string, optional): Instance ID for per-instance status

- **stats**: Get detailed statistics of a function (throughput, processing latency)
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - `instanceId` (string, optional): Instance ID for per-instance stats

- **create**: Deploy a new function with specified parameters
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, optional): The function name (can be inferred from `classname`)
  - `classname` (string, optional): The fully qualified class name implementing the function
  - `functionType` (string, optional): Built-in function type (translated to `builtin://`)
  - `inputs` (array, optional): The input topics for the function
  - `topicsPattern` (string, optional): Topic pattern to consume from
  - `inputSpecs` (object, optional): Map of input topics to consumer config
  - `output` (string, optional): The output topic for function results
  - `jar` (string, optional): Path to the JAR file for Java functions
  - `py` (string, optional): Path to the Python file for Python functions
  - `go` (string, optional): Path to the Go binary for Go functions
  - `parallelism` (number, optional): The parallelism factor of the function (default: 1)
  - `cpu` (number, optional): CPU cores per instance
  - `ram` (number, optional): RAM bytes per instance
  - `disk` (number, optional): Disk bytes per instance
  - `userConfig` (object, optional): User-defined config key/values
  - `producerConfig` (object, optional): Custom producer configuration
  - `logTopic` (string, optional): Topic for function logs
  - `schemaType` (string, optional): Output schema type or class
  - `outputSerdeClassName` (string, optional): Output SerDe class
  - `customSerdeInputs` (object, optional): Map of input topics to SerDe class
  - `customSchemaInputs` (object, optional): Map of input topics to Schema class
  - `customSchemaOutputs` (object, optional): Map of output topics to schema properties
  - `inputTypeClassName` (string, optional): Input type class name
  - `outputTypeClassName` (string, optional): Output type class name
  - `processingGuarantees` (string, optional): Delivery semantics
  - `retainOrdering` (boolean, optional): Process messages in order
  - `retainKeyOrdering` (boolean, optional): Process messages in key order
  - `batchBuilder` (string, optional): Batch builder type
  - `forwardSourceMessageProperty` (boolean, optional): Forward properties to output
  - `autoAck` (boolean, optional): Automatically acknowledge messages
  - `subsName` (string, optional): Subscription name for inputs
  - `subsPosition` (string, optional): Subscription position
  - `skipToLatest` (boolean, optional): Skip to latest on restart
  - `timeoutMs` (number, optional): Message timeout in ms
  - `maxMessageRetries` (number, optional): Max retries
  - `deadLetterTopic` (string, optional): Dead letter topic
  - `customRuntimeOptions` (string, optional): Custom runtime options
  - `secrets` (object, optional): Secrets map
  - `cleanupSubscription` (boolean, optional): Clean up subscription on delete
  - `windowLengthCount` (number, optional): Window length count
  - `windowLengthDurationMs` (number, optional): Window length duration in ms
  - `slidingIntervalCount` (number, optional): Sliding interval count
  - `slidingIntervalDurationMs` (number, optional): Sliding interval duration in ms
  - `functionConfigFile` (string, optional): YAML config file path

- **update**: Update an existing function
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - Parameters similar to `create` operation
  - `updateAuthData` (boolean, optional): Whether to update auth data

- **delete**: Delete a function
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)

- **download**: Download function package data from Pulsar to local storage
  - `destinationFile` (string, required): Local file path where the downloaded content should be written
  - `path` (string, optional): Direct Pulsar package storage path to download from
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name` when downloading by function identity
  - `tenant` (string, optional): The tenant name (default: `public`) when downloading by function identity
  - `namespace` (string, optional): The namespace name (default: `default`) when downloading by function identity
  - `name` (string, required unless `path` or `fqfn` is provided): The function name when downloading by function identity

- **start**: Start a stopped function
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)

- **stop**: Stop a running function
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)

- **restart**: Restart a function
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)

- **querystate**: Query state stored by a stateful function for a specific key
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - `key` (string, required): The state key to query

- **putstate**: Store state in a function's state store
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - `key` (string, required): The state key
  - `value` (string, required): The state value

- **trigger**: Manually trigger a function with a specific value
  - `fqfn` (string, optional): Fully qualified function name `tenant/namespace/name`
  - `tenant` (string, optional): The tenant name (default: `public`)
  - `namespace` (string, optional): The namespace name (default: `default`)
  - `name` (string, required): The function name (unless `fqfn` is provided)
  - `topic` (string, optional): The specific topic to trigger on
  - `triggerValue` (string, optional): The value to trigger the function with
  - `triggerFile` (string, optional): File path containing the trigger value

- **upload**: Upload a local file into Pulsar function package storage
  - `sourceFile` (string, required): Local file path whose content should be uploaded
  - `path` (string, required): Pulsar package storage path where the file should be stored
