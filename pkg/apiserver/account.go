package apiserver

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"net/http"
)

func (s *ApiServer) Account(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookies())
	cookie, err := r.Cookie("token")
	fmt.Println(err)
	if err != nil || cookie.Value == "" {
		response.ErrorForbidden(w)
		return
	}
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Secret), nil
	})
	fmt.Println(err)
	if err != nil || !token.Valid {
		response.ErrorForbidden(w)
		return
	}
	fmt.Println(claims.Subject)
	response.ErrorForbidden(w)
}
