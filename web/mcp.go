package web

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/mostlygeek/vibecities/db"
)

type PageSetParams struct {
	Path  string `json:"path"`
	Title string `json:"title"`
	Data  string `json:"data"`
}

type PageGetParams struct {
	Path string `json:"path"`
}

type PageDeleteParams struct {
	Path string `json:"path"`
}

type PageSetFromFileParams struct {
	Filepath string `json:"filepath"`
	Path     string `json:"path"`
}

type PageEditParams struct {
	Path    string `json:"path"`
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

func NewMCPServer(db db.Store) *mcp.Server {
	srv := mcp.NewServer("mcpcities", "0.0.1", nil)

	// Add page list tool
	listTool := mcp.NewServerTool("page_list", "List all web pages in the database",
		func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[struct{}]) (*mcp.CallToolResultFor[string], error) {
			records := db.List()
			jsonData, err := json.MarshalIndent(records, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error marshaling web pages: %v", err)
			}
			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(jsonData)},
				},
			}, nil
		})
	srv.AddTools(listTool)

	// Add page set tool
	setTool := mcp.NewServerTool("page_set", "Set the HTML source of a web page",
		func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[PageSetParams]) (*mcp.CallToolResultFor[string], error) {
			args := params.Arguments
			if err := db.Set(args.Path, args.Title, args.Data); err != nil {
				return nil, fmt.Errorf("error setting web page: %v", err)
			}

			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Successfully set web page at path: %s", args.Path)},
				},
			}, nil
		},
		mcp.Input(
			mcp.Property("path", mcp.Required(true), mcp.Description("url path for the web page")),
			mcp.Property("title", mcp.Required(true), mcp.Description("A short title of the web page")),
			mcp.Property("data", mcp.Required(true), mcp.Description("HTML source to store")),
		))
	srv.AddTools(setTool)

	// Add page get tool
	getTool := mcp.NewServerTool("page_get", "Get a web page from the database",
		func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[PageGetParams]) (*mcp.CallToolResultFor[string], error) {
			args := params.Arguments
			record, ok := db.Get(args.Path)
			if !ok {
				return nil, fmt.Errorf("web page not found at path: %s", args.Path)
			}

			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: record.Data},
				},
			}, nil
		},
		mcp.Input(
			mcp.Property("path", mcp.Required(true), mcp.Description("url path for the web page")),
		))
	srv.AddTools(getTool)

	// Add page delete tool
	deleteTool := mcp.NewServerTool("page_delete", "Delete a web page from the database",
		func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[PageDeleteParams]) (*mcp.CallToolResultFor[string], error) {
			args := params.Arguments
			if err := db.Delete(args.Path); err != nil {
				return nil, fmt.Errorf("error deleting web page: %v", err)
			}

			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Successfully deleted web page at path: %s", args.Path)},
				},
			}, nil
		},
		mcp.Input(
			mcp.Property("path", mcp.Required(true), mcp.Description("url path for the web page to delete")),
		))
	srv.AddTools(deleteTool)

	// Add page edit tool
	editTool := mcp.NewServerTool("page_edit", "Edit webpage by searching and replacing text. Searches for exact text match and replaces it.",
		func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[PageEditParams]) (*mcp.CallToolResultFor[string], error) {
			args := params.Arguments

			// Get current page content
			dbRecord, found := db.Get(args.Path)
			if !found {
				return nil, fmt.Errorf("page not found at path: %s", args.Path)
			}

			// Check if old text exists in the body
			if !strings.Contains(dbRecord.Data, args.OldText) {
				return nil, fmt.Errorf("text to replace not found in page body")
			}

			// Replace the text
			newBody := strings.Replace(dbRecord.Data, args.OldText, args.NewText, 1)

			// Save the updated page
			if err := db.Set(args.Path, dbRecord.Title, newBody); err != nil {
				return nil, fmt.Errorf("error updating web page: %v", err)
			}

			return &mcp.CallToolResultFor[string]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Page edited successfully"},
				},
			}, nil
		},
		mcp.Input(
			mcp.Property("path", mcp.Required(true), mcp.Description("url path for the web page to edit")),
			mcp.Property("old_text", mcp.Required(true), mcp.Description("exact text to find and replace")),
			mcp.Property("new_text", mcp.Required(true), mcp.Description("new text to replace the old text with")),
		))
	srv.AddTools(editTool)

	return srv
}
