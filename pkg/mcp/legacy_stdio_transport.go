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
	"bytes"
	"context"
	"io"
	stdlog "log"
	"sync"

	"github.com/mark3labs/mcp-go/server"
	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Run starts the legacy mark3labs MCP server using a go-sdk transport.
func (s *LegacyServer) Run(ctx context.Context, transport sdk.Transport, errLogger *stdlog.Logger) error {
	conn, err := transport.Connect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	stdioServer := server.NewStdioServer(s.MCPServer)
	if errLogger != nil {
		stdioServer.SetErrorLogger(errLogger)
	}

	reader := &jsonrpcReader{ctx: ctx, conn: conn}
	writer := &jsonrpcWriter{ctx: ctx, conn: conn}
	return stdioServer.Listen(ctx, reader, writer)
}

type jsonrpcReader struct {
	ctx  context.Context
	conn sdk.Connection
	buf  []byte
}

func (r *jsonrpcReader) Read(p []byte) (int, error) {
	for len(r.buf) == 0 {
		if r.ctx.Err() != nil {
			return 0, io.EOF
		}

		msg, err := r.conn.Read(r.ctx)
		if err != nil {
			if r.ctx.Err() != nil {
				return 0, io.EOF
			}
			return 0, err
		}

		data, err := jsonrpc.EncodeMessage(msg)
		if err != nil {
			return 0, err
		}
		r.buf = make([]byte, len(data)+1)
		copy(r.buf, data)
		r.buf[len(data)] = '\n'
	}

	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

type jsonrpcWriter struct {
	ctx  context.Context
	conn sdk.Connection
	mu   sync.Mutex
	buf  []byte
}

func (w *jsonrpcWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf = append(w.buf, p...)
	for {
		index := bytes.IndexByte(w.buf, '\n')
		if index < 0 {
			break
		}

		line := bytes.TrimSpace(w.buf[:index])
		w.buf = w.buf[index+1:]
		if len(line) == 0 {
			continue
		}

		msg, err := jsonrpc.DecodeMessage(line)
		if err != nil {
			return 0, err
		}
		if err := w.conn.Write(w.ctx, msg); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}
