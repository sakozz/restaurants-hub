package database

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"golang.org/x/exp/constraints"
)

type DataValue interface {
	string | constraints.Float | constraints.Integer | constraints.Complex
}

type SqlBuilder interface {
	Filter(tableName string, params url.Values) string
	Insert(tableName string, data interface{}) string
	Update(tableName string, id *int64, data interface{}) string
	Find(tableName string, params map[string]interface{}) string
	SearchBy(tableName string, params map[string]interface{}) string
}

type sqlBuilder struct {
}

func NewSqlBuilder() SqlBuilder {
	return &sqlBuilder{}
}

func (builder *sqlBuilder) Find(tableName string, params map[string]interface{}) string {
	exp := filtersToSql(params)

	sql, _, _ := goqu.From(tableName).Where(exp).ToSQL()
	return sql
}

func (builder *sqlBuilder) SearchBy(tableName string, params map[string]interface{}) string {
	exp := filtersToSql(params)
	sql, _, _ := goqu.From(tableName).Where(exp).ToSQL()
	return sql
}

func (builder *sqlBuilder) Filter(tableName string, params url.Values) string {
	query := map[string]interface{}{}
	for key, value := range params {
		query[key] = value
	}

	exp := filtersToSql(query)
	sql, _, _ := goqu.From(tableName).Where(exp).ToSQL()
	return sql
}

func (builder *sqlBuilder) Insert(tableName string, data interface{}) string {
	ds := goqu.Insert(tableName).Rows(data).Returning(goqu.T(tableName).All())

	insertSQL, _, _ := ds.ToSQL()
	return insertSQL
}

func (builder *sqlBuilder) Update(tableName string, id *int64, data interface{}) string {
	fmt.Println(data)
	ds := goqu.Update(tableName).Set(data).Where(goqu.Ex{
		"id": id,
	}).Returning(goqu.T(tableName).All())

	insertSQL, _, _ := ds.ToSQL()
	return insertSQL
}

func filtersToSql(params map[string]interface{}) exp.Ex {
	exp := goqu.Ex{}
	optMapping := map[string]string{
		"prefix":   "like",
		"contains": "contains",
		"in":       "in",
		"gt":       "gt",
		"gte":      "gte",
	}

	for key, value := range params {
		splits := strings.Split(key, "__")
		attr := splits[0]
		if len(splits) > 1 {
			opt := splits[1]
			exp[attr] = goqu.Op{optMapping[opt]: value}
		} else {
			exp[attr] = value
		}
	}
	return exp
}
