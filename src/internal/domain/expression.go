package domain

import "github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"

type Expression struct {
	Id     *valueobject.ID `json:"id" db:"id"`
	LandId *valueobject.ID `json:"langId" db:"lang_id"`
	Value  string          `json:"value" db:"value"`
	Lang   *Language       `json:"lang"`
}

type ExpressionRepo interface {
	Create(*Expression) (*Expression, error)
	Get(*Expression) (*Expression, error)
	Search(*valueobject.ID, string) ([]*Expression, error)
}
