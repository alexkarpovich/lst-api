package app

import (
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

// NodeVisibility allow/deny users to use this node
type NodeVisibility uint

const (
	// NodePublic allows share the node with other users
	NodePublic NodeVisibility = iota
	// NodePrivate deny share the node
	NodePrivate
)

type NodeType uint

const (
	NodeFolder NodeType = iota
	NodeSlice
)

type Expression struct {
	Id             *valueobject.ID  `json:"id" db:"id"`
	Value          string           `json:"value" db:"value"`
	Transcriptions []*Transcription `json:"transcriptions"`
	Translations   []*Translation   `json:"translations"`
	CreatedAt      time.Time        `json:"createdAt" db:"created_at"`
}

type Translation struct {
	Id             *valueobject.ID  `json:"id" db:"id"`
	Value          string           `json:"value"`
	Transcriptions []*Transcription `json:"transcriptions" db:"transcriptions"`
	Comment        string           `json:"comment" db:"comment"`
	CreatedAt      time.Time        `json:"createdAt" db:"created_at"`
}

type Transcription struct {
	Id    *valueobject.ID `json:"id" db:"id"`
	Value string          `json:"value" db:"value"`
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

type Node struct {
	Id          *valueobject.ID `json:"id" db:"id"`
	TextId      *valueobject.ID `json:"textId" db:"text_id"`
	Type        NodeType        `json:"type" db:"type"`
	Name        string          `json:"name" db:"name"`
	Path        string          `json:"path" db:"path"`
	Visibility  NodeVisibility  `json:"visibility" db:"visibility"`
	Text        *Text           `json:"text" db:"text"`
	Expressions []*Expression   `json:"expressions"`
	CreatedAt   time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" db:"updated_at"`
}

type NodeView struct {
	Expressions []*Expression `json:"expressions"`
}

type FlatNode struct {
	Id         *valueobject.ID `json:"id" db:"id"`
	Type       NodeType        `json:"type" db:"type"`
	Name       string          `json:"name" db:"name"`
	Path       string          `json:"path" db:"path"`
	Count      uint            `json:"count" db:"count"`
	Visibility NodeVisibility  `json:"visibility" db:"visibility"`
}

type NodeRepo interface {
	Create(*valueobject.ID, Node) (*Node, error)
	Get(*valueobject.ID) (*Node, error)
	View([]valueobject.ID) (*NodeView, error)
	List(*valueobject.ID) ([]*FlatNode, error)
	FilterSliceIds([]valueobject.ID) ([]valueobject.ID, error)
	Update(FlatNode) error
	AttachExpression(*valueobject.ID, Expression) (*Expression, error)
	DetachExpression(*valueobject.ID, *valueobject.ID) error
	TranslationsBySlices([]valueobject.ID) ([]*Translation, error)
	AvailableTranslations(*valueobject.ID, *valueobject.ID) ([]*Translation, error)
	AttachTranslation(*valueobject.ID, *valueobject.ID, Translation) (*Translation, error)
	DetachTranslation(*valueobject.ID, *valueobject.ID) error
	AttachText(*valueobject.ID, Text) (*Text, error)
	DetachText(*valueobject.ID) error
}
