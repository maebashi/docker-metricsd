package main

import (
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/maebashi/docker-metricsd/utils"
)

func main() {
	var address string
	flag.StringVar(&address, "A", ":12375", "bind to address:port")
	var sockpath = flag.String("sock-path", "/var/run/docker.sock",
		"docker's socket path")
	flag.Parse()

	u, _ := url.Parse("http://127.0.0.1:2375/")
	reverse_proxy := httputil.NewSingleHostReverseProxy(u)
	reverse_proxy.Transport = &http.Transport{
		Proxy: nil,
		Dial: func(n, a string) (net.Conn, error) {
			return net.Dial("unix", *sockpath)
		},
	}

	http.HandleFunc("/", handler(reverse_proxy))
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		a := strings.Split(r.URL.Path, "/")
		if a[1] == "containers" && len(a) > 3 {
			switch {
			case a[3] == "json":
				err = handleContainersByName(p, w, r)
			default:
				p.ServeHTTP(w, r)
			}
		} else {
			p.ServeHTTP(w, r)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(([]byte)(err.Error()))
		}
	}
}

func handleContainersByName(p *httputil.ReverseProxy, w http.ResponseWriter, r *http.Request) (err error) {
	rr := httptest.NewRecorder()
	p.ServeHTTP(rr, r)

	a := strings.Split(r.URL.Path, "/")
	pid, err := utils.GetContainerPID(a[2])
	if err != nil {
		return
	}

	var data map[string]interface{}
	dec := json.NewDecoder(rr.Body)
	if err = dec.Decode(&data); err != nil {
		return
	}

	metrics := map[string]interface{}{}

	var m map[string]float64
	for _, subsystem := range []string{"memory", "cpuacct"} {
		if m, err = utils.GetCgroupStats(a[2], subsystem); err != nil {
			continue
		}
		metrics[subsystem] = m
	}

	err = utils.NetNsSynchronize(pid, func() (err error) {
		ifstats, err := utils.GetIfStats()
		if err != nil {
			return
		}
		metrics["interfaces"] = ifstats

		if addr := utils.GetIfAddr("eth0"); addr != nil {
			if nets, ok := data["NetworkSettings"]; ok {
				a := strings.Split(addr.String(), "/")
				n := nets.(map[string]interface{})
				n["IPAddress"] = a[0]
				n["IPPrefixLen"] = a[1]
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	data["Metrics"] = metrics
	s, err := json.Marshal(data)
	if err != nil {
		return
	}
	rr.Body.Reset()
	rr.Body.Write(s)
	copyResponse(w, rr)
	return
}

func copyResponse(dst http.ResponseWriter, src *httptest.ResponseRecorder) {
	blen := src.Body.Len()
	src.Header().Set("Content-Length", strconv.Itoa(blen))
	for k, vv := range src.HeaderMap {
		for _, v := range vv {
			dst.Header().Add(k, v)
		}
	}
	dst.WriteHeader(src.Code)
	io.Copy(dst, src.Body)
}
