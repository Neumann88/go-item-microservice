package repository

import (
	"context"

	"item-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

type itemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) itemRepository {
	return itemRepository{db}
}

const createItemQuery = `INSERT INTO item (id, value) VALUES (:id, :value)`

func (r itemRepository) CreateItem(ctx context.Context, newItem entity.Item) error {
	_, err := r.db.NamedExecContext(ctx, createItemQuery, newItem)
	if err != nil {
		return err
	}

	return nil
}

const getItemByIDquery = `SELECT id, value FROM item WHERE id = $1`

func (r itemRepository) GetItemByID(ctx context.Context, id int) (entity.Item, error) {
	var item entity.Item
	err := r.db.GetContext(ctx, &item, getItemByIDquery, id)
	if err != nil {
		return entity.Item{}, err
	}

	return item, nil
}
