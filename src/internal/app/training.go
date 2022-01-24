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
	Id            *valueobject.ID     `json:"id" db:"id"`
	TrainingId    *valueobject.ID     `json:"trainingId" db:"training_id"`
	TranslationId *valueobject.ID     `json:"translationId" db:"translation_id"`
	Stage         *uint               `json:"stage,omitempty" db:"stage"`
	Cycle         *uint               `json:"cycle,omitempty" db:"cycle"`
	Complete      bool                `json:"completed" db:"complete"`
	Expression    *TrainingExpression `json:"expression"`
}

type Training struct {
	Id      *valueobject.ID  `json:"id" db:"id"`
	OwnerId *valueobject.ID  `json:"ownerId" db:"owner_id"`
	Type    TrainingType     `json:"type" db:"type"`
	Nodes   []valueobject.ID `json:"nodes" db:"nodes"`
}

type TrainingRepo interface {
	Create(Training) (*Training, error)
	Get(*valueobject.ID) (*Training, error)
	GetByItemId(*valueobject.ID) (*Training, error)
	HasCreatePermission(*valueobject.ID, []valueobject.ID) bool
}

// const stageCount = Math.round(Math.log(count * 1. / MIN_CHUNK_SIZE) / Math.log(2)) + 1;
// let stages = Array.from(Array(stageCount).keys());

// let ids, chunkSize, rate;

// stages = stages.map((k) => {
// 	ids = shuffle(transIds);
// 	rate = Math.round(count / (MIN_CHUNK_SIZE * Math.pow(2, k)));
// 	chunkSize = Math.round(count / rate);

// 	return chunk(ids, chunkSize);
// });

// console.log(count, stages.map(cycles => cycles.length));
