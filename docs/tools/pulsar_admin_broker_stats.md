#### pulsar_admin_broker_stats

Unified tool for retrieving Apache Pulsar broker statistics.

- **monitoring_metrics**
  - **get**: Get broker monitoring metrics
    
- **mbeans**
  - **get**: Get JVM MBeans statistics from broker
    
- **topics**
  - **get**: Get statistics for all topics managed by the broker
    
- **allocator_stats**
  - **get**: Get memory allocator statistics
    - `allocator_name` (string, required): Name of the allocator
    
- **load_report**
  - **get**: Get broker load report 