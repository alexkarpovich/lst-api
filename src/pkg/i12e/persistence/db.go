package persistence

import (
	"database/sql"

	"pkg/domain/repo"
)

type Repos struct {
	User *repo.UserRepo
}

func NewRepos() (*Repos, error) {
	db, err := sql.Open("postgres", "")
	if err != nil {
		return nil, err
	}

	return &Repos{
		User: NewUserRepo(db),
	}, nil
}
