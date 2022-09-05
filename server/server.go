package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/teatak/cart"
	"github.com/teatak/pipe/sections"
)

const errorServerPrefix = "[ERR]  pipe.server: "
const infoServerPrefix = "[INFO] pipe.server: "

var server *Server

type Server struct {
	Logger       *log.Logger
	ShutdownCh   chan struct{}
	httpServers  []*http.Server
	shutdown     bool
	shutdownLock sync.Mutex
}

func NewServer() (*Server, error) {

	shutdownCh := make(chan struct{})

	server = &Server{
		ShutdownCh:  shutdownCh,
		shutdown:    false,
		httpServers: []*http.Server{},
	}

	logOutput := io.MultiWriter(os.Stderr)
	server.Logger = log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds)

	sections.Load()
	go server.setupCart()

	return server, nil
}

func (s *Server) handleDomain(domain *sections.Domain, c *cart.Context, n cart.Next) {
	if domain != nil {
		var matchLocation *sections.Location
		var defaultLocation *sections.Location
		match := false
		regex := false
		for _, location := range domain.Location {
			path := location.Path
			//match regex path
			if location.Path[0] == '~' {
				path = strings.TrimSpace(location.Path[1:])
				regex = true
			}
			if regex {
				match, _ = regexp.MatchString(path, c.Request.RequestURI)
			} else {
				if path == "/" {
					defaultLocation = location
				} else {
					//split by ; or ,
					arr := strings.FieldsFunc(path, func(r rune) bool {
						return r == ';' || r == ','
					})
					for _, _p := range arr {
						if strings.HasPrefix(c.Request.RequestURI, _p) {
							match = true
							break
						}
					}
				}
			}
			if match {
				matchLocation = location
				break
			}
		}

		if matchLocation == nil && defaultLocation != nil {
			matchLocation = defaultLocation
		}

		if matchLocation != nil {
			location := matchLocation
			if len(location.Header) != 0 {
				for _, header := range location.Header {
					arr := strings.SplitN(header, " ", 2)
					key := arr[0]
					value := arr[1]
					c.Response.Header().Set(key, value)
				}
			}
			if location.Return != "" {
				arr := strings.SplitN(location.Return, " ", 2)
				_type := arr[0]
				content := arr[1]
				// if len(arr) > 2 {
				// 	code, _ = strconv.Atoi(arr[1])
				// 	content = arr[2]
				// }
				errorText := ""
				switch _type {
				case "redirect":
					_arr := strings.SplitN(content, " ", 2)
					if len(_arr) > 1 {
						code, _ := strconv.Atoi(_arr[0])
						path := _arr[1]
						path = strings.ReplaceAll(path, "$host", c.Request.Host)
						path = strings.ReplaceAll(path, "$request_uri", c.Request.RequestURI)
						c.Redirect(code, path)
					} else {
						errorText = "redirect format error"
					}
				case "json":
					header := c.Response.Header()
					header["Content-Type"] = []string{"application/json; charset=utf-8"}
					_arr := strings.SplitN(content, " ", 2)
					if len(_arr) > 1 {
						code, _ := strconv.Atoi(_arr[0])
						jsonText := _arr[1]
						c.Response.WriteHeader(code)
						c.Response.Write([]byte(jsonText))
					} else {
						errorText = "json format error"
					}

				case "html":
					header := c.Response.Header()
					header["Content-Type"] = []string{"text/html; charset=utf-8"}
					_arr := strings.SplitN(content, " ", 2)
					if len(_arr) > 1 {
						code, _ := strconv.Atoi(_arr[0])
						htmlText := _arr[1]
						c.Response.WriteHeader(code)
						c.Response.Write([]byte(htmlText))
					} else {
						errorText = "html format error"
					}
				case "string":
					header := c.Response.Header()
					header["Content-Type"] = []string{"text/plain; charset=utf-8"}
					_arr := strings.SplitN(content, " ", 2)
					if len(_arr) > 1 {
						code, _ := strconv.Atoi(_arr[0])
						stringText := _arr[1]
						c.String(code, stringText)
					} else {
						errorText = "string format error"
					}
				case "file":
					c.File(content)
				case "static":
					_arr := strings.SplitN(content, " ", 2)
					relativePath := ""
					fallback := ""
					if len(_arr) > 1 {
						relativePath = _arr[0]
						fallback = _arr[1]
					} else {
						relativePath = _arr[0]
						fallback = ""
					}
					path := location.Path
					//match regex path
					if location.Path[0] == '~' {
						path = strings.TrimSpace(location.Path[1:])
					}
					c.Static(relativePath, path, true, fallback)
				case "backend":
					s.handleBackend(content, c, n)
				default:
					s.Error502(domain.Name+" return type "+_type, c.Response)
				}

				if errorText != "" {
					s.Error502(domain.Name+" "+errorText, c.Response)

				}
			}
		}
	} else {
		s.Error502("domain "+c.Request.Host, c.Response)
	}
	n()
}

func (s *Server) Error502(content string, resp http.ResponseWriter) {

	tplString := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <style>
    .title {
		display: block;
    	font-size: 2em;
    	font-weight: bold;
    	margin: 22px 0;
    }
    .content {
        margin: 10px 0;
        padding: 10px;
        background: linen;
        font-size: 14px;
    	line-height: 150%;
    }
    .content pre {
    	padding: 0;
    	margin: 0;
        white-space: pre-wrap;
		white-space: -moz-pre-wrap;
		white-space: -pre-wrap;
		white-space: -o-pre-wrap;
		word-wrap: break-word;
		word-break: break-all;
    }
    footer {
    	text-align: center;
		margin: 20px 0;
    	padding: 10px 0;
	}
    footer span {

    }
	footer a {
		display: inline-block;
		vertical-align: middle;
    }
	.center {
		margin-top: 16px;
		display: flex;
		justify-content: center;
    }
	pre {
		font-size: 10pt;
    	font-family: "Courier New", Monospace;
    	white-space: pre;
    }
    </style>
</head>
<body>
<div class="title">{{.Title}}</div>
<div class="content">{{.Content}}</div>
<footer>
	<span>powered by pipe</span>
	<a target="_blank" href="https://github.com/teatak/pipe"><svg width="22" height="22" class="octicon octicon-mark-github" viewBox="0 0 16 16" version="1.1" aria-hidden="true"><path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path></svg></a>
</footer>
</body>
</html>
	`

	tpl, err := template.New("ErrorHTML").Parse(tplString)
	if err != nil {
		panic(err)
	}

	htmlContent := fmt.Sprintf("The <b style='color:red'>%v</b> error", content)

	obj := cart.H{
		"Title":   "502 Bad Gateway",
		"Content": template.HTML(htmlContent),
	}
	header := resp.Header()
	header["Content-Type"] = []string{"text/html; charset=utf-8"}
	resp.WriteHeader(502)
	tpl.Execute(resp, obj)
	// c.Render(code, render.HTML{Template: tpl, Data: obj})
	s.Logger.Printf(errorServerPrefix+"error %v\n", content)
}

func (s *Server) handleBackend(backendName string, c *cart.Context, n cart.Next) {
	//find backend
	var backend *sections.Backend
	for _, temp := range sections.Backends {
		if temp.Name == backendName {
			backend = temp
			break
		}
	}

	if backend == nil {
		s.Error502("backend "+backendName, c.Response)
	} else {
		//find
		proxy, err := NewProxy(backend)
		if err != nil {
			s.Logger.Printf(errorServerPrefix+"error proxy %v %v\n", backendName, err)
		} else {
			proxy.ServeHTTP(c.Response, c.Request)
		}
	}
}

func (s *Server) setupCart() error {

	for _, sv := range sections.Servers {
		cart.SetMode(cart.ReleaseMode)
		r := cart.New()
		tempDomain := sv.Domain
		//构造
		r.Use("/", func(c *cart.Context, n cart.Next) {
			//timer
			since := time.Now()
			defer func() {
				d := time.Since(since)
				s.Logger.Printf(infoServerPrefix+"track %v%v %v \n", c.Request.Host, c.Request.RequestURI, d)
			}()
			findHost := false
			var defaultDoamin *sections.Domain
			for _, domain := range tempDomain {
				if domain.Name == "_" {
					defaultDoamin = domain
				} else {
					domains := strings.Split(domain.Name, " ")
					match := false
					for _, temp := range domains {
						if strings.EqualFold(temp, c.Request.Host) {
							match = true
						}
						if match {
							findHost = true
							s.handleDomain(domain, c, n)
						}
					}
					if findHost {
						break //if match then break loop
					}
				}
			}
			if !findHost {
				s.handleDomain(defaultDoamin, c, n)
			}
		})
		srv := r.ServerKeepAlive(sv.Listen)
		//构造

		s.Logger.Printf(infoServerPrefix+"start to accept http conn: %v\n", srv.Addr)
		s.httpServers = append(s.httpServers, srv)
		if sv.SSL {
			cfg := &tls.Config{}
			for _, domain := range sv.Domain {
				if domain.CertFile != "" && domain.KeyFile != "" {
					cert, err := tls.LoadX509KeyPair(domain.CertFile, domain.KeyFile)
					if err != nil {
						log.Fatal(err)
					}
					cfg.Certificates = append(cfg.Certificates, cert)
				}
			}
			cfg.BuildNameToCertificate()
			srv.TLSConfig = cfg
			go func() {
				err := srv.ListenAndServeTLS("", "")
				if err != http.ErrServerClosed {
					s.Logger.Printf(errorServerPrefix+"start http server error: %s\n", err)
				}
				s.Logger.Printf(infoServerPrefix+"stop http server %v\n", srv.Addr)
			}()
		} else {
			go func() {
				err := srv.ListenAndServe()
				if err != http.ErrServerClosed {
					s.Logger.Printf(errorServerPrefix+"start http server error: %s\n", err)
				}
				s.Logger.Printf(infoServerPrefix+"stop http server %v\n", srv.Addr)
			}()
		}
	}
	return nil
}

func (s *Server) Reload() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}
	//stop all httpServer
	for _, srv := range s.httpServers {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			s.Logger.Printf(errorServerPrefix+"stop http server err %v\n", err)
		}
	}
	//reload config
	sections.Load()
	s.httpServers = []*http.Server{}
	clearProxys()
	go s.setupCart()
	return nil
}

func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true
	s.ShutdownCh <- struct{}{}
	return nil
}
