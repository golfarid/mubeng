package proxies

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/api/utils"
	"ktbs.dev/mubeng/internal/proxymanager"
	"net/http"
	"strconv"
)

type Controller struct {
	proxyManager *proxymanager.ProxyManager
}

func New(opt *common.Options, proxyManager *proxymanager.ProxyManager) *Controller {
	log = utils.Logger(opt.Output)
	return &Controller{proxyManager: proxyManager}
}

type ProxyDto struct {
	Index int    `json:"index"`
	Url   string `json:"url"`
}

func (controller *Controller) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		controller.getProxies(w)
	case "POST":
		controller.createProxies(w, r)
	case "DELETE":
		controller.deleteProxy(w, r)
	}
}

func (controller *Controller) getProxies(w http.ResponseWriter) {
	fmt.Printf("got /get proxies request\n")
	proxies, err := controller.proxyManager.ReadProxies()
	if err != nil {
		log.Error(err)
		response := "Can't read proxies"
		http.Error(w, response, http.StatusInternalServerError)
	} else {
		var proxiesDto []ProxyDto
		for i := 0; i < len(proxies); i++ {
			proxiesDto = append(proxiesDto, ProxyDto{Index: i + 1, Url: proxies[i]})
		}

		j, _ := json.MarshalIndent(proxiesDto, "", "  ")
		_, err = io.WriteString(w, string(j))
		if err != nil {
			response := "Response serialization failed"
			http.Error(w, response, http.StatusInternalServerError)
		}
	}
}

func (controller *Controller) createProxies(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /post proxies request\n")
	decoder := json.NewDecoder(r.Body)
	var urls []string
	err := decoder.Decode(&urls)
	if err != nil {
		response := "Can't parse body"
		http.Error(w, response, http.StatusInternalServerError)
	} else {
		err := controller.proxyManager.WriteProxies(urls)
		if err != nil {
			log.Error(err)
			response := "Can't write proxies"
			http.Error(w, response, http.StatusInternalServerError)
		}
	}
}

func (controller *Controller) deleteProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /delete proxy request\n")
	vars := mux.Vars(r)
	strIndex, ok := vars["index"]
	if !ok {
		log.Error("Index is missing in parameters")
		response := "Missed index parameter"
		http.Error(w, response, http.StatusBadRequest)
	} else {
		index, err := strconv.Atoi(strIndex)
		if err != nil {
			log.Error("Failed to convert index to int")
			response := "Index must be number"
			http.Error(w, response, http.StatusBadRequest)
		}
		err = controller.proxyManager.DeleteProxy(index - 1)
		if err != nil {
			log.Error(err)
			response := "Can't write proxies"
			http.Error(w, response, http.StatusInternalServerError)
		}
	}
}
