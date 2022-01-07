package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type LangRepo struct {
	db db.DB
}

func NewLangRepo(db db.DB) *LangRepo {
	return &LangRepo{db}
}

func (r *LangRepo) List() ([]*domain.Language, error) {
	var languages []*domain.Language
	query := `SELECT * FROM languages`

	err := r.db.Db().Select(&languages, query)
	if err != nil {
		return nil, err
	}

	return languages, nil
}
