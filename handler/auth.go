package handler

import (
	"time"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type TokenInfo struct {
	user     string
	expireAt time.Time
}

var AuthTokenList map[string]TokenInfo

func init() {
	AuthTokenList = make(map[string]TokenInfo)
}

func Check(token string) (ok bool) {
	log.Infof("%+v", AuthTokenList)
	tok, ok := AuthTokenList[token]
	if !ok {
		ok = false
		return
	}
	if tok.expireAt.Before(time.Now()) {
		delete(AuthTokenList, token)
		ok = false
		return
	}
	ok = true
	return
}

func SetAuth(user string) (token string) {
	token = uuid.NewV4().String()
	tinfo := TokenInfo{
		user:     user,
		expireAt: time.Now().Add(time.Hour * 24),
	}
	AuthTokenList[token] = tinfo
	return
}
