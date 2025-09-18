package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNoSuchNotification = errors.New("there is not notification with such id")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}
