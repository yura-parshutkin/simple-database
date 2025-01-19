package pumbkin

type QueryID string

const (
	GetID QueryID = "get"
	SetID QueryID = "set"
	DelID QueryID = "del"
)

type Query struct {
	Id     QueryID
	Params []string
}

func NewQuery(id QueryID, params []string) *Query {
	return &Query{Id: id, Params: params}
}
