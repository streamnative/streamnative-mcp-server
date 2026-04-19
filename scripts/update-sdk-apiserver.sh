#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=$(cd -- "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
SDK_DIR="${ROOT_DIR}/sdk/sdk-apiserver"
CLOUD_API_SERVER_DIR="${CLOUD_API_SERVER_DIR:-${HOME}/go/src/github.com/streamnative/cloud-api-server}"
GO_BIN="${GO_BIN:-$(command -v go)}"
PYTHON_BIN="${PYTHON_BIN:-python3}"
DOCKER_BIN="${DOCKER_BIN:-docker}"
GENERATOR_IMAGE="${OPENAPI_GENERATOR_IMAGE:-openapitools/openapi-generator-cli:v6.2.0}"

if [[ ! -d "${SDK_DIR}" ]]; then
  echo "sdk/sdk-apiserver not found at ${SDK_DIR}" >&2
  exit 1
fi

if [[ ! -d "${CLOUD_API_SERVER_DIR}" ]]; then
  echo "cloud-api-server repo not found at ${CLOUD_API_SERVER_DIR}" >&2
  echo "Set CLOUD_API_SERVER_DIR to a local checkout that contains pkg/openapi." >&2
  exit 1
fi

WORK_DIR=$(mktemp -d)
EXPORTER_DIR=""
cleanup() {
  rm -rf "${WORK_DIR}"
  if [[ -n "${EXPORTER_DIR}" && -d "${EXPORTER_DIR}" ]]; then
    rm -rf "${EXPORTER_DIR}"
  fi
}
trap cleanup EXIT

EXPORTER_DIR=$(mktemp -d "${CLOUD_API_SERVER_DIR}/.tmp-export-cloud-api-instance-schemas-XXXXXX")

cat > "${EXPORTER_DIR}/main.go" <<'EOF'
package main

import (
	"encoding/json"
	"os"

	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	openapibuilder3 "k8s.io/kube-openapi/pkg/builder3"

	"github.com/streamnative/cloud-api-server/pkg/apiserver-builder-alpha/builders"
	generatedopenapi "github.com/streamnative/cloud-api-server/pkg/openapi"
)

func main() {
	cfg := genericapiserver.DefaultOpenAPIV3Config(
		generatedopenapi.GetOpenAPIDefinitions,
		openapinamer.NewDefinitionNamer(builders.Scheme),
	)
	schemas, err := openapibuilder3.BuildOpenAPIDefinitionsForResources(
		cfg,
		"github.com/streamnative/cloud-api-server/pkg/apis/cloud/v1alpha1.Instance",
		"github.com/streamnative/cloud-api-server/pkg/apis/cloud/v1alpha1.InstanceList",
		"github.com/streamnative/cloud-api-server/pkg/apis/cloud/v1alpha1.InstanceSpec",
		"github.com/streamnative/cloud-api-server/pkg/apis/cloud/v1alpha1.InstanceStatus",
	)
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schemas); err != nil {
		panic(err)
	}
}
EOF

(cd "${CLOUD_API_SERVER_DIR}" && env -u GOROOT "${GO_BIN}" run "./$(basename "${EXPORTER_DIR}")") \
  > "${WORK_DIR}/instance-schemas.json"

"${PYTHON_BIN}" - "${SDK_DIR}/swagger.json" "${WORK_DIR}/instance-schemas.json" "${WORK_DIR}/swagger.json" <<'PY'
import copy
import json
import sys

source_path, schemas_path, output_path = sys.argv[1:4]

with open(source_path, "r", encoding="utf-8") as source_file:
    document = json.load(source_file)
with open(schemas_path, "r", encoding="utf-8") as schemas_file:
    instance_schemas = json.load(schemas_file)

components = document.setdefault("components", {}).setdefault("schemas", {})
components.update(instance_schemas)

path_map = {
    "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/pulsarinstances":
        "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/instances",
    "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/pulsarinstances/{name}":
        "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/instances/{name}",
    "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/pulsarinstances/{name}/status":
        "/apis/cloud.streamnative.io/v1alpha1/namespaces/{namespace}/instances/{name}/status",
    "/apis/cloud.streamnative.io/v1alpha1/pulsarinstances":
        "/apis/cloud.streamnative.io/v1alpha1/instances",
    "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/pulsarinstances":
        "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/instances",
    "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/pulsarinstances/{name}":
        "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/instances/{name}",
    "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/pulsarinstances/{name}/status":
        "/apis/cloud.streamnative.io/v1alpha1/watch/namespaces/{namespace}/instances/{name}/status",
    "/apis/cloud.streamnative.io/v1alpha1/watch/pulsarinstances":
        "/apis/cloud.streamnative.io/v1alpha1/watch/instances",
}

replacements = [
    ("PulsarInstance", "Instance"),
    ("pulsarinstances", "instances"),
    ("pulsar instance", "instance"),
    ("Pulsar instance", "Instance"),
    ("pulsar instances", "instances"),
    ("Pulsar instances", "Instances"),
]

def rewrite_string(value: str) -> str:
    for old, new in replacements:
        value = value.replace(old, new)
    return value

def rewrite(value):
    if isinstance(value, dict):
        return {key: rewrite(item) for key, item in value.items()}
    if isinstance(value, list):
        return [rewrite(item) for item in value]
    if isinstance(value, str):
        return rewrite_string(value)
    return value

for old_path, new_path in path_map.items():
    document["paths"][new_path] = rewrite(copy.deepcopy(document["paths"][old_path]))

with open(output_path, "w", encoding="utf-8") as output_file:
    json.dump(document, output_file, indent=2)
    output_file.write("\n")
PY

cp "${WORK_DIR}/swagger.json" "${WORK_DIR}/swagger.json.unprocessed"

"${DOCKER_BIN}" run --rm \
  -u "$(id -u):$(id -g)" \
  -v "${WORK_DIR}:/local" \
  "${GENERATOR_IMAGE}" generate \
  -i /local/swagger.json \
  -g go \
  -o /local/out \
  --git-user-id GIT_USER_ID \
  --git-repo-id GIT_REPO_ID \
  --additional-properties packageName=sncloud,packageVersion=1.0.0,enumClassPrefix=true

cp "${WORK_DIR}/swagger.json" "${WORK_DIR}/out/swagger.json"
cp "${WORK_DIR}/swagger.json.unprocessed" "${WORK_DIR}/out/swagger.json.unprocessed"
cp "${SDK_DIR}/settings" "${WORK_DIR}/out/settings"

rsync -a --delete "${WORK_DIR}/out/" "${SDK_DIR}/"
