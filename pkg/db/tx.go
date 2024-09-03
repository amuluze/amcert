// Package db
// Date       : 2024/9/2 10:36
// Author     : Amu
// Description:
package db

import (
	"encoding/json"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"strings"
)

var (
	defaultBucket  = []byte("default")
	Separator      = "/"
	ErrKeyNotFound = errors.New("Key is not found")
	ErrKeyInvalid  = errors.New("Key is invalid")
)

type Tx struct {
	tx *bolt.Tx
}

func (t *Tx) Keys(prefix string, recursive bool) ([]string, error) {
	bk := t.tx.Bucket(defaultBucket)
	if bk == nil {
		return nil, ErrKeyNotFound
	}
	var rs []string
	err := bk.ForEach(func(k, v []byte) error {
		key := string(k)
		if strings.HasPrefix(key, prefix) {
			tr := key[len(prefix):]
			if !strings.Contains(tr, Separator) || recursive {
				rs = append(rs, key)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (t *Tx) Get(key string) ([]byte, error) {
	bk := t.tx.Bucket(defaultBucket)
	if bk == nil {
		return nil, ErrKeyNotFound
	}
	v := bk.Get([]byte(key))
	if v == nil {
		return nil, ErrKeyInvalid
	}
	return v, nil
}

func (t *Tx) GetString(key string) (string, error) {
	v, err := t.Get(key)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (t *Tx) GetJson(key string, out interface{}) error {
	v, err := t.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, out)
}

func (t *Tx) Put(key string, value []byte) error {
	bk := t.tx.Bucket(defaultBucket)
	if bk == nil {
		return ErrKeyNotFound
	}
	return bk.Put([]byte(key), value)
}

func (t *Tx) PutJson(key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return t.Put(key, raw)
}

func (t *Tx) PutString(key string, value string) error {
	return t.Put(key, []byte(value))
}

func (t *Tx) Delete(k string) error {
	bk := t.tx.Bucket(defaultBucket)
	if bk == nil {
		return ErrKeyNotFound
	}
	return bk.Delete([]byte(k))
}
