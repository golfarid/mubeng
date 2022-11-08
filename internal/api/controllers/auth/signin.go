package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/api/utils"
	"net/http"
)

type Controller struct {
	opt *common.Options
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func New(opt *common.Options) *Controller {
	log = utils.Logger(opt.Output)
	return &Controller{opt: opt}
}

func (controller *Controller) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		controller.signIn(w, r)
	}
}

func (controller *Controller) signIn(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if fmt.Sprintf("%s:%s", credentials.Username, credentials.Password) == controller.opt.Auth {
		claims := jwt.StandardClaims{}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := token.SignedString([]byte(controller.opt.ApiSecret))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			w.Header().Set("Authorization", "Bearer "+signedToken)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + signedToken))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Name and password do not match"))
		return
	}
}
