// forwarding proxy
// https://github.com/betalo-sweden/forwardingproxy

// Copyright (C) 2018 Betalo AB - All Rights Reserved

// Courtesy: https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
// $ openssl req -newkey rsa:2048 -nodes -keyout server.key -new -x509 -sha256 -days 3650 -out server.pem

package internal

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	//"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"

	//"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
)

// Proxy is a HTTPS forward proxy.
type Proxy struct {
	//Logger              *zap.Logger
	AuthUser            string
	AuthPass            string
	ForwardingHTTPProxy *httputil.ReverseProxy
	DestDialTimeout     time.Duration
	DestReadTimeout     time.Duration
	DestWriteTimeout    time.Duration
	ClientReadTimeout   time.Duration
	ClientWriteTimeout  time.Duration
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//p.Logger.Info("Incoming request", zap.String("host", r.Host))

	if p.AuthUser != "" && p.AuthPass != "" {
		user, pass, ok := parseBasicProxyAuth(r.Header.Get("Proxy-Authorization"))
		if !ok || user != p.AuthUser || pass != p.AuthPass {
			//p.Logger.Warn("Authorization attempt with invalid credentials")
			http.Error(w, http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired)
			return
		}
	}

	if r.URL.Scheme == "http" {
		p.handleHTTP(w, r)
	} else {
		p.handleTunneling(w, r)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	//p.Logger.Debug("Got HTTP request", zap.String("host", r.Host))
	p.ForwardingHTTPProxy.ServeHTTP(w, r)
}

func (p *Proxy) handleTunneling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		//p.Logger.Info("Method not allowed", zap.String("method", r.Method))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	//p.Logger.Debug("Connecting", zap.String("host", r.Host))

	destConn, err := net.DialTimeout("tcp", r.Host, p.DestDialTimeout)
	if err != nil {
		//p.Logger.Error("Destination dial failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//p.Logger.Debug("Connected", zap.String("host", r.Host))

	w.WriteHeader(http.StatusOK)

	//p.Logger.Debug("Hijacking", zap.String("host", r.Host))

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		//p.Logger.Error("Hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		//p.Logger.Error("Hijacking failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//p.Logger.Debug("Hijacked connection", zap.String("host", r.Host))

	now := time.Now()
	clientConn.SetReadDeadline(now.Add(p.ClientReadTimeout))
	clientConn.SetWriteDeadline(now.Add(p.ClientWriteTimeout))
	destConn.SetReadDeadline(now.Add(p.DestReadTimeout))
	destConn.SetWriteDeadline(now.Add(p.DestWriteTimeout))

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
	defer func() { _ = dest.Close() }()
	defer func() { _ = src.Close() }()
	_, _ = io.Copy(dest, src)
}

// parseBasicProxyAuth parses an HTTP Basic Authorization string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicProxyAuth(authz string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(authz, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(authz[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// NewForwardingHTTPProxy returns a new reverse proxy that takes an incoming
// request and sends it to another server, proxying the response back to the
// client.
//
// See: https://golang.org/pkg/net/http/httputil/#ReverseProxy
func NewForwardingHTTPProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	// TODO:(alesr) Use timeouts specified via flags to customize the default
	// transport used by the reverse proxy.
	return &httputil.ReverseProxy{
		//ErrorLog: logger,
		Director: director,
	}
}

// ForwardProxy is a forwarding HTTP/S proxy
func ForwardProxy(port int) {
	// var (
	// 	flagCertPath = flag.String("cert", "", "Filepath to certificate")
	// 	flagKeyPath  = flag.String("key", "", "Filepath to private key")
	// 	flagAddr     = flag.String("addr", "", "Server address")
	// 	flagAuthUser = flag.String("user", "", "Server authentication username")
	// 	flagAuthPass = flag.String("pass", "", "Server authentication password")
	// 	flagVerbose  = flag.Bool("verbose", false, "Set log level to DEBUG")

	// 	flagDestDialTimeout         = flag.Duration("destdialtimeout", 10*time.Second, "Destination dial timeout")
	// 	flagDestReadTimeout         = flag.Duration("destreadtimeout", 5*time.Second, "Destination read timeout")
	// 	flagDestWriteTimeout        = flag.Duration("destwritetimeout", 5*time.Second, "Destination write timeout")
	// 	flagClientReadTimeout       = flag.Duration("clientreadtimeout", 5*time.Second, "Client read timeout")
	// 	flagClientWriteTimeout      = flag.Duration("clientwritetimeout", 5*time.Second, "Client write timeout")
	// 	flagServerReadTimeout       = flag.Duration("serverreadtimeout", 30*time.Second, "Server read timeout")
	// 	flagServerReadHeaderTimeout = flag.Duration("serverreadheadertimeout", 30*time.Second, "Server read header timeout")
	// 	flagServerWriteTimeout      = flag.Duration("serverwritetimeout", 30*time.Second, "Server write timeout")
	// 	flagServerIdleTimeout       = flag.Duration("serveridletimeout", 30*time.Second, "Server idle timeout")

	// 	flagLetsEncrypt = flag.Bool("letsencrypt", false, "Use letsencrypt for https")
	// 	flagLEWhitelist = flag.String("lewhitelist", "", "Hostname to whitelist for letsencrypt")
	// 	flagLECacheDir  = flag.String("lecachedir", "/tmp", "Cache directory for certificates")
	// )

	// flag.Parse()

	var (
		flagCertPath = ""
		flagKeyPath  = ""
		flagAddr     = fmt.Sprintf("127.0.0.1:%v", port)
		flagAuthUser = ""
		flagAuthPass = ""
		//flagVerbose  = true

		flagDestDialTimeout         = 10 * time.Second
		flagDestReadTimeout         = 5 * time.Second
		flagDestWriteTimeout        = 5 * time.Second
		flagClientReadTimeout       = 5 * time.Second
		flagClientWriteTimeout      = 5 * time.Second
		flagServerReadTimeout       = 30 * time.Second
		flagServerReadHeaderTimeout = 30 * time.Second
		flagServerWriteTimeout      = 30 * time.Second
		flagServerIdleTimeout       = 30 * time.Second

		flagLetsEncrypt = false
		flagLEWhitelist = ""
		flagLECacheDir  = "/tmp"
	)

	//
	// c := zap.NewProductionConfig()
	// c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// if flagVerbose {
	// 	c.Level.SetLevel(zapcore.DebugLevel)
	// } else {
	// 	c.Level.SetLevel(zapcore.ErrorLevel)
	// }

	// logger, err := c.Build()
	// if err != nil {
	// 	log.Fatalln("Error: failed to initiate logger")
	// }
	// defer logger.Sync()
	//stdLogger := zap.NewStdLog(logger)

	p := &Proxy{
		ForwardingHTTPProxy: NewForwardingHTTPProxy(),
		//Logger:              logger,
		AuthUser:           flagAuthUser,
		AuthPass:           flagAuthPass,
		DestDialTimeout:    flagDestDialTimeout,
		DestReadTimeout:    flagDestReadTimeout,
		DestWriteTimeout:   flagDestWriteTimeout,
		ClientReadTimeout:  flagClientReadTimeout,
		ClientWriteTimeout: flagClientWriteTimeout,
	}

	s := &http.Server{
		Addr:    flagAddr,
		Handler: p,
		//ErrorLog:          stdLogger,
		ReadTimeout:       flagServerReadTimeout,
		ReadHeaderTimeout: flagServerReadHeaderTimeout,
		WriteTimeout:      flagServerWriteTimeout,
		IdleTimeout:       flagServerIdleTimeout,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	if flagLetsEncrypt {
		if flagLEWhitelist == "" {
			//p.Logger.Fatal("error: no -lewhitelist flag set")
		}
		if flagLECacheDir == "/tmp" {
			//p.Logger.Info("-lecachedir should be set, using '/tmp' for now...")
		}

		m := &autocert.Manager{
			Cache:      autocert.DirCache(flagLECacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(flagLEWhitelist),
		}

		s.Addr = ":https"
		s.TLSConfig = m.TLSConfig()
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		//p.Logger.Info("Server shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown failed: %v\n", err)

			//p.Logger.Error("Server shutdown failed", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	//p.Logger.Info("Server starting", zap.String("address", s.Addr))

	var svrErr error
	if flagCertPath != "" && flagKeyPath != "" || flagLetsEncrypt {
		svrErr = s.ListenAndServeTLS(flagCertPath, flagKeyPath)
	} else {
		svrErr = s.ListenAndServe()
	}

	if svrErr != http.ErrServerClosed {
		//p.Logger.Error("Listening for incoming connections failed", zap.Error(svrErr))
	}

	<-idleConnsClosed
	//p.Logger.Info("Server stopped")
	log.Println("Server stopped")
}
