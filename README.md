# VibeCities

![VibeCities Header](vibecities-header2.webp)

Ever wonder what if Geocities started in 2025? What if the only interface was to vibe code your page? Wonder no more. **VibeCities** is here to bring back the nostalgia!

## About

- written in golang, single binary to run `./vibecities`
- point your browser web pages at `http://localhost:8111/`
- supports streamable http mcp at `/mcp`
- stores web pages in sqlite (no cgo necessary)

## How to use it with claude desktop

1. Click on “Settings → Developer → Edit Config” and place the following block:

   ```json
   {
     "mcpServers": {
       "vibecities": {
         "command": "npx",
         "args": ["mcp-remote", "http://localhost:8111/mcp"]
       }
     }
   }
   ```

1. start vibecities `./vibecities`
1.

---

# Design Spec

- Bring back the 90s web!
- an MCP server for managing a small website of info
- serves web pages at /{somepage}.html
- use golang
- pages served are single page HTML
- keep it super simple

## Tools

- https://github.com/mark3labs/mcp-go .. seems to be the most popular one
- `npx @modelcontextprotocol/inspector` for testing mcp server (localhost:8111/mpc)
