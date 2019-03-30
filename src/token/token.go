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
	raw    string
	scopes []string
	ts     int64
	hmac   []byte
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

func NewToken(name string, scopes string) (*Token, error) {
	ts, key := GetKey()
	raw := fmt.Sprintf("%s,%s,%d", name, scopes, ts)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(raw))
	hmac := h.Sum(nil)

	return &Token{
		raw:    raw,
		scopes: strings.Split(scopes, " "),
		ts:     ts,
		hmac:   hmac}, nil
}

func NewStaticToken(name string, scopes string, ts int64, hmac []byte) (*Token, error) {
	raw := fmt.Sprintf("%s,%s,%d", name, scopes, ts)
	return &Token{
		raw:    raw,
		scopes: strings.Split(scopes, " "),
		ts:     ts,
		hmac:   hmac}, nil
}

func ParseToken(token string) (*Token, error) {
	data := strings.Split(token, "|")
	if len(data) != 2 {
		return nil, Error{5000, "invalid-token-info."}
	}
	raw := data[0]
	hmacValue := data[1]

	data = strings.Split(raw, ",")
	if len(data) != 3 {
		return nil, Error{5000, "invalid-token-info."}
	}

	ts, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return nil, Error{5000, "invalid-token-info."}
	}

	key, ret := GetKeyByTs(ts)
	if ret == false {
		return nil, Error{5000, "invalid-token-info."}
	}

	encoder := hmac.New(sha256.New, []byte(key))
	encoder.Write([]byte(raw))
	hash, err := hex.DecodeString(hmacValue)
	newhmac := encoder.Sum(nil)
	if err != nil || !hmac.Equal(newhmac, hash) {
		return nil, Error{5001, "failed-checksum."}
	}
	return NewStaticToken(data[0], data[1], ts, hash)
}

func (self *Token) String() string {
	return fmt.Sprintf("%s|%s", self.raw, hex.EncodeToString(self.hmac))
}





