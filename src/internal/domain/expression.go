package domain

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type Expression struct {
	Id       *valueobject.ID `json:"id" db:"id"`
	LangCode string          `json:"langCode" db:"lang"`
	Value    string          `json:"value" db:"value"`
	Lang     *Language       `json:"lang"`
}

type ExpressionRepo interface {
	Create(*Expression) (*Expression, error)
	Get(*Expression) (*Expression, error)
	Search(string, string) ([]*Expression, error)
}
