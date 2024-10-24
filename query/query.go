package query

type Filter interface {
	Apply(*Query)
}

type Query struct {
	Filters   map[string]interface{}
	SortField string
	SortOrder string
}

func NewQuery() *Query {
	return &Query{
		Filters: make(map[string]interface{}),
	}
}

type SortOption struct {
	Field string
	Order string // "asc" or "desc"
}

const (
	ASC  = "asc"
	DESC = "desc"
)
