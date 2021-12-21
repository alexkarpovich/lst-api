package repos

import (
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type ExpressionRepo struct {
	db db.DB
}

func NewExpressionRepo(db db.DB) *ExpressionRepo {
	return &ExpressionRepo{db}
}

func (r *ExpressionRepo) Create(obj *domain.Expression) (*domain.Expression, error) {
	stmt := `
		INSERT INTO expressions (lang_id, value) VALUES(:lang_id, :value)
		RETURNING id
	`
	rows, err := r.db.Db().NamedQuery(stmt, obj)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(&obj.Id)
	}

	if err = rows.Close(); err != nil {
		log.Println(err)
	}

	return obj, nil
}

func (r *ExpressionRepo) Get(obj *domain.Expression) (*domain.Expression, error) {
	var expression *domain.Expression
	query := `SELECT * FROM expressions WHERE lang_id=:lang_id AND value=:value`
	err := r.db.Db().Get(&expression, query, obj)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (r *ExpressionRepo) Search(langId *valueobject.ID, value string) ([]*domain.Expression, error) {
	var expressions []*domain.Expression
	query := `SELECT * FROM expressions WHERE lang_id=$1 AND value LIKE '%$2%'`

	err := r.db.Db().Select(&expressions, query, langId, value)
	if err != nil {
		return nil, err
	}

	return expressions, nil
}
