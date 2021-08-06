package feiralivre

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bgildson/unico-challenge/entity"
	"github.com/bgildson/unico-challenge/repository/feiralivre"
	"github.com/bgildson/unico-challenge/server/parser"
	"github.com/bgildson/unico-challenge/server/response"
)

type Controller struct {
	feiralivreRepo    feiralivre.Repository
	queryParamsParser parser.QueryParamsParser
}

func New(feiralivreRepo feiralivre.Repository, queryParamsParser parser.QueryParamsParser) *Controller {
	return &Controller{
		feiralivreRepo:    feiralivreRepo,
		queryParamsParser: queryParamsParser,
	}
}

func (c Controller) Register(app *fiber.App, path string) {
	app.Get(path, c.GetByQueryParams)
	app.Get(path+"/:id", c.GetByID)
	app.Post(path, c.Create)
	app.Put(path+"/:id", c.Update)
	app.Delete(path+"/:id", c.Remove)
}

func (c Controller) GetByQueryParams(ctx *fiber.Ctx) error {
	queryParams := c.queryParamsParser(ctx)

	res, err := c.feiralivreRepo.GetByQueryParams(queryParams)
	if err != nil {
		logrus.Errorf("could not query with %+v: %v", queryParams, err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(
				response.Generic{
					Code:    http.StatusInternalServerError,
					Message: "could not query",
				},
			)
	}

	return ctx.JSON(res)
}

func (c Controller) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("could not parse '%v' as id: %v", idParam, err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(
				response.Generic{
					Code:    http.StatusBadRequest,
					Message: "invalid id",
				},
			)
	}

	res, err := c.feiralivreRepo.GetByID(id)
	if err == sql.ErrNoRows {
		logrus.Errorf("could not get by id, feiralivre %d does not exist: %v", id, err)
		return ctx.
			Status(http.StatusNotFound).
			JSON(response.Generic{
				Code:    http.StatusNotFound,
				Message: "not found",
			})
	}
	if err != nil {
		logrus.Errorf("could not get feiralivre %d: %v", id, err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(response.Generic{
				Code:    http.StatusInternalServerError,
				Message: "could not get by id",
			})
	}

	return ctx.JSON(res)
}

func (c Controller) Create(ctx *fiber.Ctx) error {
	var fl entity.FeiraLivre
	if err := json.Unmarshal(ctx.Body(), &fl); err != nil {
		logrus.Errorf("could not parse request body %s: %v", ctx.Body(), err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(response.Generic{
				Code:    http.StatusBadRequest,
				Message: "invalid body",
			})
	}

	res, err := c.feiralivreRepo.Create(fl)
	if err != nil {
		logrus.Errorf("could not create a new feiralivre: %v", err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(response.Generic{
				Code:    http.StatusInternalServerError,
				Message: "could not create",
			})
	}

	return ctx.
		Status(http.StatusCreated).
		JSON(res)
}

func (c Controller) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("could not parse '%v' as id: %v", idParam, err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(
				response.Generic{
					Code:    http.StatusBadRequest,
					Message: "invalid id",
				},
			)
	}

	var fl entity.FeiraLivre
	if err := json.Unmarshal(ctx.Body(), &fl); err != nil {
		logrus.Errorf("could not parse request body %s: %v", ctx.Body(), err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(response.Generic{
				Code:    http.StatusBadRequest,
				Message: "invalid body",
			})
	}

	res, err := c.feiralivreRepo.Update(id, fl)
	if err == sql.ErrNoRows {
		logrus.Errorf("could not update, feiralivre %d does not exist: %v", id, err)
		return ctx.
			Status(http.StatusNotFound).
			JSON(response.Generic{
				Code:    http.StatusNotFound,
				Message: "not found",
			})
	}
	if err != nil {
		logrus.Errorf("could not update feiralivre %d: %v", id, err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(response.Generic{
				Code:    http.StatusInternalServerError,
				Message: "could not update",
			})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(res)
}

func (c Controller) Remove(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("could not parse '%v' as id: %v", idParam, err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(
				response.Generic{
					Code:    http.StatusBadRequest,
					Message: "invalid id",
				},
			)
	}

	if err := c.feiralivreRepo.Remove(id); err != nil {
		logrus.Errorf("could not remove the feiralivre %d: %v", id, err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(response.Generic{
				Code:    http.StatusInternalServerError,
				Message: "could not remove",
			})
	}

	return ctx.SendStatus(http.StatusNoContent)
}
