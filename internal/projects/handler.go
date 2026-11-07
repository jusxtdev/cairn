package projects

import "github.com/jackc/pgx/v5/pgxpool"

type ProjectsHandler struct {
	ProjectRepo *ProjectRepository
}

func NewHandler(db *pgxpool.Pool) *ProjectsHandler {
	return &ProjectsHandler{
		&ProjectRepository{
			Db: db,
		},
	}
}