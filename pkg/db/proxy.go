// Package db
// Date       : 2024/9/2 10:36
// Author     : Amu
// Description:
package db

import (
	"github.com/pkg/errors"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var defaultTimeout = 10 * time.Second
var Default = &Proxy{Timeout: defaultTimeout}

type Proxy struct {
	Timeout time.Duration

	path string
	db   *bolt.DB
	ref  uint
	mu   sync.Mutex
}

func (p *Proxy) SetPath(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ref != 0 {
		return errors.New("proxy already set")
	}
	p.path = path
	return nil
}

func (p *Proxy) Path() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.path
}

func (p *Proxy) Update(fn func(tx *Tx) error) error {
	err := p.acquire()
	if err != nil {
		return err
	}
	defer p.release()
	return p.db.Update(func(btx *bolt.Tx) error {
		tx := &Tx{tx: btx}
		return fn(tx)
	})
}

func (p *Proxy) View(fn func(tx *Tx) error) error {
	err := p.acquire()
	if err != nil {
		return err
	}
	defer p.release()
	return p.db.View(func(btx *bolt.Tx) error {
		tx := &Tx{tx: btx}
		return fn(tx)
	})
}

func (p *Proxy) acquire() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	timeout := p.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	if p.ref == 0 {
		path := p.path
		if path == "" {
			return errors.New("proxy not set")
		}
		newDB, err := bolt.Open(path, 0600, &bolt.Options{Timeout: timeout})
		if err != nil {
			return err
		}

		err = newDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(defaultBucket)
			return err
		})
		if err != nil {
			return err
		}
		p.db = newDB
	}
	p.ref++
	return nil
}

func (p *Proxy) release() {
	timeout := p.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	for i := 0; i < 3; i++ {
		err := p.doRelease()
		if err != nil {
			return
		}
		time.Sleep(timeout)
	}
}

func (p *Proxy) doRelease() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ref == 0 {
		return nil
	}
	if p.ref == 1 {
		err := p.db.Close()
		if err != nil {
			return err
		}
		p.db = nil
	}
	p.ref--
	return nil
}

func (p *Proxy) Keys(prefix string, recursive bool) ([]string, error) {
	var rs []string
	err := p.View(func(tx *Tx) (rerr error) {
		rs, rerr = tx.Keys(prefix, recursive)
		return nil
	})
	return rs, err
}

func (p *Proxy) Get(key string) ([]byte, error) {
	var rs []byte
	err := p.View(func(tx *Tx) (rerr error) {
		rs, rerr = tx.Get(key)
		return nil
	})
	return rs, err
}

func (p *Proxy) GetString(key string) (string, error) {
	var rs string
	err := p.View(func(tx *Tx) (rerr error) {
		rs, rerr = tx.GetString(key)
		return nil
	})
	return rs, err
}

func (p *Proxy) GetJson(key string, out interface{}) error {
	return p.View(func(tx *Tx) (rerr error) {
		return tx.GetJson(key, out)
	})
}

func (p *Proxy) Put(key string, value []byte) error {
	return p.View(func(tx *Tx) (rerr error) {
		return tx.Put(key, value)
	})
}

func (p *Proxy) PutString(key string, value string) error {
	return p.View(func(tx *Tx) (rerr error) {
		return tx.PutString(key, value)
	})
}

func (p *Proxy) PutJson(key string, value interface{}) error {
	return p.View(func(tx *Tx) (rerr error) {
		return tx.PutJson(key, value)
	})
}

func (p *Proxy) Delete(key string) error {
	return p.Update(func(tx *Tx) error {
		return tx.Delete(key)
	})
}
