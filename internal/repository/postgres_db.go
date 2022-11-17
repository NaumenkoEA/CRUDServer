// Package repository : file contains operations with PostgresDB
package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

// PRepository :creating new connection with PostgresDB
type PRepository struct {
	PPool *pgxpool.Pool
}

// Create : insert new user into database
func (r *PRepository) Create(ctx context.Context, person *model.Person) (string, error) {
	newID := uuid.New().String()
	_, err := r.PPool.Exec(ctx, "insert into persons(id,name,password) values($1,$2,$3)",
		newID, &person.Name, &person.Password)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}
	return newID, nil
}

// SelectAll : Print all users(ID,Name,Works) from database
func (r *PRepository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person
	rows, err := r.PPool.Query(ctx, "select id,name from persons")
	if err != nil {
		log.Errorf("database error with select all users, %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := model.Person{}
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			log.Errorf("database error with select all users, %v", err)
			return nil, err
		}
		persons = append(persons, &p)
	}

	return persons, nil
}

// Delete : delete user by his ID
func (r *PRepository) Delete(ctx context.Context, id string) error {
	a, err := r.PPool.Exec(ctx, "delete from persons where id=$1", id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("error with delete user %v", err)
		return err
	}
	return nil
}

// UpdateAuth : update user refreshToken by his ID
func (r *PRepository) UpdateAuth(ctx context.Context, id, refreshToken string) error {
	a, err := r.PPool.Exec(ctx, "update persons set refreshToken=$1 where id=$2", refreshToken, id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update user %v", err)
		return err
	}
	return nil
}

// Update update user in db
func (r *PRepository) Update(ctx context.Context, id string, p *model.Person) error {
	a, err := r.PPool.Exec(ctx, "update persons set name=$1 where id=$2", &p.Name, id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update user %v", err)
		return err
	}
	return nil
}

// SelectByID : select one user by his ID
func (r *PRepository) SelectByID(ctx context.Context, id string) (model.Person, error) {
	p := model.Person{}
	err := r.PPool.QueryRow(ctx, "select id,name,password from persons where id=$1", id).Scan(
		&p.ID, &p.Name, &p.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Person{}, fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("database error, select by id: %v", err)
		return model.Person{}, err /*p, fmt.errorf("user with this id doesn't exist")*/
	}
	return p, nil
}

// SelectByIDAuth select auth user
func (r *PRepository) SelectByIDAuth(ctx context.Context, id string) (model.Person, error) {
	p := model.Person{}
	err := r.PPool.QueryRow(ctx, "select id,refreshToken from persons where id=$1", id).Scan(&p.ID, &p.RefreshToken)

	if err != nil /*err==no-records*/ {
		if err == pgx.ErrNoRows {
			return model.Person{}, fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("database error, select by id: %v", err)
		return model.Person{}, err /*p, fmt.errorf("user with this id doesn't exist")*/
	}
	return p, nil
}

func (r *PRepository) CreateAdvert(ctx context.Context, advert *model.Advert) (string, error) {
	newID := uuid.New().String()
	_, err := r.PPool.Exec(ctx, "insert into advert(id,address,price) values($1,$2,$3)",
		newID, &advert.Address, &advert.Price)
	if err != nil {
		log.Errorf("database error with create advert: %v", err)
		return "", err
	}
	return newID, nil
}

func (r *PRepository) SelectAllAdvert(ctx context.Context) ([]*model.Advert, error) {
	var adverts []*model.Advert
	rows, err := r.PPool.Query(ctx, "select id,address,price from adverts")
	if err != nil {
		log.Errorf("database error with select all adverts, %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		advert := model.Advert{}
		err := rows.Scan(&advert.ID, &advert.Address, &advert.Price)
		if err != nil {
			log.Errorf("database error with select all adverts, %v", err)
			return nil, err
		}
		adverts = append(adverts, &advert)
	}

	return adverts, nil
}

func (r *PRepository) DeleteAdvert(ctx context.Context, id string) error {
	a, err := r.PPool.Exec(ctx, "delete from adverts where id=$1", id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("advert with this id doesnt exist")
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("advert with this id doesnt exist: %v", err)
		}
		log.Errorf("error with delete advert %v", err)
		return err
	}
	return nil
}

func (r *PRepository) UpdateAdvert(ctx context.Context, id string, advert *model.Advert) error {
	a, err := r.PPool.Exec(ctx, "update adverts set address=$1,price=$2 where id=$4", &advert.Address, &advert.Price, id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update user %v", err)
		return err
	}
	return nil
}
