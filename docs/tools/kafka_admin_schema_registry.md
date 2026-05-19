#### kafka-admin-schema-registry

**Claude connector safety:** Actual MCP tools are split into `kafka_admin_sr_read` and `kafka_admin_sr_write`. The read tool is read-only and only exposes read operations/parameters. The write tool is destructive and is not registered in read-only mode.

### `kafka_admin_sr_read`

Read Schema Registry subjects, versions, schemas, compatibility levels, and supported schema types.

- **subjects**
  - **list**: List all schema subjects

- **subject**
  - **get**: Get the latest schema for a subject
    - `subject` (string, required): Subject name

- **versions**
  - **list**: List versions for a subject
    - `subject` (string, required): Subject name

- **version**
  - **get**: Get one version of a subject schema
    - `subject` (string, required): Subject name
    - `version` (string, required): Version number or `latest`

- **compatibility**
  - **get**: Get compatibility level globally or for a subject
    - `subject` (string, optional): Subject name for subject-specific compatibility

- **types**
  - **list**: List supported schema types

### `kafka_admin_sr_write`

Register/delete schemas and change compatibility levels.

- **subject**
  - **create**: Register a new schema for a subject
    - `subject` (string, required): Subject name
    - `schemaType` (string, required): Schema type, such as `AVRO`, `JSON`, or `PROTOBUF`
    - `schema` (object, required): Schema definition
  - **delete**: Delete a schema subject
    - `subject` (string, required): Subject name

- **version**
  - **delete**: Delete one version of a subject schema
    - `subject` (string, required): Subject name
    - `version` (string, required): Version number

- **compatibility**
  - **set**: Set compatibility level globally or for a subject
    - `compatibility` (string, required): Compatibility level, such as `BACKWARD`, `FORWARD`, `FULL`, or `NONE`
    - `subject` (string, optional): Subject name for subject-specific compatibility
