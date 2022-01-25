package app

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingType uint

const (
	TrainingThrough TrainingType = iota
	TrainingReverse
	TrainingListen
	TrainingCycles
)

type TrainingExpression struct {
	Id    *valueobject.ID `json:"id"`
	Value string          `json:"value"`
}

type TrainingItem struct {
	Id           *valueobject.ID     `json:"id" db:"id"`
	TrainingId   *valueobject.ID     `json:"trainingId" db:"training_id"`
	ExpressionId *valueobject.ID     `json:"expressionId" db:"expression_id"`
	Stage        uint                `json:"stage" db:"stage"`
	Cycle        uint                `json:"cycle" db:"cycle"`
	Complete     bool                `json:"completed" db:"complete"`
	Expression   *TrainingExpression `json:"expression"`
}

type Training struct {
	Id      *valueobject.ID  `json:"id" db:"id"`
	OwnerId *valueobject.ID  `json:"ownerId" db:"owner_id"`
	Type    TrainingType     `json:"type" db:"type"`
	Slices  []valueobject.ID `json:"slices" db:"slices"`
	Items   []*TrainingItem  `json:"-"`
}

type TrainingRepo interface {
	Create(Training) (*Training, error)
	Get(*valueobject.ID) (*Training, error)
	GetByItemId(*valueobject.ID) (*Training, error)
	HasCreatePermission(*valueobject.ID, []valueobject.ID) bool
}
