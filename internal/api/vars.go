package api

import (
	"github.com/mbndr/logo"
	"ktbs.dev/mubeng/internal/proxymanager"
	"net/http"
)

var (
	api          *http.Server
	proxyManager *proxymanager.ProxyManager
	mime         = "application/json"
	log          *logo.Logger
)
