#### pulsar_admin_functions_worker

Unified tool for managing Apache Pulsar Functions Worker resources.

- **resource** (string, required): Type of functions worker resource to access
  - **function_stats**: Get statistics for all functions running on this functions worker
    - No additional parameters required
  - **monitoring_metrics**: Get metrics for monitoring function workers
    - No additional parameters required
  - **cluster**: Get information about the function worker cluster
    - No additional parameters required
  - **cluster_leader**: Get the leader of the worker cluster
    - No additional parameters required
  - **function_assignments**: Get the assignments of functions across the worker cluster
    - No additional parameters required
