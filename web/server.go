package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mostlygeek/mcpcities/db"
)

type Server struct {
	db     db.Store
	engine *gin.Engine
}

func NewServer(db db.Store) *Server {
	r := gin.Default()
	mcpSrv := NewMCPServer()
	httpMcp := server.NewStreamableHTTPServer(mcpSrv)

	s := &Server{
		db:     db,
		engine: r,
	}

	s.engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the index page!")
	})

	s.engine.POST("/mcp", gin.WrapH(httpMcp))

	// Use a parameterized route instead of wildcard
	s.engine.GET("/:path", func(c *gin.Context) {
		path := "/" + c.Param("path")

		if rec, ok := s.db.Get(path); ok {
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
