package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mostlygeek/vibecities/db"
)

func NewMCPServer(db db.Store, enableLoadPath bool) *server.MCPServer {
	srv := server.NewMCPServer(
		"mcpcities",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	// Add page list tool
	listTool := mcp.NewTool("page_list",
		mcp.WithDescription("List all web pages in the database"),
	)
	srv.AddTool(listTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		records := db.List()

		jsonData, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling web pages: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	})

	// Add page set tool
	setTool := mcp.NewTool("page_set",
		mcp.WithDescription("Set the HTML source of a web page"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("url path for the web page"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("A Short title of the web page"),
		),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("HTML source to store"),
		),
	)
	srv.AddTool(setTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, err := request.RequireString("data")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := db.Set(path, title, data); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error setting web page: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully set web page at path: %s", path)), nil
	})

	// Add page get tool
	getTool := mcp.NewTool("page_get",
		mcp.WithDescription("Get a web page from the database"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("url path for the web page"),
		),
	)
	srv.AddTool(getTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		record, ok := db.Get(path)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Web page not found at path: %s", path)), nil
		}

		return mcp.NewToolResultText(record.Data), nil
	})

	deleteTool := mcp.NewTool("page_delete",
		mcp.WithDescription("Delete a web page from the database"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("url path for the web page to delete"),
		),
	)
	srv.AddTool(deleteTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := db.Delete(path); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error deleting web page: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted web page at path: %s", path)), nil
	})

	if enableLoadPath {
		log.Println("WARNING: loadPath tool enabled - this is a security risk if not used in a trusted environment!")

		// Add loadPath tool
		loadPathTool := mcp.NewTool("page_set_from_file",
			mcp.WithDescription("Load a file from the local filesystem and store it as a web page in the database"),
			mcp.WithString("filepath",
				mcp.Required(),
				mcp.Description("Full local file system path to read from"),
			),
			mcp.WithString("path",
				mcp.Required(),
				mcp.Description("url path for the web page in the database"),
			),
		)

		srv.AddTool(loadPathTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Check if request is from localhost

			filepathParam, err := request.RequireString("filepath")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			path, err := request.RequireString("path")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Read file content
			content, err := os.ReadFile(filepathParam)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), nil
			}

			// Get title from filename
			title := filepath.Base(filepathParam)
			if idx := strings.LastIndex(title, "."); idx != -1 {
				title = title[:idx]
			}

			// Store in database
			if err := db.Set(path, title, string(content)); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Error storing web page: %v", err)), nil
			}

			return mcp.NewToolResultText(fmt.Sprintf("Successfully loaded file '%s' to path: %s", filepathParam, path)), nil
		})
	}

	return srv
}
