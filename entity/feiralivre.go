package entity

import "time"

type FeiraLivre struct {
	ID                  int       `json:"id"`
	Longitude           float64   `json:"longitude"`
	Latitude            float64   `json:"latitude"`
	SetorCensitario     int       `json:"setor_censitario"`
	AreaPonderacao      int       `json:"area_ponderacao"`
	CodigoDistrito      int       `json:"codigo_distrito"`
	Distrito            string    `json:"distrito"`
	CodigoSubprefeitura int       `json:"codigo_subprefeitura"`
	Subprefeitura       string    `json:"subprefeitura"`
	Regiao5             string    `json:"regiao5"`
	Regiao8             string    `json:"regiao8"`
	NomeFeira           string    `json:"nome_feira"`
	Registro            string    `json:"registro"`
	Logradouro          string    `json:"logradouro"`
	Numero              string    `json:"numero"`
	Bairro              string    `json:"bairro"`
	Referencia          string    `json:"referencia"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
