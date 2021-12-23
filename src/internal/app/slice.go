package app

import (
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

// FolderVisibility allow/deny users to use this folder
type SliceVisibility uint

const (
	// SlicePublic allows share the slice with other users
	SlicePublic SliceVisibility = iota
	// SlicePrivate deny share the slice
	SlicePrivate
)

type Expression struct {
	Id           *valueobject.ID `json:"id" db:"id"`
	Value        string          `json:"value" db:"value"`
	Translations []*Translation  `json:"translations"`
	CreatedAt    time.Time       `json:"createdAt" db:"created_at"`
}

type Translation struct {
	Id             *valueobject.ID `json:"id" db:"id"`
	Value          string          `json:"value"`
	Transcriptions []string        `json:"transcriptions" db:"transcriptions"`
	Comment        string          `json:"comment" db:"comment"`
	CreatedAt      time.Time       `json:"createdAt" db:"created_at"`
}

type TextTranslation struct {
	Id        *valueobject.ID `json:"id" db:"id"`
	Content   string          `json:"content" db:"content"`
	CreatedAt time.Time       `json:"createdAt" db:"createdAt"`
}

type Correction struct {
	Id        *valueobject.ID `json:"id" db:"id"`
	AuthorId  *valueobject.ID `json:"author_id" db:"author_id"`
	Title     string          `json:"title" db:"title"`
	Content   string          `json:"content" db:"content"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
}

type Text struct {
	Id          *valueobject.ID  `json:"id" db:"id"`
	AuthorId    *valueobject.ID  `json:"author_id" db:"author_id"`
	Title       string           `json:"title" db:"title"`
	Content     string           `json:"content" db:"content"`
	Translation *TextTranslation `json:"translation"`
	Corrections []*Correction    `json:"corrections"`
	CreatedAt   time.Time        `json:"createdAt" db:"created_at"`
}

type Slice struct {
	Id          *valueobject.ID `json:"id" db:"id"`
	TextId      *valueobject.ID `json:"textId" db:"text_id"`
	Name        string          `json:"name" db:"name"`
	Path        string          `json:"path" db:"path"`
	Visibility  SliceVisibility `json:"visibility" db:"visibility"`
	Text        *Text           `json:"text" db:"text"`
	Expressions []*Expression   `json:"expressions"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
}

type NestedSlice struct {
	Id         *valueobject.ID `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	Path       string          `json:"path" db:"path"`
	Count      uint            `json:"count" db:"count"`
	Visibility SliceVisibility `json:"visibility" db:"visibility"`
	Children   []*NestedSlice  `json:"children"`
}

type SliceRepo interface {
	Create(*valueobject.ID, *Slice) (*Slice, error)
	Get(*valueobject.ID) (*Slice, error)
	List(*valueobject.ID) ([]*NestedSlice, error)
	AttachExpression(*valueobject.ID, *Expression) (*Expression, error)
	DetachExpression(*valueobject.ID, *valueobject.ID) error
	AttachTranslation(*valueobject.ID, *valueobject.ID, *Translation) (*Translation, error)
	DetachTranslation(*valueobject.ID, *valueobject.ID) error
	AttachText(*valueobject.ID, *Text) (*Text, error)
	DetachText(*valueobject.ID, *valueobject.ID) error
}
