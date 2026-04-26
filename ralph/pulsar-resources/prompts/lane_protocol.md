# Lane: Protocol Foundation

This lane builds or fixes the MCP resource protocol surface for Pulsar.

Focus on:

- `resources/list`
- `resources/templates/list`
- `resources/read`
- resource registration
- feature-gate behavior
- missing-session behavior
- URI scheme and parser foundations

Do not attempt to implement every Pulsar resource family in this lane. Add only
the minimum representative resources or templates needed to prove the protocol
surface and unblock family lanes.
