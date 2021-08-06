package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"github.com/bgildson/unico-challenge/repository/feiralivre"
)

func TestQueryParams(t *testing.T) {
	paginationInicialOffset := 0
	paginationDefaultLimit := 20
	paginationMaxLimit := 42
	queryParamsParser := NewQueryParamsParser(paginationDefaultLimit, paginationMaxLimit)
	testCases := []struct {
		name     string
		setupCtx func(ctx *fiber.Ctx)
		in       *fiber.Ctx
		out      feiralivre.QueryParams
	}{
		{
			name: "when without params",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: paginationInicialOffset,
					Limit:  paginationDefaultLimit,
				},
			},
		},
		{
			name: "when passing only offset",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service?offset=3")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: 3,
					Limit:  paginationDefaultLimit,
				},
			},
		},
		{
			name: "when passing only negative offset",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service?offset=-1")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: paginationInicialOffset,
					Limit:  paginationDefaultLimit,
				},
			},
		},
		{
			name: "when passing only limit",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service?limit=11")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: paginationInicialOffset,
					Limit:  11,
				},
			},
		},
		{
			name: "when passing only limit lower than 1",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service?limit=-1")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: paginationInicialOffset,
					Limit:  paginationDefaultLimit,
				},
			},
		},
		{
			name: "when passing only limit higher than max",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI(fmt.Sprint("http://app.service?limit=", paginationMaxLimit+1))
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: paginationInicialOffset,
					Limit:  paginationMaxLimit,
				},
			},
		},
		{
			name: "when passing offset and limit",
			setupCtx: func(ctx *fiber.Ctx) {
				ctx.Request().SetRequestURI("http://app.service?offset=3&limit=11")
			},
			out: feiralivre.QueryParams{
				Pagination: feiralivre.Pagination{
					Offset: 3,
					Limit:  11,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			tc.setupCtx(ctx)

			if r := queryParamsParser(ctx); !reflect.DeepEqual(r, tc.out) {
				t.Errorf("was expecting %+v, but returns %+v", tc.out, r)
			}
		})
	}
}
