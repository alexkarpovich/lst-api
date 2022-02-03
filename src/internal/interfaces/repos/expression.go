package repos

import (
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
	"github.com/lib/pq"
)

type ExpressionRepo struct {
	db db.DB
}

func NewExpressionRepo(db db.DB) *ExpressionRepo {
	return &ExpressionRepo{db}
}

func (r *ExpressionRepo) Create(obj *domain.Expression) (*domain.Expression, error) {
	stmt := `
		INSERT INTO expressions (lang, value) VALUES(:lang, :value)
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

func (r *ExpressionRepo) Get(id *valueobject.ID) (*domain.Expression, error) {
	expression := &domain.Expression{}
	query := `SELECT * FROM expressions WHERE id=$1`
	err := r.db.Db().Get(expression, query, id)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (r *ExpressionRepo) Search(langCode string, value string) ([]*domain.Expression, error) {
	var expressions []*domain.Expression
	query := `SELECT * FROM expressions WHERE lang=$1 AND value LIKE $2`

	err := r.db.Db().Select(&expressions, query, langCode, "%"+value+"%")
	if err != nil {
		return nil, err
	}

	return expressions, nil
}

func (r *ExpressionRepo) GetTranscriptionMap(typeId *valueobject.ID, exprParts []string) (map[string][]*domain.TranscriptionItem, error) {
	query := `
		SELECT e.value, t.id, t.value FROM expressions e
		LEFT JOIN expression_transcription et ON et.expression_id=e.id
		LEFT JOIN transcriptions t ON et.transcription_id=t.id
		WHERE t.type=$1 and e.value = ANY($2) 
	`
	rows, err := r.db.Db().Query(query, typeId, pq.StringArray(exprParts))
	if err != nil {
		return nil, err
	}

	transcriptionMap := make(map[string][]*domain.TranscriptionItem)

	for rows.Next() {
		var exprValue string
		item := &domain.TranscriptionItem{}
		rows.Scan(&exprValue, &item.Id, &item.Value)
		transcriptionMap[exprValue] = append(transcriptionMap[exprValue], item)
	}

	return transcriptionMap, nil
}
