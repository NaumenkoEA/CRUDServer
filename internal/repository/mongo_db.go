// Package repository : file contains operations with MongoDB
package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MRepository create connection with MongoDB
type MRepository struct {
	MPool *mongo.Client
}

// Create add new user to db
func (m *MRepository) Create(ctx context.Context, person *model.Person) (string, error) {

	newID := uuid.New().String()
	collection := m.MPool.Database("person").Collection("person")
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newID},
		{Key: "name", Value: person.Name},
		{Key: "password", Value: person.Password},
		{Key: "refreshtoken", Value: person.RefreshToken},
	})
	if err != nil {
		return "", fmt.Errorf("mongo: unable to create new user: %v", err)
	}
	return newID, nil
}

// Update update exist user
func (m *MRepository) Update(ctx context.Context, id string, person *model.Person) error {

	collection := m.MPool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: person.Name},
		{Key: "password", Value: person.Password},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}

// UpdateAuth add user refresh token
func (m *MRepository) UpdateAuth(ctx context.Context, id, refreshToken string) error {
	collection := m.MPool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "refreshtoken", Value: refreshToken},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}

// SelectAll take all users from db
func (m *MRepository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	var users []*model.Person
	collection := m.MPool.Database("person").Collection("person")
	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo: unable to select all users %v", err)
	}
	for c.Next(ctx) {
		user := model.Person{}
		err := c.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// Delete user from db
func (m *MRepository) Delete(ctx context.Context, id string) error {
	collection := m.MPool.Database("person").Collection("person")
	_, err := collection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		return fmt.Errorf("mongo: unable to delete user, %v", err)
	}
	return nil
}

// SelectByID select exist user from db by his id
func (m *MRepository) SelectByID(ctx context.Context, id string) (model.Person, error) {
	user := model.Person{}
	collection := m.MPool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// SelectByIDAuth take from user his refresh token
func (m *MRepository) SelectByIDAuth(ctx context.Context, id string) (model.Person, error) {
	user := model.Person{}
	collection := m.MPool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id},
		{Key: "name", Value: 0},
		{Key: "works", Value: 0},
		{Key: "age", Value: 0},
		{Key: "password", Value: 0},
	}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
