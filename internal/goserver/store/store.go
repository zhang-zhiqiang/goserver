package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	S    *datastore
)

type IStore interface {
	Posts() PostStore
	Users() UserStore
}

type datastore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}

func (ds *datastore) Posts() PostStore {
	return newPosts(ds.db)
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}
