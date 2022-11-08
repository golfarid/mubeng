package auth

import (
	"github.com/dgrijalva/jwt-go"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/api/utils"
	"net/http"
	"strings"
)

type Middleware struct {
	opt *common.Options
}

func New(opt *common.Options) *Middleware {
	log = utils.Logger(opt.Output)
	return &Middleware{opt: opt}
}

func (middleware *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}

		token, err := jwt.ParseWithClaims(
			strings.Replace(tokenString, "Bearer ", "", 1),
			&jwt.StandardClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(middleware.opt.ApiSecret), nil
			},
		)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token"))
			return
		}

		_, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
