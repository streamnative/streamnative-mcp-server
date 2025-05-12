#### pulsar_client_produce

Produce messages to a Pulsar topic. This tool allows you to send messages to a specified Pulsar topic with various options to control message format, batching, and properties.

- **pulsar_client_produce**
  - `topic` (string, required): Topic to produce to
  - `messages` (array, required): Messages to send. Specify multiple times for multiple messages
  - `num-produce` (number, optional): Number of times to send message(s) (default: 1)
  - `rate` (number, optional): Rate (in msg/sec) at which to produce, 0 means produce as fast as possible (default: 0)
  - `disable-batching` (boolean, optional): Disable batch sending of messages (default: false)
  - `chunking` (boolean, optional): Should split the message and publish in chunks if message size is larger than allowed max size (default: false)
  - `separator` (string, optional): Character to split messages string on (default: none)
  - `properties` (array, optional): Properties to add, key=value format. Specify multiple times for multiple properties
  - `key` (string, optional): Partitioning key to add to each message 