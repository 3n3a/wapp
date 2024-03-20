package wapp

import (
	// "reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func Test_Action_KV(t *testing.T) {
	kv := NewKV()

	key := "test"
	value := "test-val"

	err := kv.SetString(key, value)
	assert.Equal(t, nil, err)

	val, err := kv.GetString(key)
	assert.Equal(t, nil, err)
	assert.Equal(t, value, val)
}
