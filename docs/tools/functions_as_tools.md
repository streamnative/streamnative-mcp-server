# Functions as Tools

The "Functions as Tools" feature allows the StreamNative MCP Server to dynamically discover Apache Pulsar Functions deployed in your cluster and expose them as invokable MCP tools for AI agents. This significantly enhances the capabilities of AI agents by allowing them to interact with custom business logic encapsulated in Pulsar Functions without manual tool registration for each function.

## How it Works

### 1. Function Discovery
The MCP Server automatically discovers Pulsar Functions available in the connected Pulsar cluster. It periodically polls for functions and identifies those suitable for exposure as tools.

By default, if no custom name is provided (see Customizing Tool Properties), the MCP tool name might be derived from the Function's Fully Qualified Name (FQN), such as `pulsar_function_$tenant_$namespace_$name`.

### 2. Schema Conversion
For each discovered function, the MCP Server attempts to extract its input and output schema definitions. Pulsar Functions can be defined with various schema types for their inputs and outputs (e.g., primitive types, AVRO, JSON).

The server then converts these native Pulsar schemas into a format compatible with MCP tools. This allows the AI agent to understand the expected input parameters and the structure of the output.

Supported Pulsar schema types for automatic conversion include:
*   Primitive types (String, Boolean, Numbers like INT8, INT16, INT32, INT64, FLOAT, DOUBLE)
*   AVRO
*   JSON

If a function uses an unsupported schema type for its input or output, or if schemas are not clearly defined, it might not be exposed as an MCP tool.

## Enabling the Feature
To enable this functionality, you need to specific the default `--pulsar-instance` and `--pulsar-cluster`, and include `functions-as-tools` in the `--features` flag when starting the StreamNative MCP Server.

Example:
```bash
snmcp stdio --organization my-org --key-file /path/to/key-file.json --features pulsar-admin,pulsar-client,functions-as-tools --pulsar-instance instance --pulsar-cluster cluster
```
If `functions-as-tools` is part of a broader feature set like `all` and `streamnative-cloud`, enabling `all` or `streamnative-cloud` would also activate this feature.

## Customizing Tool Properties
You can customize how your Pulsar Functions appear as MCP tools (their name and description) by providing specific runtime options when deploying or updating your functions. This is done using the `--custom-runtime-options` flag with `pulsar-admin functions create` or `pulsar-admin functions update`.

The MCP Server looks for the following environment variables within the custom runtime options:
*   `MCP_TOOL_NAME`: Specifies the desired name for the MCP tool.
*   `MCP_TOOL_DESCRIPTION`: Provides a description for the MCP tool, which helps the AI agent understand its purpose.

**Format for `--custom-runtime-options`**:
The options should be a JSON string where you define an `env` map containing `MCP_TOOL_NAME` and `MCP_TOOL_DESCRIPTION`.

**Example**:
When deploying a Pulsar Function, you can set these properties as follows:
```bash
pulsar-admin functions create \
  --tenant public \
  --namespace default \
  --name my-custom-logic-function \
  --inputs "persistent://public/default/input-topic" \
  --output "persistent://public/default/output-topic" \
  --py my_function.py \
  --classname my_function.MyFunction \
  --custom-runtime-options \
  '''
  {
    "env": {
      "MCP_TOOL_NAME": "CustomObjectFunction",
      "MCP_TOOL_DESCRIPTION": "Takes an input number and returns the value incremented by 100."
    }
  }
  '''
```
In this example:
- The MCP tool derived from `my-custom-logic-function` will be named `CustomObjectFunction`.
- Its description will be "Takes an input number and returns the value incremented by 100."

If these custom options are not provided, the MCP tool name might default to a derivative of the function's FQN, and the description might be generic and cannot help AI Agent to understand the purpose of the MCP tool.

## Considerations and Limitations

*   **Schema Definition**: For reliable schema conversion, ensure your Pulsar Functions have clearly defined input and output schemas using Pulsar's schema registry capabilities. Functions with ambiguous or `BYTES` schemas might not be converted effectively or might default to generic byte array inputs/outputs.
*   **Function State**: This feature primarily focuses on the stateless request/response invocation pattern of functions.
*   **Discovery Latency**: There might be a slight delay between deploying/updating a function and it appearing as an MCP tool, due to the server's polling interval for function discovery.
*   **Error Handling**: The MCP Server will attempt to relay errors from function executions, but the specifics might vary.
*   **Security**: Ensure that only intended functions are exposed by managing permissions within your Pulsar cluster. The MCP Server will operate with the permissions of its Pulsar client.
