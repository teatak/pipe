package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/teatak/pipe/loadbalance"
	"github.com/teatak/pipe/sections"
	"github.com/teatak/riff/api"
)

var proxys = make(map[string]*httputil.ReverseProxy)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func GetClientIp(req *http.Request) string {
	ipAddress := req.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = req.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = req.RemoteAddr
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ipAddress)); err == nil {
		return ip
	} else {
		return ""
	}
}

func NewProxy(backend *sections.Backend) (*httputil.ReverseProxy, error) {
	if proxys[backend.Name] != nil {
		return proxys[backend.Name], nil
	} else {
		director := func(req *http.Request) {
			nodes := []string{}
			_url := ""
			//load from riff
			if backend.Riff != "" {
				arr := strings.Split(backend.Riff, "@")
				serviceName := arr[0]
				riffUrl := arr[1]
				client, err := api.RiffClient(riffUrl)
				if err != nil {
					return
				}
				service := client.Services(serviceName, api.StateAlive)
				for _, node := range service.NestNodes {
					nodeString := "http://" + node.IP + ":" + strconv.Itoa(node.Port)
					nodes = append(nodes, nodeString)
				}
			} else {
				nodes = append(nodes, backend.Server...)
			}

			ip := GetClientIp(req)
			switch backend.Mode {
			case "random":
				_url, _ = loadbalance.Random(nodes)
			case "roundRobin":
				_url, _ = loadbalance.Random(nodes)
			case "hash":
				_url, _ = loadbalance.Hash(ip, nodes)
			case "consistentHash":
				_url, _ = loadbalance.ConsistentHash(ip, nodes)
			default:
				//roundRobin
				_url, _ = loadbalance.Random(nodes)

			}
			if _url == "" {
				return
			}
			target, err := url.Parse(_url)
			if err != nil {
				fmt.Println(err)
			}
			targetQuery := target.RawQuery
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
			modifyRequest(req)
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ModifyResponse = modifyResponse()
		proxy.ErrorHandler = errorHandler(backend)

		proxys[backend.Name] = proxy
		return proxy, nil
	}
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "pipe")
}

func errorHandler(backend *sections.Backend) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		server.Error502("backend "+backend.Name, w)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "pipe")
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
