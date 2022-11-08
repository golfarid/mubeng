package api

import (
	"github.com/gorilla/mux"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/api/controllers/auth"
	"ktbs.dev/mubeng/internal/api/controllers/proxies"
	jwtAuthMiddleware "ktbs.dev/mubeng/internal/api/middlewares/auth"
	"ktbs.dev/mubeng/internal/api/utils"
	"net/http"
	"os"
	"os/signal"
)

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	log = utils.Logger(opt.Output)

	proxyManager = opt.ProxyManager
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()

	authMiddleware := jwtAuthMiddleware.New(opt)
	apiRouter.Use(authMiddleware.Handle)

	authController := auth.New(opt)
	proxiesController := proxies.New(opt, proxyManager)

	// replace http.HandleFunc with myRouter.HandleFunc
	router.HandleFunc("/auth/sign_in", authController.Handler).Methods("POST")
	apiRouter.HandleFunc("/proxies", proxiesController.Handler).Methods("GET", "POST")
	apiRouter.HandleFunc("/proxies/{index}", proxiesController.Handler).Methods("DELETE")

	api = &http.Server{
		Addr:    opt.ApiAddress,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go interrupt(stop)

	log.Infof("[PID: %d] Starting API on %s", os.Getpid(), opt.ApiAddress)
	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
