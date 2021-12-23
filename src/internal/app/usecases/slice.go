package usecases

import (
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type SliceInteractor struct {
	SliceRepo      app.SliceRepo
	ExpressionRepo domain.ExpressionRepo
}

func NewSliceInteractor(pr app.SliceRepo, er domain.ExpressionRepo) *SliceInteractor {
	return &SliceInteractor{pr, er}
}

func (i *SliceInteractor) Create(groupId *valueobject.ID, s *app.Slice) (*app.Slice, error) {
	slice, err := i.SliceRepo.Create(groupId, s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return slice, nil
}

func (i *SliceInteractor) Get(sliceId *valueobject.ID) (*app.Slice, error) {
	slice, err := i.SliceRepo.Get(sliceId)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (i *SliceInteractor) AttachExpression(sliceId *valueobject.ID, expression *app.Expression) (*app.Expression, error) {
	expression, err := i.SliceRepo.AttachExpression(sliceId, expression)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (i *SliceInteractor) DetachExpression(sliceId *valueobject.ID, expressionId *valueobject.ID) error {
	err := i.SliceRepo.DetachExpression(sliceId, expressionId)
	if err != nil {
		return err
	}

	return nil
}

func (i *SliceInteractor) AttachTranslation(sliceId *valueobject.ID, expressionId *valueobject.ID, translation *app.Translation) (*app.Translation, error) {
	translation, err := i.SliceRepo.AttachTranslation(sliceId, expressionId, translation)
	if err != nil {
		return nil, err
	}

	return translation, nil
}

func (i *SliceInteractor) DetachTranslation(sliceId *valueobject.ID, translationId *valueobject.ID) error {
	err := i.SliceRepo.DetachTranslation(sliceId, translationId)
	if err != nil {
		return err
	}

	return nil
}

func (i *SliceInteractor) AttachText(sliceId *valueobject.ID, text *app.Text) (*app.Text, error) {
	text, err := i.SliceRepo.AttachText(sliceId, text)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func (i *SliceInteractor) DetachText(sliceId *valueobject.ID, textId *valueobject.ID) error {
	err := i.SliceRepo.DetachText(sliceId, textId)
	if err != nil {
		return err
	}

	return nil
}
