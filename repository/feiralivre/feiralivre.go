package feiralivre

import "github.com/bgildson/unico-challenge/entity"

// Repository represents how a feiralivre repository should be implemented
type Repository interface {
	GetByID(int) (*entity.FeiraLivre, error)
	GetByQueryParams(QueryParams) ([]entity.FeiraLivre, error)
	Create(entity.FeiraLivre) (*entity.FeiraLivre, error)
	CreateOrUpdate(feiraLive entity.FeiraLivre) (*entity.FeiraLivre, error)
	Update(int, entity.FeiraLivre) (*entity.FeiraLivre, error)
	Remove(int) error
	SyncPK() error
}

// QueryParams contains the fields that could be used to query a feiralivre
type QueryParams struct {
	Distrito  string `query:"distrito"`
	Regiao5   string `query:"regiao5"`
	NomeFeira string `query:"nome_feira"`
	Bairro    string `query:"bairro"`
	Pagination
}

// Pagination contains the fields that could be used to paginate the query result
type Pagination struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
