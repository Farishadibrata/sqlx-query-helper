package helper

import (
	"fmt"

	models "farishadibrata.com/sqlxhelper/model"
	"github.com/jmoiron/sqlx"
)

type OperatorEnum string

const (
	Eq   OperatorEnum = "="
	Ne   OperatorEnum = "!="
	Gt   OperatorEnum = ">"
	Lt   OperatorEnum = "<"
	Gte  OperatorEnum = ">="
	Lte  OperatorEnum = "<="
	Like OperatorEnum = "Like"
)

type OrderEnum string

const (
	Asc  OrderEnum = "ASC"
	Desc OrderEnum = "DESC"
)

type WhereType struct {
	Field string
	Op    OperatorEnum
	Value interface{}
}

type OrderByType struct {
	Field string
	Op    OrderEnum
}

type Pagination struct {
	Size int
	Page int
}

type QueryHelperParams struct {
	DB          *sqlx.DB
	Destination interface{}
	Query       string
	UserInfo    models.UserCredential
	Where       []WhereType
	OrderBy     OrderByType
	Pagination  Pagination
}

type QueryHelper struct {
	p                 QueryHelperParams
	CompiledQuery     string
	CompiledArguments []interface{}
	OrderBy           string
}

func NewQuery(qh QueryHelperParams) QueryHelper {
	return QueryHelper{
		p:             qh,
		CompiledQuery: qh.Query,
	}
}

func (q *QueryHelper) addOrderBy() {
	if q.p.OrderBy.Field != "" {
		q.CompiledQuery = fmt.Sprintf("%s ORDER BY %s %s", q.CompiledQuery, q.p.OrderBy.Field, string(q.p.OrderBy.Op))
	}
}

func (q *QueryHelper) addWhere() {
	params := q.p
	if len(params.Where) == 0 {
		return
	}
	q.CompiledQuery = q.CompiledQuery + " WHERE "
	for i, condition := range params.Where {
		q.CompiledArguments = append(q.CompiledArguments, condition.Value)
		q.CompiledQuery = fmt.Sprintf("%s %s %s $%d", q.CompiledQuery, condition.Field, string(condition.Op), i+1)
		if i != len(params.Where)-1 {
			q.CompiledQuery = q.CompiledQuery + " AND "
		}
	}
}

func (q *QueryHelper) addPagination() {
	q.CompiledQuery = fmt.Sprintf("%s OFFSET %d LIMIT %d", q.CompiledQuery, (q.p.Pagination.Page-1)*q.p.Pagination.Size, q.p.Pagination.Size)
}

func (q *QueryHelper) CompileQuery() (string, error) {
	q.addWhere()
	q.addOrderBy()
	q.addPagination()
	return q.CompiledQuery, nil
}

func (q *QueryHelper) FilterProjectCode() {
	params := &q.p
	params.Where = append(params.Where, WhereType{
		Field: "project_code",
		Op:    Eq,
		Value: params.UserInfo.ProjectCode,
	})
}

func (q *QueryHelper) Select() error {
	params := q.p
	compiledQuery, errorCompile := q.CompileQuery()

	if errorCompile != nil {
		return errorCompile
	}

	err := params.DB.Select(params.Destination, compiledQuery, q.CompiledArguments...)
	if err != nil {
		return err
	}
	return nil
}
