package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
)

type Server struct {
	server http.Server
	cfg    *config.Config
}

// Create a new http-proxy server ...
func NewServer() (*Server, error) {
	s := &Server{
		cfg: config.Get(),
	}
	if s.cfg.Proxy.Port == 0 || s.cfg.Proxy.Address == "" {
		return nil, errors.New("no config file specified")
	}
	s.server.Addr += s.cfg.Proxy.Address + ":" + strconv.Itoa(s.cfg.Proxy.Port)
	s.server.Handler = s
	return s, nil
}

// Start serving ...
func (s *Server) Serve() {
	print.Info("wingman dev proxy listening on " + strconv.Itoa(s.cfg.Proxy.Port))
	if err := s.server.ListenAndServe(); err != nil {
		print.SvcErr("proxy", "Proxy crashed! "+err.Error())
	}
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check whether to log requests in the terminal
	if s.cfg.Proxy.LogRequests {
		print.SvcProxy(r.Method + "@" + r.URL.String())
	}

	for _, svc := range s.cfg.Services {
		if svc.ProxyHandle != "" && strings.HasPrefix(r.URL.String(), svc.ProxyHandle) {
			switch svc.ProxyType {
			case "service":
				s.getService(w, r, svc)
			case "static":
				s.getFile(w, r.URL, svc)
			default:
			}
		}
	}
}

func (s *Server) getFile(w http.ResponseWriter, url *url.URL, service config.ServiceConfig) {
	path := strings.Replace(url.String(), service.ProxyHandle, service.ProxyStaticDir, 1)
	if _, err := os.Stat(path); err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "404 Not Found")
		return
	}
	switch filepath.Ext(path) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpeg", ".jpg":
		w.Header().Set("Content-Type", "image/jpg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg", ".svgz":
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	w.WriteHeader(200)
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(w, "could not read file: %s", err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	fmt.Fprintf(w, "%s", string(file))
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request, service config.ServiceConfig) {
	if len(r.URL.String()) > 1 {
		if service.ProxyRouteRewrite != "" && strings.Contains(service.ProxyRouteRewrite, ":") {
			parts := strings.SplitN(service.ProxyRouteRewrite, ":", 2)
			r.URL.Path = strings.Replace(
				r.URL.Path,
				"/"+parts[0]+"/",
				"/"+parts[1]+"/",
				1,
			)
		}
		r.Header.Set("Origin", "http://"+r.Host)
		proxyRoute(service.ProxyAddress, service.ProxyPort, w, r)
	}
}

// proxyRoute ...
func proxyRoute(address string, port int, w http.ResponseWriter, r *http.Request) {
	if address == "" || port == 0 {
		print.SvcErr("proxy", "could not create proxy URL, address or port mising")
		return
	}
	u, err := url.Parse("http://" + address + ":" + strconv.Itoa(port))
	if err != nil {
		print.SvcErr("proxy", "could not create proxy URL: "+err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.FlushInterval = time.Duration(1 * time.Second)
	proxy.ServeHTTP(w, r)
}
