#### pulsar_admin_schema

**Claude connector safety:** Actual MCP tools: `pulsar_admin_schema_read` (`get`) and `pulsar_admin_schema_write` (`upload`, `delete`).


Manage Apache Pulsar schemas for topics.

- **schema**
  - **get**: Get the schema for a topic
    - `topic` (string, required): The fully qualified topic name
    - `version` (number, optional): Schema version number
  - **upload**: Upload a new schema for a topic
    - `topic` (string, required): The fully qualified topic name
    - `filename` (string, required): Path to the schema definition file
  - **delete**: Delete the schema for a topic
    - `topic` (string, required): The fully qualified topic name
