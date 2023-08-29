package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"item-service/internal/entity"

	"github.com/redis/go-redis/v9"
)

type itemCache struct {
	client *redis.Client
}

func NewItemCache(client *redis.Client) itemCache {
	return itemCache{client: client}
}

func marshalItem(item entity.Item) ([]byte, error) {
	return json.Marshal(item)
}

func unmarshalItem(data []byte) (entity.Item, error) {
	var item entity.Item
	err := json.Unmarshal(data, &item)
	if err != nil {
		return entity.Item{}, err
	}

	return item, nil
}

func (c itemCache) SetItem(ctx context.Context, newItem entity.Item) error {
	key := fmt.Sprintf("item-%d", newItem.ID)

	value, err := marshalItem(newItem)
	if err != nil {
		return err
	}

	err = c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c itemCache) GetItem(ctx context.Context, id int) (entity.Item, error) {
	key := fmt.Sprintf("item-%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return entity.Item{}, err
	}

	item, err := unmarshalItem([]byte(val))
	if err != nil {
		return entity.Item{}, err
	}

	return item, nil
}
