// Package repository : file contains operations with all DBs
package repository

import (
	"awesomeProject/internal/model"
	"context"
)

// Repository middleware
type Repository interface {
	Create(ctx context.Context, person *model.Person) (string, error)
	CreateAdvert(ctx context.Context, advert *model.Advert) (string, error)

	UpdateAuth(ctx context.Context, id string, refreshToken string) error
	Update(ctx context.Context, id string, person *model.Person) error
	UpdateAdvert(ctx context.Context, id string, advert *model.Advert) error

	SelectAll(ctx context.Context) ([]*model.Person, error)
	SelectAllAdvert(ctx context.Context) ([]*model.Advert, error)
	SelectByID(ctx context.Context, id string) (model.Person, error)
	SelectAdvertByID(ctx context.Context, id string) (model.Advert, error)
	SelectByIDAuth(ctx context.Context, id string) (model.Person, error)

	Delete(ctx context.Context, id string) error
	DeleteAdvert(ctx context.Context, id string) error
}
