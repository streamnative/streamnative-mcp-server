# Lane: Pulsar Resource Family

This lane adds one coherent Pulsar resource family.

Focus on:

- mapping existing read-only Pulsar admin/client APIs to MCP resources
- keeping large collections behind templates
- returning bounded JSON payloads
- matching existing tool builder naming and feature-gate conventions
- adding focused tests for registration, URI parsing, and handler behavior

Do not introduce mutating behavior. If an existing tool operation is mutating,
it must remain a tool and must not become a resource.
