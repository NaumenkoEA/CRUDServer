// Package cache : file contains operations with cache
package cache

import (
	"awesomeProject/internal/model"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/gommon/log"
)

// UserCache struct for cache
type UserCache struct {
	redisClient *redis.Client
}

// NewCache create new redis connection
func NewCache(rdsClient *redis.Client) *UserCache {
	return &UserCache{redisClient: rdsClient}
}

// AddToCache add user to cache
func (u *UserCache) AddToCache(ctx context.Context, person *model.Person) error {
	user, err := json.Marshal(person)
	if err != nil {
		log.Errorf("cache: failed add user to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "user",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add user to cache, %e", err)
		return err
	}
	return nil
}
func (u *UserCache) AddAdvertToCache(ctx context.Context, advert *model.Advert) error {
	user, err := json.Marshal(advert)
	if err != nil {
		log.Errorf("cache: failed add advert to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "advert",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add advert to cache, %e", err)
		return err
	}
	return nil
}

// GetUserByIDFromCache get user from cache
func (u *UserCache) GetUserByIDFromCache(ctx context.Context) (model.Person, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"user", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Person{}, false, nil
		}
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Person{}, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	person := model.Person{}
	err = json.Unmarshal([]byte(msgString), &person)
	if err != nil {
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Person{}, false, err
	}
	return person, true, nil
}

func (u *UserCache) GetAdvertByIDFromCache(ctx context.Context) (model.Advert, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"advert", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Advert{}, false, nil
		}
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Advert{}, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	advert := model.Advert{}
	err = json.Unmarshal([]byte(msgString), &advert)
	if err != nil {
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Advert{}, false, err
	}
	return advert, true, nil
}

// DeleteUserFromCache delete user from cache
func (u *UserCache) DeleteUserFromCache(ctx context.Context) error {
	_, found, err := u.GetUserByIDFromCache(ctx)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	err = u.redisClient.FlushAll(ctx).Err()
	if err != nil {
		log.Errorf("failed to delete user from cache, %e", err)
		return err
	}
	return nil
}

func (u *UserCache) DeleteAdvertFromCache(ctx context.Context) error {
	_, found, err := u.GetAdvertByIDFromCache(ctx)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	err = u.redisClient.FlushAll(ctx).Err()
	if err != nil {
		log.Errorf("failed to delete user from cache, %e", err)
		return err
	}
	return nil
}

// GetAllUsersFromCache get all users from cache
func (u *UserCache) GetAllUsersFromCache(ctx context.Context) ([]*model.Person, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"all-users", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		log.Errorf("failed get user by id from cache: %e", err)
		return nil, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	var persons []*model.Person

	err = json.Unmarshal([]byte(msgString), &persons)
	if err != nil {
		log.Errorf("failed to unmarshal json, %e", err)
		return nil, false, err
	}

	return persons, true, nil
}

func (u *UserCache) GetAllAdvertsFromCache(ctx context.Context) ([]*model.Advert, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"all-adverts", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		log.Errorf("failed get advers by id from cache: %e", err)
		return nil, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	var adverts []*model.Advert

	err = json.Unmarshal([]byte(msgString), &adverts)
	if err != nil {
		log.Errorf("failed to unmarshal json, %e", err)
		return nil, false, err
	}

	return adverts, true, nil
}

// AddAllUsersToCache add all users from db to cache
func (u *UserCache) AddAllUsersToCache(ctx context.Context, person []*model.Person) error {
	user, err := json.Marshal(person)
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "all-users",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	return nil
}

func (u *UserCache) AddAllAdvertsToCache(ctx context.Context, adverts []*model.Advert) error {
	advert, err := json.Marshal(adverts)
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "all-users",
		ID:     "0-*",
		Values: map[string]interface{}{"About": advert},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	return nil
}
