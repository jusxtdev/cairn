package user

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckUserExists(db *pgxpool.Pool, r *http.Request, id string) (bool, error) {
	// MISTAKE :- used gotId as int, which resulted in 0 when doing Scan
	// FIX :- know what data types are coming in and out of a process you are calling and use the variables accordingly
	var gotId string
	err := db.QueryRow(
		r.Context(),
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