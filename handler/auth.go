package handler

import (
	"time"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type TokenInfo struct {
	user     string
	expireAt time.Time
}

var AuthTokenList map[string]TokenInfo
var mu sync.Mutex

func init() {
	AuthTokenList = make(map[string]TokenInfo)
}

func Check(token string) (ok bool) {
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	AuthTokenList[token] = tinfo
	mu.Unlock()
	return
}
