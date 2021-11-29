package persistence

import (
	"database/sql"

	"pkg/domain/entity"
	"pkg/domain/valueobject"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) Save(u *entity.User) (*entity.User, error) {
	return nil, nil
}

func (ur *UserRepo) Get(userId valueobject.ID) (*entity.User, error) {
	return nil, nil
}

func (ur *UserRepo) FindByEmail(emailAddress valueobject.EmailAddress) (*entity.User, error) {
	return nil, nil
}

func (ur *UserRepo) Delete(userId valueobject.ID) error {
	return nil
}
