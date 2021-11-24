# pipe

a reverse proxy server

```yaml
riff:
  url: riff://localhost:8610

endpoint:
  api:
    method: ipHash #robin round hash hashRing
    servers: 
      - http://api
  console:
    method: ipHash
    servers:
      - http://console
    
servers:
  - listen: "[::]:443 ssl"
    domain:
      - name: api.teatak.com
        certFile: "/data/cert/api.teatak.com"
        keyFile: "/data/cert/api.teatak.com"
        location: 
          - path: /
            to: http://api
          - path: /api
            to: http://api
      - name: _
        location: 
          - path: /
            to: http://console
```