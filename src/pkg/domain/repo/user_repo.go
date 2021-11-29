package repo

import (
	"pkg/domain/entity"
	"pkg/domain/valueobject"
)

type UserRepo interface {
	Save(u *entity.User) (*entity.User, error)
	Get(userId valueobject.ID) (*entity.User, error)
	FindByEmail(emailAddress valueobject.EmailAddress) (*entity.User, error)
	Delete(userId valueobject.ID) error
}
