# pipe

a reverse proxy server

```yaml
riff:
  url: riff://localhost:8610

endpoint:
  api:
    method: hash #robin round hash hashRing
    servers: 
      - http://api
  console:
    method: hash
    servers:
      - http://console
      
server:
  - listen: "[::]:80"
    domain:
      - name: _
        location:
            # return static folder
          - path: /static/
            return: "200 static /Users/yanggang/workspace/github.com/teatak/pipe/static/"

            # return a file
          - path: /favicon.ico
            return: "200 file /Users/yanggang/workspace/github.com/teatak/pipe/static/favicon.ico"

            # return string
          - path: ~ ^/string$ # regex
            return: "200 string hello string"

            # redirect
          - path: /https
            return: "301 https://$host$request_uri"

            # redirect
          - path: /redirect
            return: "301 /string"
            
  - listen: "[::]:443"
    ssl: true
    domain:
      - name: api.teatak.com
        certFile: "/Users/yanggang/workspace/cert/api.teatak.com.pem"
        keyFile: "/Users/yanggang/workspace/cert/api.teatak.com.key"
        location: 
          - path: /
            to: http://api
          - path: /api
            to: http://api
      - name: console.teatak.com
        certFile: "/Users/yanggang/workspace/cert/console.teatak.com.pem"
        keyFile: "/Users/yanggang/workspace/cert/console.teatak.com.key"
        location: 
          - path: /
            to: http://console
      - name: _
        certFile: "/Users/yanggang/workspace/cert/teatak.com.pem"
        keyFile: "/Users/yanggang/workspace/cert/teatak.com.key"
        location: 
          - path: /
            return: "200 text"
            header: 
              - "Content-Type text/plain; charset=utf-8"
            # to: http://console
  
```