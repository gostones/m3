package internal

import (
	"fmt"
	"github.com/dhnt/m3/internal/rp"
	"github.com/dhnt/m3/internal/tunnel"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	rpsIni = `
[common]
bind_addr = 0.0.0.0
bind_port = 7000

#
vhost_http_port = 1080
vhost_https_port = 1443

#
dashboard_addr = 0.0.0.0
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd = password

# trace, debug, info, warn, error
log_level = info
log_max_days = 1

#
#subdomain_host =
tcp_mux = true
#
`
	//server_port, instance, service_host, service_port, remote_port
	rpcIni = `
[common]
server_addr = localhost
server_port = %v
http_proxy =

[ssh_random]
type = tcp
local_ip = %v
local_port = %v
remote_port = %v
`
	//
	listenPort = 8080

	rpsPort = 7000
)

func rpServer(listen int, web string) {
	// flags := flag.NewFlagSet("server", flag.ContinueOnError)

	// //tunnel port
	// listen := flags.Int("port", parseInt(os.Getenv("PORT"), listenPort), "server listening port")

	// web := flags.String("web", os.Getenv("FG_WEB"), "web url")

	// flags.Parse(args)

	if web == "" {
		//default to dashboard
		web = "http://localhost:7500"
	}

	//
	go rp.Server(rpsIni)

	port := FreePort()
	proxy := fmt.Sprintf("http://localhost:%v", port)
	go serve(port, web)

	tunnel.TunServer(listen, proxy)
}

//
func rpClient(url, proxy, hostport string, port int) {
	// flags := flag.NewFlagSet("client", flag.ContinueOnError)

	// //
	// url := flags.String("url", os.Getenv("FG_URL"), "tunnel url")
	// proxy := flags.String("proxy", "", "http proxy")
	// hostPort := flags.String("hostport", "", "reverse proxy service host:port")
	// port := flags.Int("port", -1, "remote reverse proxy port")

	lport := FreePort()

	remote := fmt.Sprintf("localhost:%v:localhost:%v", lport, rpsPort)

	fmt.Fprintf(os.Stdout, "remote: %v\n", remote)

	go tunnel.TunClient(proxy, url, remote)

	sleep := BackoffDuration()

	for {
		hp := strings.Split(hostport, ":")
		shost := hp[0]

		sport, err := strconv.Atoi(hp[1])
		if err != nil {
			panic(err)
		}

		rp.Client(fmt.Sprintf(rpcIni, lport, shost, sport, port))

		//should never return or error
		sleep(fmt.Errorf("Reverse proxy error"))
	}
}

func rpConnect(url, proxy, ports string) {
	// flags := flag.NewFlagSet("connect", flag.ContinueOnError)

	// //
	// url := flags.String("url", os.Getenv("FG_URL"), "tunnel url")
	// proxy := flags.String("proxy", "", "http proxy")
	// ports := flags.String("ports", "", "local_port:remote_port")

	pa := strings.Split(ports, ":")

	remote := fmt.Sprintf("localhost:%v:localhost:%v", pa[0], pa[1])

	fmt.Fprintf(os.Stdout, "remote: %v\n", remote)

	tunnel.TunClient(proxy, url, remote)
}

func parseInt(s string, v int) int {
	if s == "" {
		return v
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		i = v
	}
	return i
}

func serve(port int, target string) {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	scheme := "http"
	vhost := "localhost:8080"

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		req.URL.Host = vhost
		req.URL.Scheme = scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = vhost

		proxy.ServeHTTP(res, req)
	})

	log.Printf("serve port: %v target: %v\n", port, target)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Println(err)
	}
}

//package internal

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// )

// var (
// 	hostProxy map[string]*httputil.ReverseProxy
// )

// type baseHandle struct{}

// func (h baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	host := r.Host

// 	if fn, ok := hostProxy[host]; ok {
// 		fn.ServeHTTP(w, r)
// 		return
// 	}

// 	target := "http://" + host

// 	log.Println("target: ", target)

// 	remote, err := url.Parse(target)
// 	if err != nil {
// 		log.Println("target parse fail:", err)
// 		return
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(remote)
// 	hostProxy[host] = proxy
// 	proxy.ServeHTTP(w, r)
// 	return

// }

// func reverseProxy(port int) {
// 	hostProxy = make(map[string]*httputil.ReverseProxy)
// 	h := &baseHandle{}
// 	http.Handle("/", h)

// 	server := &http.Server{
// 		Addr:    fmt.Sprintf(":%v", port),
// 		Handler: h,
// 	}

// 	//server.ListenAndServeTLS( "server.pem", "server.key")

// 	log.Fatal(server.ListenAndServe())
// }
