package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	Db *pgxpool.Pool
}

var ErrUserNotFound = errors.New("user not found")

func (s *UserService) CheckUserExistsById(ctx context.Context, id string) (bool, error) {
	// MISTAKE :- used gotId as int, which resulted in 0 when doing Scan
	// FIX :- know what data types are coming in and out of a process you are calling and use the variables accordingly
	var gotId string
	err := s.Db.QueryRow(
		ctx,
		"SELECT id FROM users WHERE id = $1",
		id,
	).Scan(&gotId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// function that finds user with the provided id
// throws error only if something happens from db.QueryRow
// let caller handle userNotFound
func (s *UserService) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
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
	Id       uuid.UUID
	Name     string
	Email    string
	HashedPassword string
}
func (s *UserService) AddUser(ctx context.Context, u CreateUserInput) error {
		query := `
	INSERT INTO users (id, name, email, password)
	VALUES ($1, $2, $3, $4) 
	`
	_, err := s.Db.Exec(ctx, query, u.Id, u.Name, u.Email, u.HashedPassword)
	return  err
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := s.Db.QueryRow(
		ctx,
		"SELECT * FROM users WHERE email=$1",
		email,
	).Scan(&user)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return User{}, ErrUserNotFound
		}
		return User{}, err
	} 
	return user, nil
}