package service

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

type Handler struct { // handler
	s *Service
}

var (
	Pool *pgxpool.Pool
)

// NewHandler :define new handlers
func NewHandler(newS *Service) *Handler {
	return &Handler{s: newS}
}

func TestMain(m *testing.M) {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatalf("Bad connection: %v", err)
	}
	Pool = pool
	run := m.Run()
	os.Exit(run)
}

func TestService_Authentication(t *testing.T) {
	rps := NewService(&repository.PRepository{Pool: Pool}, &cache.UserCache{})
	h := NewHandler(rps)
	_, _, err := h.s.Authentication(context.Background(), "a20fc586-d9d2-4969-909f-d00bf42aa88a", "tujh2004")
	require.NoError(t, err, "passwords dont match")
	_, _, err = h.s.Authentication(context.Background(), "a20fc586-d9d2-4969-909f-d00bf42aa88a", "tujh2005")
	require.Error(t, err, "passwords dont match")
	_, _, err = h.s.Authentication(context.Background(), "a20fc586-d9d2", "tujh2005")
	require.Error(t, err, "passwords dont match or this user doesnt exist")
}
func TestHashPassword(t *testing.T) {
	testNoValidData := []string{"", " "}
	_, err := HashPassword("worker")
	require.NoError(t, err, "cannot register this user")
	for _, user := range testNoValidData {
		_, err := HashPassword(user)
		require.Error(t, err, "this user can be register")
	}
}

func TestService_Registration(t *testing.T) {
	rps := NewService(&repository.PRepository{Pool: Pool}, &cache.UserCache{})
	h := NewHandler(rps)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testValidData := model.Person{Name: "Mani", Works: true, Age: 19, Password: "worker"}
	testNoValidData := model.Person{Name: "Anton", Works: true, Age: 19, Password: "worker"}
	_, err := h.s.Registration(ctx, &testValidData)
	require.NoError(t, err, "cannot register this user")
	_, err = h.s.Registration(ctx, &testNoValidData)
	require.Error(t, err, "this user can be register")
}

func TestService_RefreshToken(t *testing.T) {
	rps := NewService(&repository.PRepository{Pool: Pool}, &cache.UserCache{})
	h := NewHandler(rps)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, _, err := h.s.RefreshToken(ctx, "token")
	require.NoError(t, err, "cannot refresh your tokens")
	_, _, err = h.s.RefreshToken(ctx, "<false token>")
	require.Error(t, err, "can refresh your tokens")
	_, _, err = h.s.RefreshToken(ctx, "old token")
	require.Error(t, err, "token already valid")
}

func TestCreateJWT(t *testing.T) {
	testUser := model.Person{
		ID:       "a20fc586-d9d2-4969-909f-d00bf42aa88a",
		Name:     "Egor Tihonov",
		Works:    true,
		Age:      18,
		Password: "tujh2004",
	}
	testUserNoValidate := model.Person{
		ID:       "0",
		Name:     "Egor Tihonov",
		Works:    true,
		Age:      120,
		Password: "tujh2004",
	}
	s := NewService(&repository.PRepository{Pool: Pool}, &cache.UserCache{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, _, err := s.CreateJWT(ctx, s.rps, &testUser)
	require.NoError(t, err, "cannot create tokens")
	_, _, err = s.CreateJWT(ctx, s.rps, &testUserNoValidate)
	require.Error(t, err, "tokens create")
}
