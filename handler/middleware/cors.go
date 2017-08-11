package middleware

import "net/http"

func CORS(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Add("Access-Control-Max-Age", "3600")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-Auth-User,X-Auth-Secret,X-NEUPRepair-Token,Content-Type")
	next(w, r)
}
