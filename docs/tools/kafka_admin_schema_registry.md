#### kafka-admin-schema-registry

This tool provides access to Kafka Schema Registry operations, including managing subjects, versions, and compatibility settings.

- **subjects**
  - **list**: List all schema subjects in the Schema Registry
    - _Parameters_: None

- **subject**
  - **get**: Get the latest schema for a subject
    - `subject` (string, required): The subject name
  - **create**: Register a new schema for a subject
    - `subject` (string, required): The subject name
    - `schema` (string, required): The schema definition (in AVRO/JSON/PROTOBUF, etc.)
    - `type` (string, optional): The schema type (e.g. AVRO, JSON, PROTOBUF)
  - **delete**: Delete a schema subject
    - `subject` (string, required): The subject name

- **versions**
  - **list**: List all versions for a specific subject
    - `subject` (string, required): The subject name
  - **get**: Get a specific version of a subject's schema
    - `subject` (string, required): The subject name
    - `version` (number, required): The version number
  - **delete**: Delete a specific version of a subject's schema
    - `subject` (string, required): The subject name
    - `version` (number, required): The version number

- **compatibility**
  - **get**: Get compatibility setting for a subject
    - `subject` (string, required): The subject name
  - **set**: Set compatibility level for a subject
    - `subject` (string, required): The subject name
    - `level` (string, required): The compatibility level (e.g. BACKWARD, FORWARD, FULL, NONE)

- **types**
  - **list**: List supported schema types (e.g. AVRO, JSON, PROTOBUF)
    - _Parameters_: None 