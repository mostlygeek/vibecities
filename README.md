# VibeCities

![VibeCities Header](vibecities-header2.webp)

What if Geocities started in 2025? What if the only interface was to vibe code your page? Wonder no more. **VibeCities** is here to bring back the nostalgia!

## About

- written in golang, single binary to run `./vibecities`
- point your browser web pages at `http://localhost:1337/`
- supports streamable http mcp at `/mcp`
- stores web pages in sqlite (no cgo necessary)

## Testing with Claude Desktop!

These instructions use Claude Desktop as an MCP client.

1. Build the binary: `go build -o bin/vibecities cmd/server/main.go`
1. Start vibecities `./bin/vibecities`
1. Set MCP server in Claude Desktop: “Settings → Developer → Edit Config”:

   ```json
   {
     "mcpServers": {
       "vibecities": {
         "command": "npx",
         "args": ["mcp-remote", "http://localhost:1337/mcp"]
       }
     }
   }
   ```

1. Restart claude desktop (Cmd+r on mac, ctrl+r on windows)
1. Ask claude to make you a website, Try this:

   ```
   Create a new page in vibecities under /geocities. The content should history of geocities it's rise and eventual fall, it's affect on culture and impact on the web. I want the design to be a 90s style website. Make it a nice personal geocities webpage. A night time theme. Throw in some funny 90s cultural references. Skip the artifact creation and recapping what you did.
   ```

1. Open `http://localhost:1337/` in your browser

## Fast! (M1 MBP)

```
$ ab -n 1000 -c 5 http://127.0.0.1:1337/modem
This is ApacheBench, Version 2.3 <$Revision: 1913912 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
...

Server Software:
Server Hostname:        127.0.0.1
Server Port:            1337

Document Path:          /modem
Document Length:        10675 bytes

Concurrency Level:      5
Time taken for tests:   0.106 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      10771000 bytes
HTML transferred:       10675000 bytes
Requests per second:    9458.68 [#/sec] (mean)
Time per request:       0.529 [ms] (mean)
Time per request:       0.106 [ms] (mean, across all concurrent requests)
Transfer rate:          99491.64 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       2
Processing:     0    0   0.4      0       7
Waiting:        0    0   0.4      0       7
Total:          0    1   0.5      0       8
ERROR: The median and mean for the total time are more than twice the standard
       deviation apart. These results are NOT reliable.

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      1
  80%      1
  90%      1
  95%      1
  98%      2
  99%      2
 100%      8 (longest request)
```

---

## Design Specs

(initial ideas for how it should work)

- Bring back the 90s web!
- an MCP server for managing a small website of info
- serves web pages at /{somepage}.html
- use golang
- pages served are single page HTML
- keep it super simple

## Tools

- https://github.com/mark3labs/mcp-go .. seems to be the most popular one
- `npx @modelcontextprotocol/inspector` for testing mcp server (localhost:1337/mpc)
