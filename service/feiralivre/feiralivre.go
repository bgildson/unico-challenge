package feiralivre

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/spf13/afero"

	"github.com/bgildson/unico-challenge/entity"
	"github.com/bgildson/unico-challenge/repository/feiralivre"
)

// Service represents how a feiralivre service should be implemented
type Service interface {
	Import(path string) (message string, err error)
}

type service struct {
	fs   afero.Fs
	repo feiralivre.Repository
}

// New creates a service for feiralivre
func New(fs afero.Fs, repo feiralivre.Repository) Service {
	return &service{
		fs:   fs,
		repo: repo,
	}
}

func (s service) parseColsToFeiraLivre(cols []string) (*entity.FeiraLivre, error) {
	if len(cols) < 17 {
		return nil, fmt.Errorf("the number of cols must be 17 or more, was received %d", len(cols))
	}

	id, err := strconv.Atoi(cols[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse ID column: %v", err)
	}

	long, err := strconv.ParseFloat(cols[1], 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse LONG column: %v", err)
	}

	lat, err := strconv.ParseFloat(cols[2], 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse LAT column: %v", err)
	}

	setcens, err := strconv.Atoi(cols[3])
	if err != nil {
		return nil, fmt.Errorf("could not parse SETCENS column: %v", err)
	}

	areap, err := strconv.Atoi(cols[4])
	if err != nil {
		return nil, fmt.Errorf("could not parse AREAP column: %v", err)
	}

	coddist, err := strconv.Atoi(cols[5])
	if err != nil {
		return nil, fmt.Errorf("could not parse CODDIST column: %v", err)
	}

	codsubpref, err := strconv.Atoi(cols[7])
	if err != nil {
		return nil, fmt.Errorf("could not parse CODSUBPREF column: %v", err)
	}

	return &entity.FeiraLivre{
		ID:                  id,
		Latitude:            lat,
		Longitude:           long,
		SetorCensitario:     setcens,
		AreaPonderacao:      areap,
		CodigoDistrito:      coddist,
		Distrito:            cols[6],
		CodigoSubprefeitura: codsubpref,
		Subprefeitura:       cols[8],
		Regiao5:             cols[9],
		Regiao8:             cols[10],
		NomeFeira:           cols[11],
		Registro:            cols[12],
		Logradouro:          cols[13],
		Numero:              cols[14],
		Bairro:              cols[15],
		Referencia:          cols[16],
	}, nil
}

func (s service) readFile(path string, flChan chan<- *entity.FeiraLivre, errChan chan<- error) {
	defer close(flChan)
	defer close(errChan)

	f, err := s.fs.Open(path)
	if err != nil {
		errChan <- fmt.Errorf("could not open csv file: %v", err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// skip header
	if _, err := reader.Read(); err != nil {
		errChan <- fmt.Errorf("could not read csv header: %v", err)
		return
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errChan <- fmt.Errorf("could not read csv row: %v", err)
			continue
		}
		fl, err := s.parseColsToFeiraLivre(row)
		if err != nil {
			errChan <- fmt.Errorf("could not parse row to feiralivre: %v", err)
			continue
		}
		flChan <- fl
	}
}

// Import implements the csv import operation
func (s service) Import(path string) (string, error) {
	// read
	flChan := make(chan *entity.FeiraLivre)
	readErrChan := make(chan error)
	go s.readFile(path, flChan, readErrChan)

	// count read errors
	var readErrCount int64
	go func() {
		for range readErrChan {
			atomic.AddInt64(&readErrCount, 1)
		}
	}()

	// import (worker pool)
	poolSize := 8
	wg := sync.WaitGroup{}
	wg.Add(poolSize)
	var flCount int64
	var importErrCount int64
	for i := 0; i < poolSize; i++ {
		go func() {
			defer wg.Done()
			for fl := range flChan {
				atomic.AddInt64(&flCount, 1)
				if _, err := s.repo.CreateOrUpdate(*fl); err != nil {
					atomic.AddInt64(&importErrCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// sync feira_livre table pk
	var extraErr string
	if err := s.repo.SyncPK(); err != nil {
		extraErr = fmt.Sprintf("could not sync feira_livre table pk: %v\n", err)
	}

	return fmt.Sprintf(
		"Import finished! Read %d registers, %d imported and %d errors.\n%s",
		flCount+atomic.LoadInt64(&readErrCount),
		flCount-atomic.LoadInt64(&importErrCount),
		atomic.LoadInt64(&readErrCount)+atomic.LoadInt64(&importErrCount),
		extraErr,
	), nil
}
