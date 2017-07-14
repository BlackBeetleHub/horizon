package core

import (
	sq "github.com/Masterminds/squirrel"
)


// TrustlinesByAddress loads all trustlines for `addy`
func (q *Q) AliasesByAddress(dest interface{}, addy string) error {
sql := selectAlias.Where("accountsourceid = ?", addy)
return q.Select(dest, sql)
}

var selectAlias = sq.Select(
"al.accountid",
"al.accountsourceid",
).From("aliases al")
