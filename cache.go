package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

type cache struct {
	locks    *sync.Map
	database *badger.DB
}

func newCache() *cache {
	o := badger.DefaultOptions

	d := fmt.Sprintf("/tmp/badger-%v", uuid.New())
	o.Dir = d
	o.ValueDir = d

	db, err := badger.Open(o)

	if err != nil {
		panic(err)
	}

	c := &cache{&sync.Map{}, db}

	runtime.SetFinalizer(c, func(*cache) {
		db.Close()
		os.RemoveAll(d)
	})

	return c
}

func (c cache) LoadOrStore(s string) (interface{}, func(interface{}), bool) {
	if x, ok := c.load(s); ok {
		return x, nil, true
	}

	g := &sync.WaitGroup{}
	g.Add(1)

	if g, ok := c.locks.LoadOrStore(s, g); ok {
		g.(*sync.WaitGroup).Wait()
		x, _ := c.load(s)

		return x, nil, true
	}

	return nil, func(x interface{}) {
		c.store(s, x)
		g.Done()
		c.locks.Delete(s)
	}, false
}

func (c cache) load(s string) (interface{}, bool) {
	bs := []byte(nil)

	err := c.database.View(func(t *badger.Txn) error {
		i, err := t.Get([]byte(s))

		if err != nil {
			return err
		}

		bs, err = i.Value()

		return err
	})

	if err == badger.ErrKeyNotFound {
		return nil, false
	}

	x := interface{}(nil)
	err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&x)

	if err != nil {
		panic(err)
	}

	return x, true
}

func (c cache) store(s string, x interface{}) {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(x)

	if err != nil {
		panic(err)
	}

	err = c.database.Update(func(t *badger.Txn) error {
		return t.Set([]byte(s), b.Bytes())
	})

	if err != nil {
		panic(err)
	}
}
