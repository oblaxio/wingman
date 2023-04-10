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

// NewServer ...
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

// Serve ...
func (s *Server) Serve() {
	print.Info("wingman dev proxy listening on " + strconv.Itoa(s.cfg.Proxy.Port))
	if err := s.server.ListenAndServe(); err != nil {
		print.SvcErr("proxy", "Proxy crashed!")
	}
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check whether to log requests in the terminal
	if s.cfg.Proxy.LogRequests {
		print.SvcProxy(r.Method + "@" + r.URL.String())
	}
	if strings.HasPrefix(r.URL.String(), "/"+s.cfg.Proxy.APIPrefix) && s.cfg.Proxy.APIPrefix != "" {
		// if it's an API request
		s.getAPI(w, r)
	} else if strings.HasPrefix(r.URL.String(), "/"+s.cfg.Proxy.Storage.Prefix) && s.cfg.Proxy.Storage.Prefix != "" {
		// if it's a storage request
		s.getStorageItem(w, r)
	} else {
		if s.cfg.Proxy.SPA.Port != 0 {
			// if it's a SPA request
			proxyRoute(s.cfg.Proxy.SPA.Address, s.cfg.Proxy.SPA.Port, w, r)
		} else {
			// if it's a static file request
			s.getFile(w, r.URL)
		}
	}
}

// getFile ...
func (s *Server) getFile(w http.ResponseWriter, url *url.URL) {
	path := s.cfg.Proxy.Static.Dir + url.String()
	if strings.HasSuffix(url.String(), "/") {
		path += s.cfg.Proxy.Static.Index
	}
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
	case ".jpeg":
	case ".jpg":
		w.Header().Set("Content-Type", "image/jpg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
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

// getAPI ...
func (s *Server) getAPI(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.String()) > 1 {
		var svc *config.ServiceConfig
		for _, v := range s.cfg.Services {
			if v.ProxyHandle == r.URL.String() {
				svc = &v
				break
			}
		}
		if svc == nil {
			print.SvcErr("proxy", "could not create proxy URL")
			return
		}
		r.URL.Path = r.URL.String()
		r.Header.Set("Origin", "http://"+r.Host)
		proxyRoute(svc.ProxyAddress, svc.ProxyPort, w, r)
	}
}

// getStorageItem ...
func (s *Server) getStorageItem(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(
		r.URL.Path,
		"/"+s.cfg.Proxy.Storage.Prefix+"/",
		"/"+s.cfg.Proxy.Storage.Bucket+"/",
		1,
	)
	proxyRoute(s.cfg.Proxy.Storage.Address, s.cfg.Proxy.Storage.Port, w, r)
}

// proxyRoute ...
func proxyRoute(address string, port int, w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse("http://" + address + ":" + strconv.Itoa(port))
	if err != nil {
		print.SvcErr("proxy", "could not create proxy URL: "+err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.FlushInterval = time.Duration(1 * time.Second)
	proxy.ServeHTTP(w, r)
}
