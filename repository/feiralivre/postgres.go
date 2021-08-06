package feiralivre

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bgildson/unico-challenge/entity"
)

const (
	QueryByID = `
SELECT
    id,
    latitude,
    longitude,
    setor_censitario,
    area_ponderacao,
    codigo_distrito,
    distrito,
    codigo_subprefeitura,
    subprefeitura,
    regiao5,
    regiao8,
    nome_feira,
    registro,
    logradouro,
    numero,
    bairro,
    referencia,
    created_at,
    updated_at
FROM feira_livre
WHERE id = $1;`
	QueryCreate = `
INSERT INTO feira_livre
    (latitude, longitude, setor_censitario, area_ponderacao, codigo_distrito, distrito, codigo_subprefeitura, subprefeitura, regiao5, regiao8, nome_feira, registro, logradouro, numero, bairro, referencia)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, created_at, updated_at;`
	QueryCreateOrUpdate = `
INSERT INTO feira_livre
    (id, latitude, longitude, setor_censitario, area_ponderacao, codigo_distrito, distrito, codigo_subprefeitura, subprefeitura, regiao5, regiao8, nome_feira, registro, logradouro, numero, bairro, referencia)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
ON CONFLICT (id)
DO
    UPDATE SET
        latitude = $18,
        longitude = $19,
        setor_censitario = $20,
        area_ponderacao = $21,
        codigo_distrito = $22,
        distrito = $23,
        codigo_subprefeitura = $24,
        subprefeitura = $25,
        regiao5 = $26,
        regiao8 = $27,
        nome_feira = $28,
        registro = $29,
        logradouro = $30,
        numero = $31,
        bairro = $32,
        referencia = $33,
        updated_at = NOW()
RETURNING id, created_at, updated_at;`
	QueryUpdate = `
UPDATE
    feira_livre
SET
    latitude = $1,
    longitude = $2,
    setor_censitario = $3,
    area_ponderacao = $4,
    codigo_distrito = $5,
    distrito = $6,
    codigo_subprefeitura = $7,
    subprefeitura = $8,
    regiao5 = $9,
    regiao8 = $10,
    nome_feira = $11,
    registro = $12,
    logradouro = $13,
    numero = $14,
    bairro = $15,
    referencia = $16,
	updated_at = NOW()
WHERE
    id = $17
RETURNING updated_at;`
	QueryRemove = `
DELETE FROM
    feira_livre
WHERE
    id = $1;`
	QuerySyncPK = `SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"feira_livre"', 'id')), (SELECT (MAX("id") + 1) FROM "feira_livre"), FALSE);`
)

func ParseQueryParamsToQuery(qp QueryParams) string {
	result := `
SELECT
    id,
    latitude,
    longitude,
    setor_censitario,
    area_ponderacao,
    codigo_distrito,
    distrito,
    codigo_subprefeitura,
    subprefeitura,
    regiao5,
    regiao8,
    nome_feira,
    registro,
    logradouro,
    numero,
    bairro,
    referencia,
    created_at,
    updated_at
FROM feira_livre`

	where := `
WHERE`

	result += where

	next := 1

	if qp.Distrito != "" {
		result += `
    distrito ILIKE '%' || $` + fmt.Sprint(next) + ` || '%' AND`
		next++
	}

	if qp.Regiao5 != "" {
		result += `
    regiao5 ILIKE '%' || $` + fmt.Sprint(next) + ` || '%' AND`
		next++
	}

	if qp.NomeFeira != "" {
		result += `
    nome_feira ILIKE '%' || $` + fmt.Sprint(next) + ` || '%' AND`
		next++
	}

	if qp.Bairro != "" {
		result += `
    bairro ILIKE '%' || $` + fmt.Sprint(next) + ` || '%' AND`
		next++
	}

	result = strings.TrimRight(result, where)
	result = strings.TrimRight(result, " AND")

	result += `
OFFSET $` + fmt.Sprint(next) + `
LIMIT $` + fmt.Sprint(next+1) + `;`

	return result
}

func ParseQueryParamsToArgs(qp QueryParams) (args []interface{}) {
	if qp.Distrito != "" {
		args = append(args, qp.Distrito)
	}

	if qp.Regiao5 != "" {
		args = append(args, qp.Regiao5)
	}

	if qp.NomeFeira != "" {
		args = append(args, qp.NomeFeira)
	}

	if qp.Bairro != "" {
		args = append(args, qp.Bairro)
	}

	args = append(args, qp.Pagination.Offset)
	args = append(args, qp.Pagination.Limit)

	return
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r postgresRepository) GetByQueryParams(qp QueryParams) ([]entity.FeiraLivre, error) {
	q := ParseQueryParamsToQuery(qp)
	a := ParseQueryParamsToArgs(qp)
	res, err := r.db.Query(q, a...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	result := []entity.FeiraLivre{}
	for res.Next() {
		f := entity.FeiraLivre{}
		err := res.Scan(
			&f.ID,
			&f.Latitude,
			&f.Longitude,
			&f.SetorCensitario,
			&f.AreaPonderacao,
			&f.CodigoDistrito,
			&f.Distrito,
			&f.CodigoSubprefeitura,
			&f.Subprefeitura,
			&f.Regiao5,
			&f.Regiao8,
			&f.NomeFeira,
			&f.Registro,
			&f.Logradouro,
			&f.Numero,
			&f.Bairro,
			&f.Referencia,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}

	return result, nil
}

func (r postgresRepository) GetByID(id int) (*entity.FeiraLivre, error) {
	row := r.db.QueryRow(QueryByID, id)
	var f entity.FeiraLivre
	err := row.Scan(
		&f.ID,
		&f.Latitude,
		&f.Longitude,
		&f.SetorCensitario,
		&f.AreaPonderacao,
		&f.CodigoDistrito,
		&f.Distrito,
		&f.CodigoSubprefeitura,
		&f.Subprefeitura,
		&f.Regiao5,
		&f.Regiao8,
		&f.NomeFeira,
		&f.Registro,
		&f.Logradouro,
		&f.Numero,
		&f.Bairro,
		&f.Referencia,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (r postgresRepository) Create(feiraLive entity.FeiraLivre) (*entity.FeiraLivre, error) {
	err := r.db.
		QueryRow(
			QueryCreate,
			feiraLive.Latitude,
			feiraLive.Longitude,
			feiraLive.SetorCensitario,
			feiraLive.AreaPonderacao,
			feiraLive.CodigoDistrito,
			feiraLive.Distrito,
			feiraLive.CodigoSubprefeitura,
			feiraLive.Subprefeitura,
			feiraLive.Regiao5,
			feiraLive.Regiao8,
			feiraLive.NomeFeira,
			feiraLive.Registro,
			feiraLive.Logradouro,
			feiraLive.Numero,
			feiraLive.Bairro,
			feiraLive.Referencia,
		).
		Scan(
			&feiraLive.ID,
			&feiraLive.CreatedAt,
			&feiraLive.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &feiraLive, nil
}

func (r postgresRepository) CreateOrUpdate(feiraLive entity.FeiraLivre) (*entity.FeiraLivre, error) {
	err := r.db.
		QueryRow(
			QueryCreateOrUpdate,
			feiraLive.ID,
			feiraLive.Latitude,
			feiraLive.Longitude,
			feiraLive.SetorCensitario,
			feiraLive.AreaPonderacao,
			feiraLive.CodigoDistrito,
			feiraLive.Distrito,
			feiraLive.CodigoSubprefeitura,
			feiraLive.Subprefeitura,
			feiraLive.Regiao5,
			feiraLive.Regiao8,
			feiraLive.NomeFeira,
			feiraLive.Registro,
			feiraLive.Logradouro,
			feiraLive.Numero,
			feiraLive.Bairro,
			feiraLive.Referencia,
			feiraLive.Latitude,
			feiraLive.Longitude,
			feiraLive.SetorCensitario,
			feiraLive.AreaPonderacao,
			feiraLive.CodigoDistrito,
			feiraLive.Distrito,
			feiraLive.CodigoSubprefeitura,
			feiraLive.Subprefeitura,
			feiraLive.Regiao5,
			feiraLive.Regiao8,
			feiraLive.NomeFeira,
			feiraLive.Registro,
			feiraLive.Logradouro,
			feiraLive.Numero,
			feiraLive.Bairro,
			feiraLive.Referencia,
		).
		Scan(
			&feiraLive.ID,
			&feiraLive.CreatedAt,
			&feiraLive.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &feiraLive, nil
}

func (r postgresRepository) Update(id int, feiraLive entity.FeiraLivre) (*entity.FeiraLivre, error) {
	err := r.db.
		QueryRow(
			QueryUpdate,
			feiraLive.Latitude,
			feiraLive.Longitude,
			feiraLive.SetorCensitario,
			feiraLive.AreaPonderacao,
			feiraLive.CodigoDistrito,
			feiraLive.Distrito,
			feiraLive.CodigoSubprefeitura,
			feiraLive.Subprefeitura,
			feiraLive.Regiao5,
			feiraLive.Regiao8,
			feiraLive.NomeFeira,
			feiraLive.Registro,
			feiraLive.Logradouro,
			feiraLive.Numero,
			feiraLive.Bairro,
			feiraLive.Referencia,
			id,
		).
		Scan(
			&feiraLive.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &feiraLive, nil
}

func (r postgresRepository) Remove(id int) error {
	_, err := r.db.Exec(QueryRemove, id)
	return err
}

func (r postgresRepository) SyncPK() error {
	_, err := r.db.Exec(QuerySyncPK)
	return err
}
