package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	go server.setupCart()

	return server, nil
}

func (s *Server) excuteDoamin(domain *sections.Domain, c *cart.Context, n cart.Next) {
	if domain != nil {
		for _, location := range domain.Location {
			match := false
			regex := false
			path := location.Path
			//match regex path
			if location.Path[0] == '~' {
				path = strings.TrimSpace(location.Path[1:])
				regex = true
			}
			if regex {
				match, _ = regexp.MatchString(path, c.Request.RequestURI)
			} else {
				if strings.HasPrefix(c.Request.RequestURI, path) {
					match = true
				}
				if path == "/" {
					match = true
				}
			}
			if match {
				if len(location.Header) != 0 {
					for _, header := range location.Header {
						arr := strings.SplitN(header, " ", 2)
						key := arr[0]
						value := arr[1]
						c.Response.Header().Set(key, value)
					}
				}
				if location.Return != "" {
					arr := strings.SplitN(location.Return, " ", 3)
					code, _ := strconv.Atoi(arr[0])
					_type := arr[1]
					content := ""
					if len(arr) > 2 {
						content = arr[2]
					}
					switch code {
					case 301, 302: //move
						path := _type
						path = strings.ReplaceAll(path, "$host", c.Request.Host)
						path = strings.ReplaceAll(path, "$request_uri", c.Request.RequestURI)
						c.Redirect(code, path)
					default:
						c.Response.WriteHeader(code)
						switch _type {
						case "json":
							header := c.Response.Header()
							header["Content-Type"] = []string{"application/json; charset=utf-8"}
							c.Response.Write([]byte(content))
						case "html":
							header := c.Response.Header()
							header["Content-Type"] = []string{"text/html; charset=utf-8"}
							c.Response.Write([]byte(content))
						case "string":
							header := c.Response.Header()
							header["Content-Type"] = []string{"text/plain; charset=utf-8"}
							c.Response.Write([]byte(content))
						case "file":
							c.File(content)
						case "static":
							c.Static(content, path, true)
						default:

						}

					}
				}
				break
			}
		}
	}
	n()
}

func (s *Server) setupCart() error {
	for _, sv := range *sections.Server {
		cart.SetMode(cart.ReleaseMode)
		r := cart.New()
		tempDomain := sv.Domain
		//构造
		r.Use("/", func(c *cart.Context, n cart.Next) {
			//timer
			since := time.Now()
			defer func() {
				d := time.Since(since)
				log.Println("[Track]", c.Request.RequestURI, d)
			}()
			findHost := false
			var defaultDoamin *sections.Domain
			for _, domain := range tempDomain {
				if strings.EqualFold(domain.Name, c.Request.Host) {
					findHost = true
					s.excuteDoamin(domain, c, n)
				}
				if domain.Name == "_" {
					defaultDoamin = domain
				}
			}
			if !findHost {
				s.excuteDoamin(defaultDoamin, c, n)
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

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	// director :=

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "pipe")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("502 bad gateway"))
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		//return errors.New("response body is invalid")
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

// func Run() {

// 	shutdownCh := make(chan struct{})

// 	for k, v := range sections.Endpoint {
// 		fmt.Println(k, v)
// 	}

// 	for k, v := range sections.Servers {
// 		fmt.Println(k, v)
// 	}
// 	// initialize a reverse proxy and pass the actual backend server url here
// 	// proxy, err := NewProxy("http://127.0.0.1:8091")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// c := cart.Default()

// 	// c.Route("/", func(r *cart.Router) {
// 	// 	c.Route("/").ANY(func(c *cart.Context, n cart.Next) {
// 	// 		fmt.Println("/")
// 	// 		proxy.ServeHTTP(c.Response, c.Request)
// 	// 	})
// 	// 	c.Route("/api").ANY(func(c *cart.Context, n cart.Next) {
// 	// 		fmt.Println("/api")
// 	// 		proxy.ServeHTTP(c.Response, c.Request)
// 	// 		n()
// 	// 	})
// 	// 	c.Route("/api/*path").ANY(func(c *cart.Context, n cart.Next) {
// 	// 		fmt.Println("/api/test")
// 	// 		proxy.ServeHTTP(c.Response, c.Request)
// 	// 		n()
// 	// 	})
// 	// })

// 	// _, _ = c.Run(":80")

// 	sigs := make(chan os.Signal, 10)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
// 	go func() {
// 		for {
// 			sig := <-sigs
// 			fmt.Println()
// 			log.Printf("get signal %v\n", sig)
// 			if sig == syscall.SIGUSR2 {
// 				//grace reload
// 				// s.Self.LoadServices()
// 				// s.Shutter()
// 				close(shutdownCh)
// 			} else {
// 				close(shutdownCh)
// 				//s.Shutdown()
// 			}
// 		}
// 	}()
// 	<-shutdownCh
// }
