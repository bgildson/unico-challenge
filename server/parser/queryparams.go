package parser

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bgildson/unico-challenge/repository/feiralivre"
)

type QueryParamsParser func(*fiber.Ctx) feiralivre.QueryParams

func NewQueryParamsParser(defaultLimit, maxLimit int) QueryParamsParser {
	return func(c *fiber.Ctx) feiralivre.QueryParams {
		var queryParams feiralivre.QueryParams
		c.QueryParser(&queryParams)

		if queryParams.Pagination.Limit < 1 {
			queryParams.Pagination.Limit = defaultLimit
		} else if queryParams.Pagination.Limit > maxLimit {
			queryParams.Pagination.Limit = maxLimit
		}

		if queryParams.Pagination.Offset < 0 {
			queryParams.Pagination.Offset = 0
		}

		return queryParams
	}
}
