package feiralivre

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	"github.com/bgildson/unico-challenge/entity"
	"github.com/bgildson/unico-challenge/repository/feiralivre"
)

func TestServiceParseColsToFeiraLivre(t *testing.T) {
	testCases := []struct {
		name     string
		in       []string
		out      *entity.FeiraLivre
		hasError bool
	}{
		{
			name:     "when the input does not have the minimum size",
			in:       []string{},
			hasError: true,
		},
		{
			name:     "when the id can not be parsed",
			in:       []string{"a", "-46548146", "-23568390", "355030885000019", "3550308005040", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the long can not be parsed",
			in:       []string{"1", "a", "-23568390", "355030885000019", "3550308005040", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the lat can not be parsed",
			in:       []string{"1", "-46548146", "a", "355030885000019", "3550308005040", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the setcens can not be parsed",
			in:       []string{"1", "-46548146", "-23568390", "a", "3550308005040", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the areap can not be parsed",
			in:       []string{"1", "-46548146", "-23568390", "355030885000019", "a", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the coddist can not be parsed",
			in:       []string{"1", "-46548146", "-23568390", "355030885000019", "3550308005040", "a", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name:     "when the codsubpref can not be parsed",
			in:       []string{"1", "-46548146", "-23568390", "355030885000019", "3550308005040", "87", "VILA FORMOSA", "a", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			hasError: true,
		},
		{
			name: "when success",
			in:   []string{"1", "-46548146", "-23568390", "355030885000019", "3550308005040", "87", "VILA FORMOSA", "26", "ARICANDUVA", "Leste", "Leste 1", "PRAÇA LEÃO X", "7216-8", "RUA CODAJÁS", "45", "VILA FORMOSA", "PRAÇA MARECHAL LEITE BANDEIRA"},
			out: &entity.FeiraLivre{
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
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service{}
			res, err := s.parseColsToFeiraLivre(tc.in)
			if tc.hasError && err == nil {
				t.Error("was expecting an error, but returns nil")
			}
			if !tc.hasError && err != nil {
				t.Errorf("was not expecting an error, but returns: %v", err)
			}
			if !reflect.DeepEqual(tc.out, res) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, res)
			}
		})
	}
}

func TestServiceReadFile(t *testing.T) {
	headersLine := "ID,LONG,LAT,SETCENS,AREAP,CODDIST,DISTRITO,CODSUBPREF,SUBPREF,REGIAO5,REGIAO8,NOME_FEIRA,REGISTRO,LOGRADOURO,NUMERO,BAIRRO,REFERENCIA"
	bodyLine := "1,-46548146,-23568390,355030885000019,3550308005040,87,VILA FORMOSA,26,ARICANDUVA,Leste,Leste 1,PRAÇA LEÃO X,7216-8,RUA CODAJÁS,45,VILA FORMOSA,PRAÇA MARECHAL LEITE BANDEIRA"
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
	testCases := []struct {
		name          string
		inPathToWrite string
		inPathToRead  string
		inContent     string
		outFls        []entity.FeiraLivre
		outErr        []error
	}{
		{
			name:         "when the file does not exists",
			inPathToRead: "/does-not-exists.csv",
			outErr:       []error{errors.New("file does not exists")},
		},
		{
			name:          "when the file contains an invalid content in the header",
			inPathToWrite: "/my.csv",
			inPathToRead:  "/my.csv",
			inContent:     `"`,
			outErr:        []error{errors.New("invalid content in the header")},
		},
		{
			name:          "when the file contains an invalid content in the body",
			inPathToWrite: "/my.csv",
			inPathToRead:  "/my.csv",
			inContent:     headersLine + "\n" + `"`,
			outErr:        []error{errors.New("invalid content in the body")},
		},
		{
			name:          "when can not parse the row",
			inPathToWrite: "/my.csv",
			inPathToRead:  "/my.csv",
			inContent:     headersLine + "\n,,,,,,,,,,,,,,,,\n",
			outErr:        []error{errors.New("can not parse the row")},
		},
		{
			name:          "when success",
			inPathToWrite: "/my.csv",
			inPathToRead:  "/my.csv",
			inContent:     headersLine + "\n" + bodyLine,
			outFls:        []entity.FeiraLivre{fl},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afero.WriteFile(fs, tc.inPathToWrite, []byte(tc.inContent), 0644)
			s := service{fs: fs}

			flChan := make(chan *entity.FeiraLivre, 1)
			errChan := make(chan error, 1)

			go s.readFile(tc.inPathToRead, flChan, errChan)

			var fls []entity.FeiraLivre
			for fl := range flChan {
				fls = append(fls, *fl)
			}
			if !reflect.DeepEqual(tc.outFls, fls) {
				t.Errorf("was expecting %+v, but returns %+v", tc.outFls, fls)
			}

			var errs []error
			for err := range errChan {
				errs = append(errs, err)
			}
			if len(tc.outErr) != len(errs) {
				t.Errorf("was expecting %d errors, but returns %d errors (%+v)", len(tc.outErr), len(errs), errs)
			}
		})
	}
}

func TestServiceImport(t *testing.T) {
	headersLine := "ID,LONG,LAT,SETCENS,AREAP,CODDIST,DISTRITO,CODSUBPREF,SUBPREF,REGIAO5,REGIAO8,NOME_FEIRA,REGISTRO,LOGRADOURO,NUMERO,BAIRRO,REFERENCIA"
	bodyLine := "1,-46548146,-23568390,355030885000019,3550308005040,87,VILA FORMOSA,26,ARICANDUVA,Leste,Leste 1,PRAÇA LEÃO X,7216-8,RUA CODAJÁS,45,VILA FORMOSA,PRAÇA MARECHAL LEITE BANDEIRA"
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
	testCases := []struct {
		name       string
		setupMocks func(fs afero.Fs, a *feiralivre.MockRepository)
		in         string
		out        string
	}{
		{
			name: "when the file is empty",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine), 0644)
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 0 registers, 0 imported and 0 errors.\n",
		},
		{
			name: "when the file has one row with error when read",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n,,,,,,,,,,,,,,,,"), 0644)
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 1 registers, 0 imported and 1 errors.\n",
		},
		{
			name: "when the file has one row with error when import",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n"+bodyLine), 0644)
				repo.
					EXPECT().
					CreateOrUpdate(fl).
					Return(nil, errors.New("unexpected error"))
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 1 registers, 0 imported and 1 errors.\n",
		},
		{
			name: "when the file has one row ok",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n"+bodyLine), 0644)
				repo.
					EXPECT().
					CreateOrUpdate(fl).
					Return(&fl, nil)
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 1 registers, 1 imported and 0 errors.\n",
		},
		{
			name: "when the file has one row ok and occur an error when sync",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n"+bodyLine), 0644)
				repo.
					EXPECT().
					CreateOrUpdate(fl).
					Return(&fl, nil)
				repo.
					EXPECT().
					SyncPK().
					Return(errors.New("unexpected error"))
			},
			in:  "/my.csv",
			out: "Import finished! Read 1 registers, 1 imported and 0 errors.\ncould not sync feira_livre table pk: unexpected error\n",
		},
		{
			name: "when the file has one row ok and one with error",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n"+bodyLine+"\n,,,,,,,,,,,,,,,,"), 0644)
				repo.
					EXPECT().
					CreateOrUpdate(fl).
					Return(&fl, nil)
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 2 registers, 1 imported and 1 errors.\n",
		},
		{
			name: "when the file has one row with error and one ok",
			setupMocks: func(fs afero.Fs, repo *feiralivre.MockRepository) {
				afero.WriteFile(fs, "/my.csv", []byte(headersLine+"\n,,,,,,,,,,,,,,,,\n"+bodyLine), 0644)
				repo.
					EXPECT().
					CreateOrUpdate(fl).
					Return(&fl, nil)
				repo.
					EXPECT().
					SyncPK().
					Return(nil)
			},
			in:  "/my.csv",
			out: "Import finished! Read 2 registers, 1 imported and 1 errors.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			ctrl := gomock.NewController(t)
			repo := feiralivre.NewMockRepository(ctrl)
			tc.setupMocks(fs, repo)
			svc := New(fs, repo)

			msg, err := svc.Import(tc.in)

			if msg != tc.out {
				t.Errorf("was expecting:\n%s\nbut returns:\n%s\n", tc.out, msg)
			}
			if err != nil {
				t.Errorf("was expecting an empty error, but returns: %v", err)
			}
		})
	}
}
