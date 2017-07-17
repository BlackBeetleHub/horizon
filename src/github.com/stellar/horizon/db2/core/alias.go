package core

import (
	sq "github.com/Masterminds/squirrel"
)

// AccountDataByKey loads a row from `accountdata`, by key
func (q *Q) AliasByKey(dest interface{}, addy string) error {
	sql := selectAlias.Limit(1).
		Where("accountid = ?", addy)
	return q.Get(dest, sql)
}

// TrustlinesByAddress loads all trustlines for `addy`
func (q *Q) AliasesByAddress(dest interface{}, addy string) error {
sql := selectAlias.Where("accountid = ?", addy)
return q.Select(dest, sql)
}

var selectAlias = sq.Select(
"al.accountid",
		 "al.aliasid",
).From("aliases al")
