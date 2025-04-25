// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

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
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
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
