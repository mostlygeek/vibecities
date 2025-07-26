package web

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mostlygeek/vibecities/db"
)

type Server struct {
	db     db.Store
	mcp    *server.MCPServer
	engine *gin.Engine
}

func NewServer(db db.Store, mcp *server.MCPServer) *Server {
	r := gin.Default()
	httpMcp := server.NewStreamableHTTPServer(mcp)

	s := &Server{
		db:     db,
		mcp:    mcp,
		engine: r,
	}

	s.engine.GET("/", s.IndexHandler)

	s.engine.POST("/mcp", gin.WrapH(httpMcp))

	// Use a parameterized route instead of wildcard
	s.engine.GET("/:path", func(c *gin.Context) {
		path := "/" + c.Param("path")

		if rec, ok := s.db.Get(path); ok {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, rec.Data)
		} else {
			c.String(http.StatusNotFound, "")
		}
	})

	return s
}

// ServeHTTP implements http.Handler for compatibility.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}
func (s *Server) IndexHandler(c *gin.Context) {
	records := s.db.List()

	html := `<!DOCTYPE html>
<html>
<head>
    <title>MCPCities</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f5f5f5; font-weight: bold; }
        tr:hover { background-color: #f9f9f9; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .timestamp { color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>VibeCities</h1>`

	// Helper function for relative time
	formatRelativeTime := func(t time.Time) string {
		duration := time.Since(t)

		if duration < time.Minute {
			return "just now"
		} else if duration < time.Hour {
			minutes := int(duration.Minutes())
			if minutes == 1 {
				return "1 minute ago"
			}
			return fmt.Sprintf("%d minutes ago", minutes)
		} else if duration < 24*time.Hour {
			hours := int(duration.Hours())
			if hours == 1 {
				return "1 hour ago"
			}
			return fmt.Sprintf("%d hours ago", hours)
		} else if duration < 7*24*time.Hour {
			days := int(duration.Hours() / 24)
			if days == 1 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", days)
		} else if duration < 30*24*time.Hour {
			weeks := int(duration.Hours() / (7 * 24))
			if weeks == 1 {
				return "1 week ago"
			}
			return fmt.Sprintf("%d weeks ago", weeks)
		} else if duration < 365*24*time.Hour {
			months := int(duration.Hours() / (30 * 24))
			if months == 1 {
				return "1 month ago"
			}
			return fmt.Sprintf("%d months ago", months)
		} else {
			years := int(duration.Hours() / (365 * 24))
			if years == 1 {
				return "1 year ago"
			}
			return fmt.Sprintf("%d years ago", years)
		}
	}

	if len(records) == 0 {
		html += `    <p>No pages available yet.</p>`
	} else {
		// Convert map to slice for sorting
		type recordItem struct {
			path    string
			title   string
			created time.Time
			updated time.Time
		}
		var items []recordItem
		for path, record := range records {
			items = append(items, recordItem{
				path:    path,
				title:   record.Title,
				created: record.Created,
				updated: record.Updated,
			})
		}

		// Sort by path
		sort.Slice(items, func(i, j int) bool {
			return items[i].path < items[j].path
		})

		html += `    <table>
        <thead>
            <tr>
                <th>Path</th>
                <th>Title</th>
                <th>Created</th>
                <th>Updated</th>
            </tr>
        </thead>
        <tbody>`

		for _, item := range items {
			html += `            <tr>
                <td><a href="` + item.path + `" target="_blank">` + item.path + `</a></td>
                <td>` + item.title + `</td>
                <td class="timestamp">` + formatRelativeTime(item.created) + `</td>
                <td class="timestamp">` + formatRelativeTime(item.updated) + `</td>
            </tr>`
		}
		html += `        </tbody>
    </table>`
	}

	html += `</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
