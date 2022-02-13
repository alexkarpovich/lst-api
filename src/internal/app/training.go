package app

import (
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingType uint

const (
	TrainingDirect TrainingType = iota
	TrainingReverse
	TrainingListen
	TrainingCycles
)

type TrainingExpression struct {
	Id    *valueobject.ID `json:"id"`
	Value string          `json:"value"`
}

type TrainingAnswer struct {
	Id    *valueobject.ID `json:"id"`
	Value string          `json:"value"`
}

type TrainingMeta struct {
	StageCount      uint `json:"stageCount"`
	UniqueItemCount uint `json:"uniqueItemCount"`
	CompleteCount   uint `json:"completeCount"`
}

type TrainingItem struct {
	Id           *valueobject.ID     `json:"id" db:"id"`
	TrainingId   *valueobject.ID     `json:"trainingId" db:"training_id"`
	ExpressionId *valueobject.ID     `json:"expressionId" db:"expression_id"`
	Stage        uint                `json:"stage" db:"stage"`
	Cycle        uint                `json:"cycle" db:"cycle"`
	Complete     bool                `json:"complete" db:"complete"`
	Expression   *TrainingExpression `json:"expression"`
}

type Training struct {
	Id      *valueobject.ID  `json:"id" db:"id"`
	OwnerId *valueobject.ID  `json:"ownerId" db:"owner_id"`
	Type    TrainingType     `json:"type" db:"type"`
	Slices  []valueobject.ID `json:"slices" db:"slices"`
	Items   []*TrainingItem  `json:"-"`
	Meta    *TrainingMeta    `json:"meta" db:"meta"`
}

type TrainingRepo interface {
	Create(Training) (*Training, error)
	Reset(*valueobject.ID) error
	List(*valueobject.ID) ([]*Training, error)
	Get(*valueobject.ID) (*Training, error)
	GetByItemId(*valueobject.ID) (*Training, error)
	GetBySlices(Training) (*Training, error)
	NextItem(*valueobject.ID) (*TrainingItem, error)
	ItemAnswers(*valueobject.ID) ([]*TrainingAnswer, error)
	MarkItemAsComplete(*valueobject.ID) error
	HasCreatePermission(*valueobject.ID, []valueobject.ID) bool
}
