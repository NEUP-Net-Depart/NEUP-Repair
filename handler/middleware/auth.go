package middleware

import (
	"net/http"

	"github.com/NEUP-Net-Depart/NEUP-Repair/handler"
	log "github.com/Sirupsen/logrus"
)

const (
	AuthHeader = "X-NEUPRepair-Token"
)

func CheckAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Header.Get(AuthHeader)
	log.Infof("token = %s", token)
	if ok := handler.Check(token); !ok {
		log.Info("Not a authenticated user")
		r.Header.Del("AuthOK")
		next(w, r)
		return
	}
	r.Header.Set("AuthOK", "true")
	next(w, r)
	return
}

func RedirectIfNotAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

}
