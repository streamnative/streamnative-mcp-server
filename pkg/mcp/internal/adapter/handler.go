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

package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolHandlerFuncV1 represents the mark3labs/mcp-go handler signature.
// Existing handlers use this signature:
// func(ctx context.Context, request CallToolRequest) (*CallToolResult, error)
type ToolHandlerFuncV1 interface{}

// AdaptHandlerV1ToV2 adapts a mark3labs-style handler to go-sdk's ToolHandler.
//
// The key differences:
// - mark3labs: func(ctx, request) -> (result, error)
// - go-sdk:     func(ctx, request) -> (result, error)
//
// Both have similar signatures, but the types are from different packages.
// This adapter handles the type conversion.
func AdaptHandlerV1ToV2(handler ToolHandlerFuncV1) mcpsdk.ToolHandler {
	// Try to assert as go-sdk handler first
	if h, ok := handler.(func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error)); ok {
		return h
	}
	// If not a go-sdk handler, return an error handler
	// The actual migration will happen when we update tool builders
	return func(ctx context.Context, request *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{
				&mcpsdk.TextContent{
					Text: "handler needs migration to go-sdk signature",
				},
			},
			IsError: true,
		}, nil
	}
}

// NewTextResult creates a text result for go-sdk.
func NewTextResult(text string) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{
				Text: text,
			},
		},
	}
}

// NewErrorResult creates an error result for go-sdk.
func NewErrorResult(format string, args ...any) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{
				Text: fmt.Sprintf(format, args...),
			},
		},
		IsError: true,
	}
}

// AdaptToolDefV1ToV2 converts a mark3labs Tool definition to go-sdk Tool.
// This handles the conversion of tool metadata and input schema.
func AdaptToolDefV1ToV2(toolV1 interface{}) *mcpsdk.Tool {
	// During migration, we'll need to convert tool definitions
	// For now, tool builders will need to be updated to use go-sdk types
	return nil
}

// RequireString extracts a required string argument from the request.
func RequireString(request *mcpsdk.CallToolRequest, name string) (string, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", name)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s is not of type string", name)
	}
	return str, nil
}

// RequireOptionalString extracts an optional string argument from the request.
func RequireOptionalString(request *mcpsdk.CallToolRequest, name string) (string, bool, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return "", false, nil
	}
	str, ok := val.(string)
	if !ok {
		return "", false, fmt.Errorf("parameter %s is not of type string", name)
	}
	return str, true, nil
}

// GetArgumentsMap extracts all arguments as a map from the request.
func GetArgumentsMap(request *mcpsdk.CallToolRequest) (map[string]any, error) {
	if request.Params.Arguments == nil || len(request.Params.Arguments) == 0 {
		return make(map[string]any), nil
	}
	var args map[string]any
	if err := json.Unmarshal(request.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}
	return args, nil
}

// GetString extracts an optional string argument from the request.
func GetString(request *mcpsdk.CallToolRequest, name string, defaultValue string) string {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return defaultValue
	}
	val, ok := args[name]
	if !ok {
		return defaultValue
	}
	str, ok := val.(string)
	if !ok {
		return defaultValue
	}
	return str
}

// GetFloat extracts an optional float argument from the request.
func GetFloat(request *mcpsdk.CallToolRequest, name string, defaultValue float64) float64 {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return defaultValue
	}
	val, ok := args[name]
	if !ok {
		return defaultValue
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return defaultValue
	}
}

// GetInt extracts an optional int argument from the request.
func GetInt(request *mcpsdk.CallToolRequest, name string, defaultValue int) int {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return defaultValue
	}
	val, ok := args[name]
	if !ok {
		return defaultValue
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultValue
	}
}

// GetBool extracts an optional bool argument from the request.
func GetBool(request *mcpsdk.CallToolRequest, name string, defaultValue bool) bool {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return defaultValue
	}
	val, ok := args[name]
	if !ok {
		return defaultValue
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	default:
		return defaultValue
	}
}

// GetStringSlice extracts an optional string slice argument from the request.
func GetStringSlice(request *mcpsdk.CallToolRequest, name string, defaultValue []string) []string {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return defaultValue
	}
	val, ok := args[name]
	if !ok {
		return defaultValue
	}
	// Try to convert to []string
	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	default:
		return defaultValue
	}
}

// RequireInt extracts a required int argument from the request.
func RequireInt(request *mcpsdk.CallToolRequest, name string) (int, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return 0, fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return 0, fmt.Errorf("missing required parameter: %s", name)
	}
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("parameter %s is not of type int", name)
	}
}

// RequireFloat extracts a required float argument from the request.
func RequireFloat(request *mcpsdk.CallToolRequest, name string) (float64, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return 0, fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return 0, fmt.Errorf("missing required parameter: %s", name)
	}
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("parameter %s is not of type float", name)
	}
}

// RequireStringSlice extracts a required string slice argument from the request.
func RequireStringSlice(request *mcpsdk.CallToolRequest, name string) ([]string, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: %s", name)
	}
	switch v := val.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				return nil, fmt.Errorf("parameter %s contains non-string element", name)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("parameter %s is not of type []string", name)
	}
}

// RequireObject extracts a required object argument from the request.
func RequireObject(request *mcpsdk.CallToolRequest, name string) (map[string]any, error) {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}
	val, ok := args[name]
	if !ok {
		return nil, fmt.Errorf("missing required parameter: %s", name)
	}
	obj, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("parameter %s is not of type object", name)
	}
	return obj, nil
}

// GetObject extracts an optional object argument from the request.
func GetObject(request *mcpsdk.CallToolRequest, name string) map[string]any {
	args, err := GetArgumentsMap(request)
	if err != nil {
		return nil
	}
	val, ok := args[name]
	if !ok {
		return nil
	}
	obj, ok := val.(map[string]any)
	if !ok {
		return nil
	}
	return obj
}
