package feiralivre

import "github.com/bgildson/unico-challenge/entity"

type Repository interface {
	GetByID(int) (*entity.FeiraLivre, error)
	GetByQueryParams(QueryParams) ([]entity.FeiraLivre, error)
	Create(entity.FeiraLivre) (*entity.FeiraLivre, error)
	CreateOrUpdate(feiraLive entity.FeiraLivre) (*entity.FeiraLivre, error)
	Update(int, entity.FeiraLivre) (*entity.FeiraLivre, error)
	Remove(int) error
	SyncPK() error
}

type QueryParams struct {
	Distrito  string `query:"distrito"`
	Regiao5   string `query:"regiao5"`
	NomeFeira string `query:"nome_feira"`
	Bairro    string `query:"bairro"`
	Pagination
}

type Pagination struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
