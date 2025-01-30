package compute

import (
	"errors"
	"fmt"
	"simpledatabase/pkg/pumbkin/dto"
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Parse(command string) (*dto.Query, error) {
	queryParts := strings.Split(command, " ")
	if len(queryParts) < 2 {
		return nil, errors.New("invalid command format")
	}
	params := queryParts[1:]
	id := dto.QueryID(strings.ToLower(queryParts[0]))
	switch id {
	case dto.DelID:
		if len(params) != 1 {
			return nil, errors.New("invalid command format, expect 1 parameter for del")
		}
	case dto.GetID:
		if len(params) != 1 {
			return nil, errors.New("invalid command format, expect 1 parameter for get")
		}
	case dto.SetID:
		if len(params) != 2 {
			return nil, errors.New("invalid command format, expect 2 parameters for set")
		}
	default:
		return nil, fmt.Errorf("unknown command: %s", id)
	}
	return dto.NewQuery(id, params), nil
}
