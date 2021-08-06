package feiralivre

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/bgildson/unico-challenge/entity"
)

func TestParseQueryParamsToQuery(t *testing.T) {
	testCases := []struct {
		name string
		in   QueryParams
		out  string
	}{
		{
			name: "without query fields",
			in:   QueryParams{},
			out: `
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
OFFSET $1
LIMIT $2;`,
		},
		{
			name: "query by distrito",
			in: QueryParams{
				Distrito: "any",
			},
			out: `
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
WHERE
    distrito ILIKE '%' || $1 || '%'
OFFSET $2
LIMIT $3;`,
		},
		{
			name: "query by regiao5",
			in: QueryParams{
				Regiao5: "any",
			},
			out: `
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
WHERE
    regiao5 ILIKE '%' || $1 || '%'
OFFSET $2
LIMIT $3;`,
		},
		{
			name: "query by nome_feira",
			in: QueryParams{
				NomeFeira: "any",
			},
			out: `
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
WHERE
    nome_feira ILIKE '%' || $1 || '%'
OFFSET $2
LIMIT $3;`,
		},
		{
			name: "query by bairro",
			in: QueryParams{
				Bairro: "any",
			},
			out: `
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
WHERE
    bairro ILIKE '%' || $1 || '%'
OFFSET $2
LIMIT $3;`,
		},
		{
			name: "query by distrito and regiao5",
			in: QueryParams{
				Distrito: "any",
				Regiao5:  "any",
			},
			out: `
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
WHERE
    distrito ILIKE '%' || $1 || '%' AND
    regiao5 ILIKE '%' || $2 || '%'
OFFSET $3
LIMIT $4;`,
		},
		{
			name: "query by all fields",
			in: QueryParams{
				Distrito:  "any",
				Regiao5:   "any",
				NomeFeira: "any",
				Bairro:    "any",
			},
			out: `
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
WHERE
    distrito ILIKE '%' || $1 || '%' AND
    regiao5 ILIKE '%' || $2 || '%' AND
    nome_feira ILIKE '%' || $3 || '%' AND
    bairro ILIKE '%' || $4 || '%'
OFFSET $5
LIMIT $6;`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := ParseQueryParamsToQuery(tc.in); r != tc.out {
				t.Errorf("was expecting:%s\nbut returns: %s", tc.out, r)
			}
		})
	}
}

func TestParseQueryParamsToArgs(t *testing.T) {
	testCases := []struct {
		name string
		in   QueryParams
		out  []interface{}
	}{
		{
			name: "when without query fields",
			in: QueryParams{
				Pagination: Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			out: []interface{}{0, 10},
		},
		{
			name: "when only with distrito field",
			in: QueryParams{
				Distrito: "distrito",
				Pagination: Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			out: []interface{}{"distrito", 0, 10},
		},
		{
			name: "when only with regiao5 field",
			in: QueryParams{
				Regiao5: "regiao5",
				Pagination: Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			out: []interface{}{"regiao5", 0, 10},
		},
		{
			name: "when only with nome_feira field",
			in: QueryParams{
				NomeFeira: "nome_feira",
				Pagination: Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			out: []interface{}{"nome_feira", 0, 10},
		},
		{
			name: "when only with bairro field",
			in: QueryParams{
				Bairro: "bairro",
				Pagination: Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			out: []interface{}{"bairro", 0, 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := ParseQueryParamsToArgs(tc.in); !reflect.DeepEqual(tc.out, r) {
				t.Errorf("was expecting %v, but returns %v", tc.out, r)
			}
		})
	}
}

func TestPostgresRepositoryGetByQueryParams(t *testing.T) {
	fl := entity.FeiraLivre{
		ID:                  1,
		Latitude:            -23568390,
		Longitude:           -46548146,
		SetorCensitario:     355030885000019,
		AreaPonderacao:      3550308005040,
		CodigoDistrito:      87,
		Distrito:            "VILA FORMOSA",
		CodigoSubprefeitura: 26,
		Subprefeitura:       "ARICANDUVA",
		Regiao5:             "Leste",
		Regiao8:             "Leste 1",
		NomeFeira:           "PRAÇA LEÃO X",
		Registro:            "7216-8",
		Logradouro:          "RUA CODAJÁS",
		Numero:              "45",
		Bairro:              "VILA FORMOSA",
		Referencia:          "PRAÇA MARECHAL LEITE BANDEIRA",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	queryParams := QueryParams{
		Distrito: "any",
		Pagination: Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	query := ParseQueryParamsToQuery(queryParams)
	args := ParseQueryParamsToArgs(queryParams)
	var argsDriverValue []driver.Value
	for _, v := range args {
		argsDriverValue = append(argsDriverValue, v)
	}
	cols := []string{"id", "latitude", "longitude", "setor_censitario", "area_ponderacao", "codigo_distrito", "distrito", "codigo_subprefeitura", "subprefeitura", "regiao5", "regiao8", "nome_feira", "registro", "logradouro", "numero", "bairro", "referencia", "created_at", "updated_at"}
	vals := []driver.Value{fl.ID, fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia, fl.CreatedAt, fl.UpdatedAt}
	valsErr := []driver.Value{"", -46548146, -23568390, 355030885000019, 3550308005040, 87, "VILA FORMOSA", 26, "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", 45, "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA", time.Now().String(), time.Now().String()}
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         QueryParams
		out        []entity.FeiraLivre
		hasError   bool
	}{
		{
			name: "when occur an error to apply query",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WithArgs(args).WillReturnError(errors.New("unexpected error"))
			},
			in:       queryParams,
			out:      nil,
			hasError: true,
		},
		{
			name: "when occur an error to scan any row",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(valsErr...)
				mock.ExpectQuery(query).WithArgs(argsDriverValue...).WillReturnRows(rows)
			},
			in:       queryParams,
			out:      nil,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(vals...)
				mock.ExpectQuery(query).WithArgs(argsDriverValue...).WillReturnRows(rows)
			},
			in:       queryParams,
			out:      []entity.FeiraLivre{fl},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			res, err := repo.GetByQueryParams(tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestPostgresRepositoryGetByID(t *testing.T) {
	fl := entity.FeiraLivre{
		ID:                  1,
		Latitude:            -23568390,
		Longitude:           -46548146,
		SetorCensitario:     355030885000019,
		AreaPonderacao:      3550308005040,
		CodigoDistrito:      87,
		Distrito:            "VILA FORMOSA",
		CodigoSubprefeitura: 26,
		Subprefeitura:       "ARICANDUVA",
		Regiao5:             "Leste",
		Regiao8:             "Leste 1",
		NomeFeira:           "PRAÇA LEÃO X",
		Registro:            "7216-8",
		Logradouro:          "RUA CODAJÁS",
		Numero:              "45",
		Bairro:              "VILA FORMOSA",
		Referencia:          "PRAÇA MARECHAL LEITE BANDEIRA",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	cols := []string{"id", "latitude", "longitude", "setor_censitario", "area_ponderacao", "codigo_distrito", "distrito", "codigo_subprefeitura", "subprefeitura", "regiao5", "regiao8", "nome_feira", "registro", "logradouro", "numero", "bairro", "referencia", "created_at", "updated_at"}
	vals := []driver.Value{fl.ID, fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia, fl.CreatedAt, fl.UpdatedAt}
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         int
		out        *entity.FeiraLivre
		hasError   bool
	}{
		{
			name: "when there's not a register",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(QueryByID)).WithArgs(fl.ID).WillReturnError(sql.ErrNoRows)
			},
			in:       fl.ID,
			out:      nil,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(vals...)
				mock.ExpectQuery(regexp.QuoteMeta(QueryByID)).WithArgs(fl.ID).WillReturnRows(rows)
			},
			in:       fl.ID,
			out:      &fl,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			res, err := repo.GetByID(tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestPostgresRepositoryCreate(t *testing.T) {
	fl := entity.FeiraLivre{
		ID:                  1,
		Latitude:            -23568390,
		Longitude:           -46548146,
		SetorCensitario:     355030885000019,
		AreaPonderacao:      3550308005040,
		CodigoDistrito:      87,
		Distrito:            "VILA FORMOSA",
		CodigoSubprefeitura: 26,
		Subprefeitura:       "ARICANDUVA",
		Regiao5:             "Leste",
		Regiao8:             "Leste 1",
		NomeFeira:           "PRAÇA LEÃO X",
		Registro:            "7216-8",
		Logradouro:          "RUA CODAJÁS",
		Numero:              "45",
		Bairro:              "VILA FORMOSA",
		Referencia:          "PRAÇA MARECHAL LEITE BANDEIRA",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	cols := []string{"id", "created_at", "updated_at"}
	vals := []driver.Value{fl.ID, fl.CreatedAt, fl.UpdatedAt}
	argsDriverValue := []driver.Value{fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia}
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         entity.FeiraLivre
		out        *entity.FeiraLivre
		hasError   bool
	}{
		{
			name: "when db returns an error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(QueryCreate)).WithArgs(argsDriverValue...).WillReturnError(sql.ErrConnDone)
			},
			in:       fl,
			out:      nil,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(vals...)
				mock.ExpectQuery(regexp.QuoteMeta(QueryCreate)).WithArgs(argsDriverValue...).WillReturnRows(rows)
			},
			in:       fl,
			out:      &fl,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			res, err := repo.Create(tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestPostgresRepositoryCreateOrUpdate(t *testing.T) {
	fl := entity.FeiraLivre{
		ID:                  1,
		Latitude:            -23568390,
		Longitude:           -46548146,
		SetorCensitario:     355030885000019,
		AreaPonderacao:      3550308005040,
		CodigoDistrito:      87,
		Distrito:            "VILA FORMOSA",
		CodigoSubprefeitura: 26,
		Subprefeitura:       "ARICANDUVA",
		Regiao5:             "Leste",
		Regiao8:             "Leste 1",
		NomeFeira:           "PRAÇA LEÃO X",
		Registro:            "7216-8",
		Logradouro:          "RUA CODAJÁS",
		Numero:              "45",
		Bairro:              "VILA FORMOSA",
		Referencia:          "PRAÇA MARECHAL LEITE BANDEIRA",
	}
	cols := []string{"id", "created_at", "updated_at"}
	vals := []driver.Value{fl.ID, fl.CreatedAt, fl.UpdatedAt}
	argsDriverValue := []driver.Value{fl.ID, fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia, fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia}
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         entity.FeiraLivre
		out        *entity.FeiraLivre
		hasError   bool
	}{
		{
			name: "when db returns an error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(QueryCreateOrUpdate)).WithArgs(argsDriverValue...).WillReturnError(sql.ErrConnDone)
			},
			in:       fl,
			out:      nil,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(vals...)
				mock.ExpectQuery(regexp.QuoteMeta(QueryCreateOrUpdate)).WithArgs(argsDriverValue...).WillReturnRows(rows)
			},
			in:       fl,
			out:      &fl,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			res, err := repo.CreateOrUpdate(tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestPostgresRepositoryUpdate(t *testing.T) {
	fl := entity.FeiraLivre{
		ID:                  1,
		Latitude:            -23568390,
		Longitude:           -46548146,
		SetorCensitario:     355030885000019,
		AreaPonderacao:      3550308005040,
		CodigoDistrito:      87,
		Distrito:            "VILA FORMOSA",
		CodigoSubprefeitura: 26,
		Subprefeitura:       "ARICANDUVA",
		Regiao5:             "Leste",
		Regiao8:             "Leste 1",
		NomeFeira:           "PRAÇA LEÃO X alterado",
		Registro:            "7216-8",
		Logradouro:          "RUA CODAJÁS",
		Numero:              "45",
		Bairro:              "VILA FORMOSA",
		Referencia:          "PRAÇA MARECHAL LEITE BANDEIRA",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	newFl := fl
	newFl.UpdatedAt = newFl.UpdatedAt.Add(10 * time.Minute)
	cols := []string{"updated_at"}
	vals := []driver.Value{newFl.UpdatedAt}
	argsDriverValue := []driver.Value{fl.Latitude, fl.Longitude, fl.SetorCensitario, fl.AreaPonderacao, fl.CodigoDistrito, fl.Distrito, fl.CodigoSubprefeitura, fl.Subprefeitura, fl.Regiao5, fl.Regiao8, fl.NomeFeira, fl.Registro, fl.Logradouro, fl.Numero, fl.Bairro, fl.Referencia, fl.ID}
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         entity.FeiraLivre
		out        *entity.FeiraLivre
		hasError   bool
	}{
		{
			name: "when db returns an error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(QueryUpdate)).WithArgs(argsDriverValue...).WillReturnError(sql.ErrConnDone)
			},
			in:       fl,
			out:      nil,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(vals...)
				mock.ExpectQuery(regexp.QuoteMeta(QueryUpdate)).WithArgs(argsDriverValue...).WillReturnRows(rows)
			},
			in:       fl,
			out:      &newFl,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			res, err := repo.Update(tc.in.ID, tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestPostgresRepositoryRemove(t *testing.T) {
	flID := 1
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		in         int
		hasError   bool
	}{
		{
			name: "when db returns an error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(QueryRemove)).WithArgs(flID).WillReturnError(sql.ErrConnDone)
			},
			in:       flID,
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(QueryRemove)).WithArgs(flID).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			in:       flID,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			err = repo.Remove(tc.in)
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
		})
	}
}

func TestPostgresRepositorySyncPK(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock)
		hasError   bool
	}{
		{
			name: "when db returns an error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(QuerySyncPK)).WillReturnError(sql.ErrConnDone)
			},
			hasError: true,
		},
		{
			name: "when success",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(QuerySyncPK)).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("could not mock sql: %v", err)
			}
			defer db.Close()

			tc.setupMocks(mock)

			repo := NewPostgresRepository(db)

			err = repo.SyncPK()
			if tc.hasError && err == nil {
				t.Errorf("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns %v", err)
			}
		})
	}
}
