package srv

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vislee/goframe/util"
	"io"
	"net"
	"net/http"
	httpprof "net/http/pprof"
	"net/url"
	"runtime"
	"strconv"
)

type httpSrv struct {
	baseSrv
	listener net.Listener
}

func NewHttpSrv(cfg *Config, log *util.BLogs) *httpSrv {
	return &httpSrv{baseSrv: baseSrv{cfgInfo: cfg, logs: log}}
}

func (hs *httpSrv) InitListener() error {
	listenHost, ok := hs.cfgInfo.HttpListen[hs.serviceCode]
	if !ok {
		var errMsg string = "InitListener addr nil."
		hs.logs.Warn(errMsg)
		return errors.New(errMsg)
	}

	var err error
	hs.listener, err = net.Listen("tcp", listenHost.Host)
	if err != nil {
		errmsg := fmt.Sprintf("listen host:%s error. error: %s", listenHost.Host, err.Error())
		hs.logs.Error(errmsg)
		return errors.New(errmsg)
	}
	hs.logs.Debug("listen host: %s ok", listenHost.Host)
	return nil
}

func (hs *httpSrv) Close() {
	if hs.listener != nil {
		hs.listener.Close()
	}
}

func (hs *httpSrv) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	hs.logs.Debug("remote ip: %s, url: %s", req.RemoteAddr, req.URL.Path)
	hs.routerV01(w, req)
}

func (hs *httpSrv) httpMain() {
	server := &http.Server{
		Handler: hs,
	}
	err := server.Serve(hs.listener)
	if err != nil {
		hs.logs.Error("http serve run error. error: %s", err.Error())
	}
}

func (hs *httpSrv) routerV01(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/ping":
		hs.pingHandler(w, req)
	case "/debug/pprof":
		setPub(w)
		httpprof.Index(w, req)
	case "/debug/pprof/block":
		hs.debugRuntime(w, req, "block")
	case "/debug/pprof/goroutine":
		hs.debugRuntime(w, req, "goroutine")
	case "/debug/pprof/heap":
		hs.debugRuntime(w, req, "heap")
	case "/debug/pprof/threadcreate":
		hs.debugRuntime(w, req, "threadcreate")
	case "/debug/pprof/profile":
		setPub(w)
		httpprof.Profile(w, req)
	case "/debug/pprof/cmdline":
		setPub(w)
		httpprof.Cmdline(w, req)
	case "/debug/pprof/symbol":
		setPub(w)
		httpprof.Symbol(w, req)
	default:
		RespondV1(w, 404, "page not found")
	}
}

func (hs *httpSrv) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Length", "4")
	io.WriteString(w, "PONG")
}

func (hs *httpSrv) debugRuntime(w http.ResponseWriter, req *http.Request, s string) {
	setPub(w)
	if s == "heap" {
		runtime.GC()
	}
	httpprof.Handler(s).ServeHTTP(w, req)
	return
}

func (hs *httpSrv) acceptVersion(req *http.Request) int {
	if req.Header.Get("accept") == "application/GoFrame; version=1.0" {
		return 1
	}

	return 0
}

func RespondV1(w http.ResponseWriter, code int, data interface{}) {
	var response []byte
	var err error
	var isJSON bool

	if code == 200 {
		switch data.(type) {
		case string:
			response = []byte(data.(string))
		case []byte:
			response = data.([]byte)
		case nil:
			response = []byte{}
		default:
			isJSON = true
			response, err = json.Marshal(data)
			if err != nil {
				code = 500
				data = err
			}
		}
	}

	if code != 200 {
		isJSON = true
		response = []byte(fmt.Sprintf(`{"message":"%s"}`, data))
	}

	if isJSON {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.Header().Set("Server", "GoFrame 1.0 @liwq")
	w.Header().Set("X-GoFrame-Content-Type", "GoFrame; version=1.0")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.WriteHeader(code)
	w.Write(response)
}

func setPub(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Server", "GoFrame 1.0 @liwq")
	w.Header().Set("X-GoFrame-Content-Type", "GoFrame; version=1.0")
}
