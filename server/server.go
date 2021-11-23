package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/teatak/pipe/sections"
)

var server *Server

type Server struct {
	Logger       *log.Logger
	ShutdownCh   chan struct{}
	shutdown     bool
	shutdownLock sync.Mutex
}

func NewServer() (*Server, error) {

	shutdownCh := make(chan struct{})

	server = &Server{
		ShutdownCh: shutdownCh,
	}

	logOutput := io.MultiWriter(os.Stderr)
	server.Logger = log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds)

	return server, nil
}

func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true
	close(s.ShutdownCh)
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

func Run() {

	shutdownCh := make(chan struct{})

	for k, v := range sections.Endpoint {
		fmt.Println(k, v)
	}

	for k, v := range sections.Servers {
		fmt.Println(k, v)
	}
	// initialize a reverse proxy and pass the actual backend server url here
	// proxy, err := NewProxy("http://127.0.0.1:8091")
	// if err != nil {
	// 	panic(err)
	// }

	// c := cart.Default()

	// c.Route("/", func(r *cart.Router) {
	// 	c.Route("/").ANY(func(c *cart.Context, n cart.Next) {
	// 		fmt.Println("/")
	// 		proxy.ServeHTTP(c.Response, c.Request)
	// 	})
	// 	c.Route("/api").ANY(func(c *cart.Context, n cart.Next) {
	// 		fmt.Println("/api")
	// 		proxy.ServeHTTP(c.Response, c.Request)
	// 		n()
	// 	})
	// 	c.Route("/api/*path").ANY(func(c *cart.Context, n cart.Next) {
	// 		fmt.Println("/api/test")
	// 		proxy.ServeHTTP(c.Response, c.Request)
	// 		n()
	// 	})
	// })

	// _, _ = c.Run(":80")

	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			fmt.Println()
			log.Printf("get signal %v\n", sig)
			if sig == syscall.SIGUSR2 {
				//grace reload
				// s.Self.LoadServices()
				// s.Shutter()
				close(shutdownCh)
			} else {
				close(shutdownCh)
				//s.Shutdown()
			}
		}
	}()
	<-shutdownCh
}
