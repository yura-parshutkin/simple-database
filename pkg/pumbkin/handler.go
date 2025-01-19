package pumbkin

import (
	"context"
	"fmt"
	"strconv"
)

type parser interface {
	Parse(command string) (*Query, error)
}

type storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) (bool, error)
}

type Handler struct {
	storage storage
	parser  parser
}

func NewHandler(
	storage storage,
	parser parser,
) *Handler {
	return &Handler{
		storage: storage,
		parser:  parser,
	}
}

func (h *Handler) Handle(ctx context.Context, command string) (string, error) {
	query, errPr := h.parser.Parse(command)
	if errPr != nil {
		return "", fmt.Errorf("parse query: %w", errPr)
	}
	switch query.Id {
	case DelID:
		ok, err := h.storage.Delete(ctx, query.Params[0])
		return strconv.FormatBool(ok), err
	case GetID:
		res, err := h.storage.Get(ctx, query.Params[0])
		return res, err
	case SetID:
		err := h.storage.Set(ctx, query.Params[0], query.Params[1])
		return "true", err
	default:
		return "", fmt.Errorf("unknown command: %s", query.Id)
	}
}
