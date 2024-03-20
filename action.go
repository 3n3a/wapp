package wapp

import (
	// "errors"
	// "log"
	// "strings"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ActionType string

const (
	ActionTypePre ActionType = "pre"
	ActionTypeMain ActionType = "main"
	ActionTypePost ActionType = "post"
)

type KV struct {
	store map[string][]byte
}

func (kv *KV) init() {
	kv.store = make(map[string][]byte)
}

func (kv *KV) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}

	val, ok := kv.store[key]
	if !ok {
		return nil, nil
	}

	return val, nil
}

func (kv *KV) GetString(key string) (string, error) {
	val, err := kv.Get(key)
	if err != nil {
		return "", nil
	}

	return string(val), nil
}

func (kv *KV) GetInt(key string) (int, error) {
	val, err := kv.Get(key)
	if err != nil {
		return 0, nil
	}

	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, nil
	}

	return valInt, nil
}

func (kv *KV) Set(key string, val []byte) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0{
		return nil
	}

	kv.store[key] = []byte(val)

	return nil
}

func (kv *KV) SetString(key string, val string) error {
	return kv.Set(key, []byte(val))
}

func (kv *KV) SetInt(key string, val int) error {
	return kv.Set(key, []byte(strconv.Itoa(val)))
}

func (kv *KV) Delete(key string) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 {
		return nil
	}

	delete(kv.store, key)

	return nil
}

func NewKV() *KV {
	kv := &KV{}
	kv.init()
	return kv
}

type ActionCtx struct {
	*fiber.Ctx

	Store *KV
}

type ActionFunc = func(*ActionCtx) error

// Main Action Container
type Action struct {
	// function that is executed when you know
	f ActionFunc
}

// NewAction creates and initializes a new Action
func NewAction(f ActionFunc) *Action {
	action := &Action{}

	action.f = f

	return action
}