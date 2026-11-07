package projects

import "github.com/jackc/pgx/v5/pgxpool"

type ProjectRepository struct {
	Db *pgxpool.Pool
}