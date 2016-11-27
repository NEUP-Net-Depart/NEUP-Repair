package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NEUP-Net-Depart/NEUP-Repair-backend/model"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type AddOrderResponse struct {
	SecretID string `json:"secret_id"`
	QRcode   string `json:"qrcode"`
}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg,omitempty"`
}

func ParseOrder(r *http.Request) (err error, v model.Order) {
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&v)
	if err != nil {
		err = errors.Wrap(err, "parse order error")
	}
	return
}

func (resp Response) Write(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		err = errors.Wrap(err, "response write error")
		log.Error(err)
		resp.WriteError(w)
	}
	return
}

func (resp Response) WriteError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
