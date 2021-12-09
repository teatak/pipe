# pipe

a reverse proxy server

```yaml
riff:
  url: riff://localhost:8630

backend:
  - name: test
    mode: hash #random 随机 roundRobin 轮训 hash 哈希 consistentHash 一致性哈希
    server: 
      - http://test
  - name: riff
    mode: hash
    server:
      - http://127.0.0.1:8610
      - http://127.0.0.1:8610

server:
  - listen: "[::]:80"
    domain:
      - name: _
        location:
            # return static folder
          - path: /static/
            return: "static /Users/yanggang/workspace/github.com/teatak/pipe/static/"
            # return a file
          - path: /favicon.ico
            return: "file /Users/yanggang/workspace/github.com/teatak/pipe/static/favicon.ico"
            # return string
          - path: ~ ^/string$ # regex
            return: "string 200 hello string"
            # redirect
          - path: /https
            return: "redirect 301 https://$host$request_uri"
            # redirect
          - path: /redirect
            return: "redirect 301 /string" 
          - path: /backend
            return: "backend test"
            # default
          - path: /
            return: "string 200 default"
  - listen: "[::]:443"
    ssl: true
    domain:
      - name: dev.x-t.top console.teatak.com
        certFile: "/Users/yanggang/workspace/cert/dev.x-t.top.pem"
        keyFile: "/Users/yanggang/workspace/cert/dev.x-t.top.key"
        location: 
          - path: /
            return: "backend test"
          - path: /ws
            return: "backend riff"
          - path: /api
            return: "backend riff"
          - path: /console
            return: "backend riff"
          - path: /static/
            return: "backend riff"
          - path: /bad
            return: "backend bad"
          - path: /test
            return: "backend test"
      - name: teatak.com www.teatak.com
        certFile: "/Users/yanggang/workspace/cert/teatak.com.pem"
        keyFile: "/Users/yanggang/workspace/cert/teatak.com.key"
        location: 
          - path: /
            return: "backend test"
      - name: _
        certFile: "/Users/yanggang/workspace/cert/teatak.com.pem"
        keyFile: "/Users/yanggang/workspace/cert/teatak.com.key"
        location: 
          - path: /
            return: "string 200 index"
            header: 
              - "Content-Type text/plain; charset=utf-8"
            # proxy: http://console
    
```