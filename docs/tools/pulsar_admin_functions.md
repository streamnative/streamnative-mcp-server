#### pulsar_admin_functions

Manage Apache Pulsar Functions for stream processing. Pulsar Functions are lightweight compute processes that can consume messages from one or more Pulsar topics, apply user-defined processing logic, and produce results to another topic. Functions support Java, Python, and Go runtimes, enabling complex event processing, data transformations, filtering, and integration with external systems.

This tool provides a comprehensive set of operations to manage the entire function lifecycle:

- **list**: List all functions in a namespace
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name

- **get**: Get function configuration
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **status**: Get runtime status of a function (instances, metrics)
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **stats**: Get detailed statistics of a function (throughput, processing latency)
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **create**: Deploy a new function with specified parameters
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name
  - `classname` (string, required): The fully qualified class name implementing the function
  - `inputs` (array, required): The input topics for the function
  - `output` (string, optional): The output topic for function results
  - `jar` (string, optional): Path to the JAR file for Java functions
  - `py` (string, optional): Path to the Python file for Python functions
  - `go` (string, optional): Path to the Go binary for Go functions
  - `parallelism` (number, optional): The parallelism factor of the function (default: 1)
  - `userConfig` (object, optional): User-defined config key/values

- **update**: Update an existing function
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name
  - Parameters similar to `create` operation

- **delete**: Delete a function
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **start**: Start a stopped function
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **stop**: Stop a running function
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **restart**: Restart a function
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name

- **querystate**: Query state stored by a stateful function for a specific key
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name
  - `key` (string, required): The state key to query

- **putstate**: Store state in a function's state store
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name
  - `key` (string, required): The state key
  - `value` (string, required): The state value

- **trigger**: Manually trigger a function with a specific value
  - `tenant` (string, required): The tenant name
  - `namespace` (string, required): The namespace name
  - `name` (string, required): The function name
  - `topic` (string, optional): The specific topic to trigger on
  - `triggerValue` (string, optional): The value to trigger the function with 