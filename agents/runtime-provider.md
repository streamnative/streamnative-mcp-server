# Embedded Cloud runtime providers

## Interface and ownership

An embedding application supplies `config.CloudProvider` with:

- `Identity`: the configured effective identity used by instructions and whoami.
- `TokenSource`: concurrency-safe control-plane authentication, also used by logs.
- `Transport`: optional TLS/proxy transport; credential headers are supplied by TokenSource.
- `Resolver`: a session-scoped `ClusterResolver`. Organization and identity are fixed
  by the application, not provided by individual tool requests.

`ResolveCluster(ctx, instance, cluster)` returns `ResolvedRuntime`: Pulsar
connection settings and optional Kafka settings, with dynamic TokenSources. It
must validate resource ownership, readiness, endpoints and authorization policy;
never modify the active clients, launch an interactive login, or return raw
client credentials. Persistent credential selection and cache policy are owned
by the embedding application.

`Options.Complete` validates an explicit provider without opening a keyring,
reading ambient credentials or running a key-file flow. Explicit provider mode
cannot be combined with key-file or external modes. Provider configuration is
process-local and excluded from serialization. Modern HTTP remains external-only;
injected Cloud is not a new modern protocol profile.

The standalone adapter still accepts service-account key files and resolves
PulsarInstance OAuth2 status. Its grants and runtime renewal state are process-local,
not shared through an audience-only persistent cache. It does not interpret
application-specific resource annotations or derive application-specific audiences.
Restarting standalone mode may require a new runtime token exchange.

## Generic runtime lifetime

- Both adapters cross the same resolver interface. The MCP kernel prepares
  clients, then atomically publishes a `config.RuntimeBinding` generation.
- Failed resolution or preparation closes candidates and preserves the previous
  generation. Successful switches/reset wait for in-flight tool/resource/prompt
  leases, publish the new generation, and close retired clients.
- Injected resolvers retain token renewal ownership; Pulsar/Kafka/SR/Connect clients
  obtain current tokens at authentication time. The authentication layer adds no
  business-operation retries. Protocol retry policies are unchanged.
- Direct integrations outside MCP middleware must use `AcquireRuntimeBinding`
  and its snapshot for the entire request. Do not retain clients past the lease
  or publish/reset while holding a read lease. External modes have their own lifetime.
- Transports use `Server.Close` for cleanup. Initial selection does not imply
  locking: `--lock-cluster-context` and read-only mode disable context mutation tools.

## Compatibility

The existing standalone key-file entrypoint, context tool schema and external
protocol settings remain supported. Embedding callers need a released dependency
containing this interface before updating their module pin. Local workspace builds
are development validation, not a substitute for a released module.
