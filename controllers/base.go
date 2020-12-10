package controllers

import (
	"myapp/libs"
	"myapp/models"

	"github.com/kataras/iris/v12"
)

func GetCommonListSearch(ctx iris.Context) *models.Search {
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")

	relation := ctx.FormValue("relation")
	return &models.Search{
		Sort:      sort,
		Offset:    offset,
		Limit:     limit,
		OrderBy:   orderBy,
		Relations: models.GetRelations(relation, nil),
	}
}
