# VibeCities

![VibeCities Header](vibecities-header2.webp)

**VibeCities** is bringing back the nostalgia of the personal web page.
- Notepad.exe → LLM
- FTP → MCP
- GeoCities → VibeCities

## About

VibeCities gives your LLM a one step web publishing super power! 

- Single, easy to run binary `./vibecities`.
- Streamable http mcp server at `/mcp`
- Stores web pages in one sqlite file (no cgo necessary)

It started with a question: what MCP server can I build in less than 3 hours? This was my original design spec: 

- an MCP server that CRUDs some resource
- use golang
- keep it super simple
- make it composable with other tools
- do something fun!

## Batteries Included: vibecities.db

- VibeCities comes some super cool sites out of the box
- `./vibecities-darwin-arm64 -db vibecities.db`
- Point your browser at `http://localhost:1337`
  
## Try it with Claude Desktop

> [!NOTE]
> Any MCP client will work. Claude Desktop was picked to make it easy to play with quickly. 

1. Build the vibecities binary: `make all`
1. Start vibecities `./build/vibecities-darwin-arm64` (use appropriate binary for your platform)
1. Set MCP config in Claude Desktop: “Settings → Developer → Edit Config”:

   ```json
   {
     "mcpServers": {
       "vibecities": {
         "command": "npx",
         "args": ["mcp-remote", "http://localhost:1337/mcp", "--allow-http"]
       }
     }
   }
   ```

   -  `mcp-remote` is required as Claude desktop does not support remote transports 
   -  `--allow-http` is required _only_ when vibecities is running somewhere other than localhost.

1. Restart claude desktop (Cmd+r on mac, ctrl+r on windows)
1. Ask claude to search the web and make a new page, for example:
   <img width="1112" height="912" alt="image" src="https://github.com/user-attachments/assets/0b1a2a62-db90-49ee-9916-2203a2ffcb41" />


1. Open `http://localhost:1337/bbs` in your browser and feel the vibes...
   <img width="1112" height="990" alt="image" src="https://github.com/user-attachments/assets/e50a2e08-4f70-4652-93dc-3800f97e0b62" />


## Benchmarks (M1 Macbook Pro)

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

## Tools

- https://github.com/mark3labs/mcp-go .. seems to be the most popular one
- `npx @modelcontextprotocol/inspector` for testing mcp server (localhost:1337/mpc)
