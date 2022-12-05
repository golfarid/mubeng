package server

import (
	"github.com/gorilla/mux"
	"ktbs.dev/mubeng/internal/api/controllers/auth"
	"ktbs.dev/mubeng/internal/api/controllers/proxies"
	jwtAuthMiddleware "ktbs.dev/mubeng/internal/api/middlewares/auth"
	"net/http"
	"os"
	"os/signal"

	"github.com/elazarl/goproxy"
	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
	"ktbs.dev/mubeng/common"
)

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	cli := logo.NewReceiver(os.Stderr, "")
	cli.Color = true
	cli.Level = logo.DEBUG

	file, _ := logo.Open(opt.Output)
	out := logo.NewReceiver(file, "")
	out.Format = "%s: %s"

	dump = &httpretty.Logger{
		RequestHeader:  true,
		ResponseHeader: true,
		Colors:         true,
	}

	handler = &Proxy{}
	handler.Options = opt
	handler.HTTPProxy = goproxy.NewProxyHttpServer()
	handler.HTTPProxy.OnRequest().DoFunc(handler.onRequest)
	handler.HTTPProxy.OnRequest().HandleConnectFunc(handler.onConnect)
	handler.HTTPProxy.OnResponse().DoFunc(handler.onResponse)

	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()

	if opt.Auth != "" {
		authMiddleware := jwtAuthMiddleware.New(opt)
		apiRouter.Use(authMiddleware.Handle)
		authController := auth.New(opt)
		router.HandleFunc("/auth/sign_in", authController.Handler).Methods("POST")
	}

	proxiesController := proxies.New(opt, opt.ProxyManager)
	apiRouter.HandleFunc("/proxies", proxiesController.Handler).Methods("GET", "POST", "DELETE")

	router.HandleFunc("/", nonProxy)
	handler.HTTPProxy.NonproxyHandler = router

	server = &http.Server{
		Addr:    opt.Address,
		Handler: handler.HTTPProxy,
	}

	log = logo.NewLogger(cli, out)

	if opt.Watch {
		watcher, err := opt.ProxyManager.Watch()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		go watch(watcher)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go interrupt(stop)

	log.Infof("[PID: %d] Starting proxy server on %s", os.Getpid(), opt.Address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
