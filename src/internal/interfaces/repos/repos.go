package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type Repos struct {
	User       app.UserRepo
	Group      app.GroupRepo
	Node       app.NodeRepo
	Expression domain.ExpressionRepo
	Lang       domain.LangRepo
	Training   app.TrainingRepo
}

func NewRepos(db db.DB) *Repos {
	return &Repos{
		User:       NewUserRepo(db),
		Group:      NewGroupRepo(db),
		Node:       NewNodeRepo(db),
		Expression: NewExpressionRepo(db),
		Lang:       NewLangRepo(db),
		Training:   NewTrainingRepo(db),
	}
}
