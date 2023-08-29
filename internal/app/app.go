package app

import (
	itemCommand "item-service/internal/app/item/command"
	itemQuery "item-service/internal/app/item/query"

	"item-service/internal/adapters/cache"
	"item-service/internal/adapters/repository"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Command struct {
	CreateItem itemCommand.CreateItemHandler
}

type Query struct {
	GetItemByID itemQuery.GetItemByIDHandler
}

type Application struct {
	Command Command
	Query   Query
}

func New(db *sqlx.DB, cacheClient *redis.Client) Application {
	itemRepo := repository.NewItemRepository(db)
	itemCache := cache.NewItemCache(cacheClient)

	createItem := itemCommand.NewCreateItemHandler(itemRepo, itemCache)
	getItemByID := itemQuery.NewGetItemByIDHandler(itemRepo, itemCache)
	return Application{
		Command: Command{
			CreateItem: createItem,
		},
		Query: Query{
			GetItemByID: getItemByID,
		},
	}
}
