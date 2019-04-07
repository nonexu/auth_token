package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	raw      string
	clientId string
	userName string
	scopes   map[string]bool
	deviceId string
	ts       int64
	hmac     []byte
}

type Error struct {
	ErrCode int16  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (err Error) Error() string {
	if data, err := json.Marshal(err); err == nil {
		return string(data)
	}
	return ""
}

func NewToken(clientId string, userName string, scopes []string, deviceId string) (*Token, error) {
	ts, key := GetKey()
	scopeStr := strings.Join(scopes, " ")
	raw := fmt.Sprintf("%s,%s,%s,%s,%s,%d", getRandString(10), clientId, userName, scopeStr, deviceId, ts)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(raw))
	hmac := h.Sum(nil)

	tk := &Token{
		raw:      raw,
		clientId: clientId,
		userName: userName,
		scopes:   make(map[string]bool),
		deviceId: deviceId,
		ts:       ts,
		hmac:     hmac}
	for _, sp := range scopes {
		tk.scopes[sp] = true
	}
	return tk, nil
}

func NewStaticToken(randValue string, clientId string, userName string, scopeStr string, deviceId string, ts int64, hmac []byte) (*Token, error) {
	raw := fmt.Sprintf("%s,%s,%s,%s,%s,%d", randValue, clientId, userName, scopeStr, deviceId, ts)
	scopes := strings.Split(scopeStr, " ")

	tk := &Token{
		raw:      raw,
		clientId: clientId,
		userName: userName,
		scopes:   make(map[string]bool),
		deviceId: deviceId,
		ts:       ts,
		hmac:     hmac}
	for _, sp := range scopes {
		tk.scopes[sp] = true
	}

	return tk, nil
}

func ParseToken(token string) (*Token, error) {
	data := strings.Split(token, "|")
	if len(data) != 2 {
		return nil, NewError(ERROR_TOKEN_INVALIDE)
	}

	raw := data[0]
	hmacV := data[1]

	data = strings.Split(raw, ",")
	if len(data) != 6 {
		return nil, NewError(ERROR_TOKEN_INVALIDE)
	}

	ts, err := strconv.ParseInt(data[5], 10, 64)
	if err != nil {
		return nil, NewError(ERROR_TOKEN_INVALIDE)
	}

	key, ret := GetKeyByTs(ts)
	if ret == false {
		return nil, NewError(ERROR_TOKEN_EXPIRED)
	}

	encoder := hmac.New(sha256.New, []byte(key))
	encoder.Write([]byte(raw))
	hash, err := hex.DecodeString(hmacV)
	newhmac := encoder.Sum(nil)
	if err != nil || !hmac.Equal(newhmac, hash) {
		return nil, NewError(ERROR_TOKEN_INVALIDE)
	}
	return NewStaticToken(data[0], data[1], data[2], data[3], data[4], ts, hash)
}

func (self *Token) String() string {
	return fmt.Sprintf("%s|%s", self.raw, hex.EncodeToString(self.hmac))
}

func (self *Token) CheckScopes(scopes []string) bool {
	for _, sp := range scopes {
		if _, ok := self.scopes[sp]; !ok {
			return false
		}
	}
	return true
}
