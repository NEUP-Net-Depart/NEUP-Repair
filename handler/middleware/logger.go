package middleware

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
)

func Logger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Infof("Request From %s, Request URI %s, Header %v", r.Header.Get("Origin"), r.RequestURI, r.Header)
	next(w, r)
}
