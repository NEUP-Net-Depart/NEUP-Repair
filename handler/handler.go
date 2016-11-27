package handler

import (
	"net/http"
	"os"

	"encoding/base64"

	"github.com/NEUP-Net-Depart/NEUP-Repair-backend/model"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
)

// So small, no need middleware just now

func AddOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	aor := AddOrderResponse{}
	resp := Response{}
	er, o := ParseOrder(r)
	if er != nil {
		err = errors.Wrap(er, "add order error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "invalid json"
		resp.Write(w, r)
		return
	}
	oldo, err := model.OrderByStuID(model.GlobalDB.Unsafe(), o.StuID)

	// Still cannot apply
	log.Info(oldo)
	if oldo.ID != -1 {
		resp.Success = false
		resp.Msg = "之前的维修预约还没有完成，请完成之前的申请之后再次申请"
		resp.Write(w, r)
		return
	}
	if err != nil {
		err = errors.Wrap(err, "add order error")
		log.Error(err)
		resp.WriteError(w)
		return
	}

	o.SecretID = uuid.NewV4().String()
	err = o.Insert(model.GlobalDB)
	if err != nil {
		err = errors.Wrap(err, "add order error")
		log.Error(err)
		resp.WriteError(w)
		return
	}

	// Generate QR Code
	secretPath := r.Header.Get("Origin") + "/order.html" + "?secret=" + o.SecretID
	log.Info(secretPath)
	qrpath := "qr." + o.SecretID + ".png"
	err = qrcode.WriteFile(secretPath, qrcode.Medium, 256, qrpath)
	if err != nil {
		err = errors.Wrap(err, "add order error: qr encode error")
		log.Error(err)
		resp.WriteError(w)
		return
	}
	f, err := os.Open(qrpath)
	if err != nil {
		err = errors.Wrap(err, "add order error")
		log.Error(err)
		resp.WriteError(w)
	}
	defer f.Close()
	err = os.Remove(qrpath)
	if err != nil {
		err = errors.Wrap(err, "add order error: cannot remove QRCode file")
		log.Error(err)
	}
	info, err := f.Stat()
	if err != nil {
		err = errors.Wrap(err, "add order error: cannot get file szie")
		log.Error(err)
		resp.WriteError(w)
	}
	encbuf := make([]byte, info.Size())
	f.Read(encbuf)
	encstr := base64.StdEncoding.EncodeToString(encbuf)
	if err != nil {
		err = errors.Wrap(err, "add order error: base64 encode error")
		log.Error(err)
		resp.WriteError(w)
	}
	resp.Success = true
	aor.QRcode = encstr
	aor.SecretID = o.SecretID
	resp.Data = aor
	resp.Write(w, r)
	return
}

func OrdersByPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func OrdersSecret(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	secretID := ps.ByName("secret")
	o, err := model.OrderBySecret(model.GlobalDB.Unsafe(), secretID)
	if err != nil {
		err = errors.Wrap(err, "get order by secret error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "服务器错误"
		resp.Code = 500
		resp.Write(w, r)
		return
	}
	if o.ID == -1 {
		resp.Success = false
		resp.Msg = "未找到此预约，请检查是否本预约已经完成"
		resp.Code = 400
		resp.Write(w, r)
		return
	}
	resp.Data = o
	resp.Success = true
	resp.Code = 200
	resp.Write(w, r)
	return
}

func Auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func FinishOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	secretID := ps.ByName("secret")
	err := model.UpdateOrderDoneFlagBySecret(model.GlobalDB.Unsafe(), secretID)
	if err != nil {
		err = errors.Wrap(err, "get order by secret error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "更新order状态失败 请检查你的secretID是否正确,如有疑问请联系管理员"
		resp.Code = 500
		resp.Write(w, r)
		return
	}
	resp.Success = true
	resp.Code = 200
	resp.Write(w, r)
	return
}

func RateOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func Comments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
