package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

func (s *UserRepository) CheckUserExistsById(ctx context.Context, id string) (bool, error) {
	// MISTAKE :- used gotId as int, which resulted in 0 when doing Scan
	// FIX :- know what data types are coming in and out of a process you are calling and use the variables accordingly
	var gotId string
	err := s.Db.QueryRow(
		ctx,
		"SELECT id FROM users WHERE id = $1",
		id,
	).Scan(&gotId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// function that finds user with the provided id
// throws error only if something happens from db.QueryRow
// let caller handle userNotFound
func (s *UserRepository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.Db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT email FROM users WHERE email = $1)",
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type CreateUserInput struct {
	Name           string
	Email          string
	HashedPassword string
}

func (s *UserRepository) AddUser(ctx context.Context, u CreateUserInput) (string, error) {
	query := `
	INSERT INTO users (name, email, password)
	VALUES ($1, $2, $3) 
	RETURNING id
	`
	var id string
	err := s.Db.QueryRow(ctx, query, u.Name, u.Email, u.HashedPassword).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *UserRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.Db.QueryRow(
		ctx,
		"SELECT * FROM users WHERE email=$1",
		email,
	).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (s *UserRepository) GetUserById(ctx context.Context, id string) (User, error) {
	var user User
	err := s.Db.QueryRow(
		ctx,
		"SELECT * FROM users WHERE id=$1",
		id,
	).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}