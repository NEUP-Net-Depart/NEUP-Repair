package handler

import "time"

type TokenInfo struct {
	user     string
	expireAt time.Time
}

var AuthTokenList map[string]TokenInfo

func init() {
	AuthTokenList = make(map[string]TokenInfo)
}

func Check(token string) (ok bool) {
	tok, ok := AuthTokenList[token]
	if !ok {
		ok = false
		return
	}
	if tok.expireAt.After(time.Now()) {
		delete(AuthTokenList, token)
		ok = false
		return
	}
	ok = true
	return
}
