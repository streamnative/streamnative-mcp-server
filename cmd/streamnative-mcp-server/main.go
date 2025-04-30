// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/streamnative/streamnative-mcp-server/pkg/cmd/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
)

var version = "version"
var commit = "commit"
var date = "date"

func main() {

	// Create configuration options
	configOpts := config.NewConfigOptions()

	// Create root command
	rootCmd := newRootCommand(configOpts)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// newRootCommand creates and returns the root command
func newRootCommand(configOpts *config.Options) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "snmcp",
		Short: "StreamNative Cloud MCP Server for AI agent integration",
		Long: `StreamNative Cloud MCP Server provides resources and tools for AI agents 
to interact with StreamNative Cloud resources and services.`,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
		Version: fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s", version, commit, date),
	}

	// Add global flags
	configOpts.AddFlags(rootCmd)

	o := mcp.NewMcpServerOptions(configOpts)
	o.AddFlags(rootCmd)
	// Add subcommands
	rootCmd.AddCommand(mcp.NewCmdMcpStdioServer(o))

	rootCmd.SetVersionTemplate("{{.Short}}\n{{.Version}}\n")

	return rootCmd
}
