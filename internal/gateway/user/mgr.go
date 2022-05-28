package user

import (
	"github.com/go-chi/chi/v5"

	"github.com/ocean-rw/ocean/internal/gateway/db/common"
)

type Mgr struct {
	db common.UserTableIF
}

func New(db common.UserTableIF) (*Mgr, error) {
	return &Mgr{db: db}, nil
}

func (m *Mgr) RegisterRouter(r *chi.Mux) {

}
