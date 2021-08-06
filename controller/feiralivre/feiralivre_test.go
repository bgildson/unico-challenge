package feiralivre

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	"github.com/bgildson/unico-challenge/entity"
	"github.com/bgildson/unico-challenge/repository/feiralivre"
	"github.com/bgildson/unico-challenge/server/parser"
)

func MatchRoute(app *fiber.App, method string, path string) bool {
	stack := app.Stack()
	for m := range stack {
		for _, r := range stack[m] {
			if r.Method == method && r.Path == path {
				return true
			}
		}
	}
	return false
}

func TestControllerRegister(t *testing.T) {
	path := "/feiras-livres"
	testCases := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "when looking for query route",
			method: http.MethodGet,
			path:   path,
		},
		{
			name:   "when looking for get by id route",
			method: http.MethodGet,
			path:   path + "/:id",
		},
		{
			name:   "when looking for create route",
			method: http.MethodPost,
			path:   path,
		},
		{
			name:   "when looking for update",
			method: http.MethodPut,
			path:   path + "/:id",
		},
		{
			name:   "when looking for remove route",
			method: http.MethodDelete,
			path:   path + "/:id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := New(nil, nil)

			app := fiber.New()

			controller.Register(app, path)

			if !MatchRoute(app, tc.method, tc.path) {
				t.Errorf("there's no route registered in method %s and path %s", tc.method, tc.path)
			}
		})
	}
}

func TestControllerGetByQueryParams(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	path := "/feiras-livres"
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
	testCases := []struct {
		name       string
		setupMocks func(repo *feiralivre.MockRepository)
		outStatus  int
		outBody    interface{}
	}{
		{
			name: "when occur an error in repository",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					GetByQueryParams(gomock.Any()).
					Return(nil, errors.New("unexpected error"))
			},
			outStatus: http.StatusInternalServerError,
			outBody: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "could not query",
			},
		},
		{
			name: "when success",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					GetByQueryParams(gomock.Any()).
					Return([]entity.FeiraLivre{fl}, nil)
			},
			outStatus: http.StatusOK,
			outBody: []map[string]interface{}{
				{
					"id":                   fl.ID,
					"latitude":             fl.Latitude,
					"longitude":            fl.Longitude,
					"setor_censitario":     fl.SetorCensitario,
					"area_ponderacao":      fl.AreaPonderacao,
					"codigo_distrito":      fl.CodigoDistrito,
					"distrito":             fl.Distrito,
					"codigo_subprefeitura": fl.CodigoSubprefeitura,
					"subprefeitura":        fl.Subprefeitura,
					"regiao5":              fl.Regiao5,
					"regiao8":              fl.Regiao8,
					"nome_feira":           fl.NomeFeira,
					"registro":             fl.Registro,
					"logradouro":           fl.Logradouro,
					"numero":               fl.Numero,
					"bairro":               fl.Bairro,
					"referencia":           fl.Referencia,
					"created_at":           fl.CreatedAt,
					"updated_at":           fl.UpdatedAt,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			parser := parser.NewQueryParamsParser(10, 42)
			controller := New(repo, parser)

			app := fiber.New()

			controller.Register(app, path)

			res, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if res.StatusCode != tc.outStatus {
				t.Errorf("was expecting %v, but returns %v", tc.outStatus, res.StatusCode)
			}
			var body interface{}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Errorf("could not decode body: %v", err)
			}

			b1, _ := json.Marshal(tc.outBody)
			b2, _ := json.Marshal(body)
			if string(b1) != string(b2) {
				t.Errorf("was expecting %s, but returns %s", b1, b2)
			}
		})
	}
}

func TestControllerGetByID(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	path := "/feiras-livres"
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
	testCases := []struct {
		name       string
		setupMocks func(repo *feiralivre.MockRepository)
		in         string
		outStatus  int
		outBody    interface{}
	}{
		{
			name:       "when passing an invalid id",
			setupMocks: func(repo *feiralivre.MockRepository) {},
			in:         "a",
			outStatus:  http.StatusBadRequest,
			outBody: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid id",
			},
		},
		{
			name: "when register does not exists",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					GetByID(-1).
					Return(nil, sql.ErrNoRows)
			},
			in:        "-1",
			outStatus: http.StatusNotFound,
			outBody: map[string]interface{}{
				"code":    http.StatusNotFound,
				"message": "not found",
			},
		},
		{
			name: "when repository returns an unexpected error",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					GetByID(1).
					Return(nil, errors.New("unexpected error"))
			},
			in:        "1",
			outStatus: http.StatusInternalServerError,
			outBody: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "could not get by id",
			},
		},
		{
			name: "when success",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					GetByID(1).
					Return(&fl, nil)
			},
			in:        "1",
			outStatus: http.StatusOK,
			outBody: map[string]interface{}{
				"id":                   fl.ID,
				"latitude":             fl.Latitude,
				"longitude":            fl.Longitude,
				"setor_censitario":     fl.SetorCensitario,
				"area_ponderacao":      fl.AreaPonderacao,
				"codigo_distrito":      fl.CodigoDistrito,
				"distrito":             fl.Distrito,
				"codigo_subprefeitura": fl.CodigoSubprefeitura,
				"subprefeitura":        fl.Subprefeitura,
				"regiao5":              fl.Regiao5,
				"regiao8":              fl.Regiao8,
				"nome_feira":           fl.NomeFeira,
				"registro":             fl.Registro,
				"logradouro":           fl.Logradouro,
				"numero":               fl.Numero,
				"bairro":               fl.Bairro,
				"referencia":           fl.Referencia,
				"created_at":           fl.CreatedAt,
				"updated_at":           fl.UpdatedAt,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			controller := New(repo, nil)

			app := fiber.New()

			controller.Register(app, path)

			res, err := app.Test(httptest.NewRequest(http.MethodGet, path+"/"+tc.in, nil))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if res.StatusCode != tc.outStatus {
				t.Errorf("was expecting %v, but returns %v", tc.outStatus, res.StatusCode)
			}
			var body interface{}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Errorf("could not decode body: %v", err)
			}

			b1, _ := json.Marshal(tc.outBody)
			b2, _ := json.Marshal(body)
			if string(b1) != string(b2) {
				t.Errorf("was expecting %s, but returns %s", b1, b2)
			}
		})
	}
}

func TestControllerCreate(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	path := "/feiras-livres"
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
	flAsBody := entity.FeiraLivre{
		Latitude:            fl.Latitude,
		Longitude:           fl.Longitude,
		SetorCensitario:     fl.SetorCensitario,
		AreaPonderacao:      fl.AreaPonderacao,
		CodigoDistrito:      fl.CodigoDistrito,
		Distrito:            fl.Distrito,
		CodigoSubprefeitura: fl.CodigoSubprefeitura,
		Subprefeitura:       fl.Subprefeitura,
		Regiao5:             fl.Regiao5,
		Regiao8:             fl.Regiao8,
		NomeFeira:           fl.NomeFeira,
		Registro:            fl.Registro,
		Logradouro:          fl.Logradouro,
		Numero:              fl.Numero,
		Bairro:              fl.Bairro,
		Referencia:          fl.Referencia,
	}
	body := map[string]interface{}{
		"latitude":             fl.Latitude,
		"longitude":            fl.Longitude,
		"setor_censitario":     fl.SetorCensitario,
		"area_ponderacao":      fl.AreaPonderacao,
		"codigo_distrito":      fl.CodigoDistrito,
		"distrito":             fl.Distrito,
		"codigo_subprefeitura": fl.CodigoSubprefeitura,
		"subprefeitura":        fl.Subprefeitura,
		"regiao5":              fl.Regiao5,
		"regiao8":              fl.Regiao8,
		"nome_feira":           fl.NomeFeira,
		"registro":             fl.Registro,
		"logradouro":           fl.Logradouro,
		"numero":               fl.Numero,
		"bairro":               fl.Bairro,
		"referencia":           fl.Referencia,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testCases := []struct {
		name       string
		setupMocks func(repo *feiralivre.MockRepository)
		in         string
		outStatus  int
		outBody    interface{}
	}{
		{
			name:       "when invalid body",
			setupMocks: func(repo *feiralivre.MockRepository) {},
			in:         ":invalid:",
			outStatus:  http.StatusBadRequest,
			outBody: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid body",
			},
		},
		{
			name: "when repository returns an error",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Create(flAsBody).
					Return(nil, errors.New("unexpected error"))
			},
			in:        string(bodyJSON),
			outStatus: http.StatusInternalServerError,
			outBody: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "could not create",
			},
		},
		{
			name: "when success",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Create(flAsBody).
					Return(&fl, nil)
			},
			in:        string(bodyJSON),
			outStatus: http.StatusCreated,
			outBody: map[string]interface{}{
				"id":                   fl.ID,
				"latitude":             fl.Latitude,
				"longitude":            fl.Longitude,
				"setor_censitario":     fl.SetorCensitario,
				"area_ponderacao":      fl.AreaPonderacao,
				"codigo_distrito":      fl.CodigoDistrito,
				"distrito":             fl.Distrito,
				"codigo_subprefeitura": fl.CodigoSubprefeitura,
				"subprefeitura":        fl.Subprefeitura,
				"regiao5":              fl.Regiao5,
				"regiao8":              fl.Regiao8,
				"nome_feira":           fl.NomeFeira,
				"registro":             fl.Registro,
				"logradouro":           fl.Logradouro,
				"numero":               fl.Numero,
				"bairro":               fl.Bairro,
				"referencia":           fl.Referencia,
				"created_at":           fl.CreatedAt,
				"updated_at":           fl.UpdatedAt,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			controller := New(repo, nil)

			app := fiber.New()

			controller.Register(app, path)

			res, err := app.Test(httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte(tc.in))))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if res.StatusCode != tc.outStatus {
				t.Errorf("was expecting %v, but returns %v", tc.outStatus, res.StatusCode)
			}
			var body interface{}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Errorf("could not decode body: %v", err)
			}

			b1, _ := json.Marshal(tc.outBody)
			b2, _ := json.Marshal(body)
			if string(b1) != string(b2) {
				t.Errorf("was expecting %s, but returns %s", b1, b2)
			}
		})
	}
}

func TestControllerUpdate(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	path := "/feiras-livres"
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
	flAsBody := entity.FeiraLivre{
		Latitude:            fl.Latitude,
		Longitude:           fl.Longitude,
		SetorCensitario:     fl.SetorCensitario,
		AreaPonderacao:      fl.AreaPonderacao,
		CodigoDistrito:      fl.CodigoDistrito,
		Distrito:            fl.Distrito,
		CodigoSubprefeitura: fl.CodigoSubprefeitura,
		Subprefeitura:       fl.Subprefeitura,
		Regiao5:             fl.Regiao5,
		Regiao8:             fl.Regiao8,
		NomeFeira:           fl.NomeFeira,
		Registro:            fl.Registro,
		Logradouro:          fl.Logradouro,
		Numero:              fl.Numero,
		Bairro:              fl.Bairro,
		Referencia:          fl.Referencia,
	}
	body := map[string]interface{}{
		"latitude":             fl.Latitude,
		"longitude":            fl.Longitude,
		"setor_censitario":     fl.SetorCensitario,
		"area_ponderacao":      fl.AreaPonderacao,
		"codigo_distrito":      fl.CodigoDistrito,
		"distrito":             fl.Distrito,
		"codigo_subprefeitura": fl.CodigoSubprefeitura,
		"subprefeitura":        fl.Subprefeitura,
		"regiao5":              fl.Regiao5,
		"regiao8":              fl.Regiao8,
		"nome_feira":           fl.NomeFeira,
		"registro":             fl.Registro,
		"logradouro":           fl.Logradouro,
		"numero":               fl.Numero,
		"bairro":               fl.Bairro,
		"referencia":           fl.Referencia,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testCases := []struct {
		name       string
		setupMocks func(repo *feiralivre.MockRepository)
		inID       string
		inBody     []byte
		outStatus  int
		outBody    interface{}
	}{
		{
			name:       "when invalid id",
			setupMocks: func(repo *feiralivre.MockRepository) {},
			inID:       "a",
			inBody:     bodyJSON,
			outStatus:  http.StatusBadRequest,
			outBody: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid id",
			},
		},
		{
			name:       "when invalid body",
			setupMocks: func(repo *feiralivre.MockRepository) {},
			inID:       fmt.Sprint(fl.ID),
			inBody:     []byte(":invalid:"),
			outStatus:  http.StatusBadRequest,
			outBody: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid body",
			},
		},
		{
			name: "when register does not exists",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Update(fl.ID, flAsBody).
					Return(nil, sql.ErrNoRows)
			},
			inID:      fmt.Sprint(fl.ID),
			inBody:    bodyJSON,
			outStatus: http.StatusNotFound,
			outBody: map[string]interface{}{
				"code":    http.StatusNotFound,
				"message": "not found",
			},
		},
		{
			name: "when repository returns an error",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Update(fl.ID, flAsBody).
					Return(nil, errors.New("unexpected error"))
			},
			inID:      fmt.Sprint(fl.ID),
			inBody:    bodyJSON,
			outStatus: http.StatusInternalServerError,
			outBody: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "could not update",
			},
		},
		{
			name: "when success",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Update(fl.ID, flAsBody).
					Return(&fl, nil)
			},
			inID:      fmt.Sprint(fl.ID),
			inBody:    bodyJSON,
			outStatus: http.StatusOK,
			outBody: map[string]interface{}{
				"id":                   fl.ID,
				"latitude":             fl.Latitude,
				"longitude":            fl.Longitude,
				"setor_censitario":     fl.SetorCensitario,
				"area_ponderacao":      fl.AreaPonderacao,
				"codigo_distrito":      fl.CodigoDistrito,
				"distrito":             fl.Distrito,
				"codigo_subprefeitura": fl.CodigoSubprefeitura,
				"subprefeitura":        fl.Subprefeitura,
				"regiao5":              fl.Regiao5,
				"regiao8":              fl.Regiao8,
				"nome_feira":           fl.NomeFeira,
				"registro":             fl.Registro,
				"logradouro":           fl.Logradouro,
				"numero":               fl.Numero,
				"bairro":               fl.Bairro,
				"referencia":           fl.Referencia,
				"created_at":           fl.CreatedAt,
				"updated_at":           fl.UpdatedAt,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			controller := New(repo, nil)

			app := fiber.New()

			controller.Register(app, path)

			res, err := app.Test(httptest.NewRequest(http.MethodPut, path+"/"+tc.inID, bytes.NewReader(tc.inBody)))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if res.StatusCode != tc.outStatus {
				t.Errorf("was expecting %v, but returns %v", tc.outStatus, res.StatusCode)
			}
			var body interface{}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Errorf("could not decode body: %v", err)
			}

			b1, _ := json.Marshal(tc.outBody)
			b2, _ := json.Marshal(body)
			if string(b1) != string(b2) {
				t.Errorf("was expecting %s, but returns %s", b1, b2)
			}
		})
	}
}

func TestControllerRemove(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)
	path := "/feiras-livres"
	testCases := []struct {
		name       string
		setupMocks func(repo *feiralivre.MockRepository)
		in         string
		outStatus  int
		outBody    interface{}
	}{
		{
			name:       "when invalid id",
			setupMocks: func(repo *feiralivre.MockRepository) {},
			in:         "a",
			outStatus:  http.StatusBadRequest,
			outBody: map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid id",
			},
		},
		{
			name: "when repository returns an error",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Remove(1).
					Return(errors.New("unexpected error"))
			},
			in:        "1",
			outStatus: http.StatusInternalServerError,
			outBody: map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "could not remove",
			},
		},
		{
			name: "when success",
			setupMocks: func(repo *feiralivre.MockRepository) {
				repo.
					EXPECT().
					Remove(1).
					Return(nil)
			},
			in:        "1",
			outStatus: http.StatusNoContent,
			outBody:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(repo)
			controller := New(repo, nil)

			app := fiber.New()

			controller.Register(app, path)

			res, err := app.Test(httptest.NewRequest(http.MethodDelete, path+"/"+tc.in, nil))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if res.StatusCode != tc.outStatus {
				t.Errorf("was expecting %v, but returns %v", tc.outStatus, res.StatusCode)
			}
			if res.StatusCode == http.StatusNoContent {
				return
			}
			var body interface{}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Errorf("could not decode body: %v", err)
			}

			b1, _ := json.Marshal(tc.outBody)
			b2, _ := json.Marshal(body)
			if string(b1) != string(b2) {
				t.Errorf("was expecting %s, but returns %s", b1, b2)
			}
		})
	}
}
