package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mostlygeek/mcpcities/db"
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
			c.Header("Content-Type", "text/html")
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
    <title>Available Pages</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        ul { list-style-type: none; padding: 0; }
        li { margin: 10px 0; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>Available Pages</h1>`

	if len(records) == 0 {
		html += `    <p>No pages available yet.</p>`
	} else {
		html += `    <ul>`
		for path := range records {
			html += `        <li><a href="` + path + `">` + path + `</a></li>`
		}
		html += `    </ul>`
	}

	html += `</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
