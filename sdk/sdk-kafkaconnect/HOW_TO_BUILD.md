1. Update the `kafka-connect-admin.json` from https://github.com/streamnative/function-mesh-worker-service/blob/master/openapi/kafka-connect-admin.json
2. Using following command in `pkg/kafkaconnect` to generate the client sdk

```
openapi-generator-cli generate \
  -i kafka-connect-admin.json \
  -g go \
  -o ./ \
  --additional-properties=packageName=kafkaconnect,isGoSubmodule=true --git-repo-id cloud-cli/pkg --git-user-id streamnative
```