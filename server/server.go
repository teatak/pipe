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

func (s *Server) setupCart() error {
	for _, sv := range *sections.Server {
		cart.SetMode(cart.ReleaseMode)
		r := cart.New()

		//构造
		r.Use("/", func(c *cart.Context, n cart.Next) {
			switch c.Request.Host {

			}
			// sort.Slice(sv.Domain, func(i, j int) bool {
			// 	return sv.Domain[j].Name == "_"
			// })
			//match host
			//sort domain name
			findHost := false
			for _, domain := range sv.Domain {
				for _, _ = range domain.Location {
					//fmt.Println(path, localtion)
				}
				if strings.EqualFold(domain.Name, c.Request.Host) {
					findHost = true
					if c.Request.TLS == nil {
						c.Redirect(301, "https://"+c.Request.Host+c.Request.RequestURI)
					} else {
						c.JSON(200, cart.H{"code": 200})
					}
				}
			}
			if !findHost {
				if c.Request.TLS == nil {
					c.Redirect(301, "https://"+c.Request.Host+c.Request.RequestURI)
				} else {
					c.JSON(200, cart.H{"code": 200})
				}
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
