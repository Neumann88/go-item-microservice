package query

import (
	"context"
	"item-service/internal/entity"
)

type repository interface {
	GetItemByID(ctx context.Context, id int) (entity.Item, error)
}

type cache interface {
	GetItem(ctx context.Context, id int) (entity.Item, error)
}

type GetItemByIDHandler struct {
	repository repository
	cache      cache
}

func NewGetItemByIDHandler(repo repository, cache cache) GetItemByIDHandler {
	return GetItemByIDHandler{repository: repo, cache: cache}
}

func (h GetItemByIDHandler) Handle(ctx context.Context, id int) (entity.Item, error) {
	item, err := h.cache.GetItem(ctx, id)
	if err == nil {
		return item, nil
	}

	item, err = h.repository.GetItemByID(ctx, id)
	if err != nil {
		return entity.Item{}, err
	}

	return item, nil
}
