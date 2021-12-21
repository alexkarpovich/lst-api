package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type Repos struct {
	User       app.UserRepo
	Group      app.GroupRepo
	Slice      app.SliceRepo
	Expression domain.ExpressionRepo
}

func NewRepos(db db.DB) *Repos {
	return &Repos{
		User:       NewUserRepo(db),
		Group:      NewGroupRepo(db),
		Slice:      NewSliceRepo(db),
		Expression: NewExpressionRepo(db),
	}
}
