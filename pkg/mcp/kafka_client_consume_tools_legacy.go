// Copyright 2025 StreamNative
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders"
	kafkabuilders "github.com/streamnative/streamnative-mcp-server/pkg/mcp/builders/kafka"
)

// KafkaClientAddConsumeToolsLegacy adds Kafka client consume tools to the legacy MCP server.
func KafkaClientAddConsumeToolsLegacy(s *server.MCPServer, readOnly bool, logrusLogger *logrus.Logger, features []string) {
	builder := kafkabuilders.NewKafkaConsumeLegacyToolBuilder()
	config := builders.ToolBuildConfig{
		ReadOnly: readOnly,
		Features: features,
		Options: map[string]interface{}{
			"logger": logrusLogger,
		},
	}

	tools, err := builder.BuildTools(context.Background(), config)
	if err != nil {
		return
	}

	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}
}
