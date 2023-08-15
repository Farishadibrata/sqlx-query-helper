package main

import (
	"fmt"

	"farishadibrata.com/sqlxhelper/helper"
	models "farishadibrata.com/sqlxhelper/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name="%s"`, "localhost", "5432", "postgres", "mysecretpassword", "training", "exapp")
	db, connError := sqlx.Open("postgres", psqlInfo)
	if connError != nil {
		panic(connError)
	}

	another := []models.ExTable{}
	query := helper.NewQuery(helper.QueryHelperParams{
		DB:          db,
		Destination: &another,
		Query:       "SELECT Description, random_number FROM tr.ex_1",
		Where:       []helper.WhereType{{Field: "random_number", Op: helper.Gte, Value: 40}, {Field: "random_number", Op: helper.Lte, Value: 50}},
		OrderBy:     helper.OrderByType{Field: "id", Op: helper.Desc},
		Pagination:  helper.Pagination{Size: 10, Page: 1},
		UserInfo: models.UserCredential{
			ProjectCode: "516",
		},
	})

	query.FilterProjectCode()
	errorQuery := query.Select()
	if errorQuery != nil {
		panic(errorQuery)
	}
	spew.Dump(another)
}
