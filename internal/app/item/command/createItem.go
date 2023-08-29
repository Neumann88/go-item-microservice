package command

import (
	"context"
	"item-service/internal/entity"
)

type repository interface {
	CreateItem(ctx context.Context, newItem entity.Item) error
}

type cache interface {
	SetItem(ctx context.Context, newItem entity.Item) error
}

type CreateItemHandler struct {
	repository repository
	cache      cache
}

func NewCreateItemHandler(repo repository, cache cache) CreateItemHandler {
	return CreateItemHandler{repository: repo, cache: cache}
}

func (h CreateItemHandler) Handle(ctx context.Context, newItem entity.Item) error {
	err := h.repository.CreateItem(ctx, newItem)
	if err != nil {
		return err
	}

	_ = h.cache.SetItem(ctx, newItem)
	return nil
}
