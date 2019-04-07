package token

import (
	"math/rand"
	"sync"
	"time"
)

type KeyInfo struct {
	keys  map[int64]string
	mutex sync.Mutex
}

var keys *KeyInfo

const (
	base      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	SecondDay = 24 * 3600
)

func init() {
	keys = &KeyInfo{
		keys: make(map[int64]string),
	}
	keys.InitKeys()
	go keys.UpdateKeyLoop()
}

func (this *KeyInfo) InitKeys() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.keys = make(map[int64]string)
	for i := 0; i <= 2; i++ {
		ts := time.Now().Unix() - int64(2*i*SecondDay)
		value := getRandString(10)
		this.keys[ts] = value
	}
}

func (this *KeyInfo) GetKey() (int64, string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ts := int64(0)
	serial := ""
	for key, value := range this.keys {
		if key > ts {
			ts = key
			serial = value
		}
	}
	return ts, serial
}

func (this *KeyInfo) GetKeyByTs(ts int64) (string, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	key, ok := this.keys[ts]
	return key, ok
}

func (this *KeyInfo) UpdateKeyLoop() {
	for {
		select {
		case <-time.After(2 * time.Second):
			this.UpdateKey()
		}
	}
}

func (this *KeyInfo) UpdateKey() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tmpKey := make([]int64, 0)
	tdTs := time.Now().Unix() - int64(2*SecondDay)
	for ts, _ := range this.keys {
		if ts <= tdTs {
			tmpKey = append(tmpKey, ts)
		}
	}

	for _, key := range tmpKey {
		delete(this.keys, key)
	}

	keyLen := len(this.keys)
	for i := 1; i <= 3-keyLen; i++ {
		ts := time.Now().Unix() - int64((i-1)*SecondDay)
		value := getRandString(10)
		this.keys[ts] = value
	}
}

func getRandString(length int) string {
	array := make([]byte, length)
	diclen := len(base)
	for i := 0; i < length; i++ {
		array[i] = base[rand.Int()%diclen]
	}
	return string(array)
}

func GetKey() (int64, string) {
	return keys.GetKey()
}

func GetKeyByTs(ts int64) (string, bool) {
	return keys.GetKeyByTs(ts)
}
