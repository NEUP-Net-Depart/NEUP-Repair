package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"encoding/base64"

	"github.com/NEUP-Net-Depart/NEUP-Repair/config"
	"github.com/NEUP-Net-Depart/NEUP-Repair/model"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
)

var GlobalAnnounce string

func init() {
	GlobalAnnounce = "硬件维修工作正常进行~ 欢迎大家报名"
}

// So small, no need middleware just now

func Announce(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	resp.Code = http.StatusOK
	resp.Success = true
	resp.Data = GlobalAnnounce
	resp.Write(w, r)
	return
}

func SetAnnounce(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	if r.Header.Get("AuthOK") != "true" {
		resp.Code = http.StatusUnauthorized
		resp.Success = false
		resp.Msg = "没有权限对公告进行修改哦"
		resp.Write(w, r)
		return
	}
	an, er := ParseString(r)
	if er != nil {
		err := errors.Wrap(er, "update announce error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "invalid json"
		resp.Write(w, r)
		return
	}
	GlobalAnnounce = an.Announce
	resp.Success = true
	resp.Code = http.StatusOK
	resp.Write(w, r)
	return
}

func AddOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	aor := AddOrderResponse{}
	resp := Response{}
	o, er := ParseOrder(r)
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
	log.Info(r.Header)
	log.Info(r.URL)
	log.Info(r.RequestURI)

	// Here we remove all the query params in origin URL
	secretPath := config.GlobalConfig.FrontendRoot + "#/orders/" + o.SecretID
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
		return
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
		return
	}
	encbuf := make([]byte, info.Size())
	f.Read(encbuf)
	encstr := base64.StdEncoding.EncodeToString(encbuf)
	if err != nil {
		err = errors.Wrap(err, "add order error: base64 encode error")
		log.Error(err)
		resp.WriteError(w)
		return
	}
	resp.Success = true
	aor.QRcode = encstr
	aor.SecretID = o.SecretID
	resp.Data = aor
	resp.Write(w, r)
	return
}

func OrdersByPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//var blackList = map[string]bool{}
	perpg := 50
	resp := Response{}
	p := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(p, 10, strconv.IntSize)
	var ol interface{}
	if r.Header.Get("AuthOK") == "true" {
		ol, err = model.OrderPager(model.GlobalDB, int(page), perpg)
	} else {
		ol, err = model.OrderGuardedPager(model.GlobalDB, int(page), perpg)
	}
	if err != nil {
		err = errors.Wrap(err, "get orders by page error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "服务器错误"
		resp.Write(w, r)
		return
	}
	tot, err := model.OrderCount(model.GlobalDB)
	//log.Infof("%+v", ol)
	pgcnt := tot / perpg
	if tot%perpg != 0 {
		pgcnt++
	}
	resp.Success = true
	resp.Data = struct {
		Data      interface{} `json:"data"`
		PageCount int         `json:"page_count"`
		ItemCount int         `json:"item_count"`
		PageOn    int64       `json:"page_on"`
	}{
		Data:      ol,
		PageCount: pgcnt,
		ItemCount: tot,
		PageOn:    page,
	}
	resp.Write(w, r)
	return
}

func OrdersSecret(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	secretID := ps.ByName("secret")
	o, err := model.OrderBySecret(model.GlobalDB, secretID)
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
	resp := Response{}
	if r.Header.Get("AuthOK") == "true" {
		resp.Success = true
		resp.Msg = "您已经登录了"
		resp.Code = http.StatusOK
		resp.Write(w, r)
		return
	}

	au, err := ParseAuth(r)
	if err != nil {
		err = errors.Wrap(err, "add order error")
		log.Error(err)
		resp.Success = false
		resp.Msg = "invalid json"
		resp.Code = http.StatusBadRequest
		resp.Write(w, r)
		return
	}
	// Encoded Password
	enc := sha256.New()
	enc.Write([]byte(au.Password))
	dest := fmt.Sprintf("%x", enc.Sum(nil))
	realpass, err := model.PasswordByName(model.GlobalDB, au.UserName)
	if err != nil {
		log.Error(err)
		resp.Success = false
		resp.Msg = "错误的用户名或者密码"
		resp.Code = http.StatusUnauthorized
		resp.Write(w, r)
		return
	}
	token := SetAuth(au.UserName)
	if dest == realpass {
		resp.Success = true
		resp.Data = token //Login Code
		resp.Code = http.StatusOK
		resp.Write(w, r)
		return
	}
	// Not success
	resp.Success = false
	resp.Msg = "错误的用户名或者密码"
	resp.Code = http.StatusUnauthorized
	resp.Write(w, r)
	return
}

func FinishOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := Response{}
	secretID := ps.ByName("secret")
	if r.Header.Get("AuthOK") == "true" {
		err := model.UpdateOrderDoneFlagBySecret(model.GlobalDB.Unsafe(), secretID)
		if err != nil {
			err = errors.Wrap(err, "get order by secret error")
			log.Error(err)
			resp.Success = false
			resp.Msg = "更新order状态失败 请检查你的secretID是否正确,如有疑问请联系管理员"
			resp.Code = http.StatusInternalServerError
			resp.Write(w, r)
			return
		}
	} else {
		resp.Success = false
		resp.Msg = "没有权限进行此操作"
		resp.Code = http.StatusUnauthorized
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
