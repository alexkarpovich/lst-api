package usecases

import (
	"errors"
	"log"
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type NodeInteractor struct {
	NodeRepo       app.NodeRepo
	GroupRepo      app.GroupRepo
	ExpressionRepo domain.ExpressionRepo
}

func NewNodeInteractor(pr app.NodeRepo, gr app.GroupRepo, er domain.ExpressionRepo) *NodeInteractor {
	return &NodeInteractor{pr, gr, er}
}

func (i *NodeInteractor) Create(groupId *valueobject.ID, s app.Node) (*app.Node, error) {
	slice, err := i.NodeRepo.Create(groupId, s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return slice, nil
}

func (i *NodeInteractor) Get(nodeId *valueobject.ID) (*app.Node, error) {
	node, err := i.NodeRepo.Get(nodeId)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (i *NodeInteractor) View(ids []valueobject.ID) (*app.NodeView, error) {
	nodesView, err := i.NodeRepo.View(ids)
	if err != nil {
		return nil, err
	}

	return nodesView, nil
}

func (i *NodeInteractor) Update(actorId *valueobject.ID, node app.FlatNode) error {
	member, err := i.GroupRepo.FindMemberByNodeId(node.Id, actorId)
	if err != nil {
		return err
	}

	if member.Role == app.UserReader {
		return errors.New("Fobidden, only admin or editor can edit node.")
	}

	err = i.NodeRepo.Update(node)
	if err != nil {
		return err
	}

	return nil
}

func (i *NodeInteractor) AttachExpression(nodeId *valueobject.ID, inExpr app.Expression) (*app.Expression, error) {
	if inExpr.Id == nil {
		inExpr.Value = strings.TrimSpace(inExpr.Value)

		if inExpr.Value == "" {
			return nil, errors.New("You need to specify expression id or value.")
		}
	}

	expression, err := i.NodeRepo.AttachExpression(nodeId, inExpr)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (i *NodeInteractor) DetachExpression(nodeId *valueobject.ID, expressionId *valueobject.ID) error {
	err := i.NodeRepo.DetachExpression(nodeId, expressionId)
	if err != nil {
		return err
	}

	return nil
}

func (i *NodeInteractor) AvailableTranslations(nodeId *valueobject.ID, expressionId *valueobject.ID) ([]*app.Translation, error) {
	translations, err := i.NodeRepo.AvailableTranslations(nodeId, expressionId)
	if err != nil {
		return nil, err
	}

	return translations, nil
}

func (i *NodeInteractor) AttachTranslation(nodeId *valueobject.ID, expressionId *valueobject.ID, inTranslation app.Translation) (*app.Translation, error) {
	if inTranslation.Id == nil {
		inTranslation.Value = strings.TrimSpace(inTranslation.Value)

		if inTranslation.Value == "" {
			return nil, errors.New("You need to specify translation id or value.")
		}
	}

	translation, err := i.NodeRepo.AttachTranslation(nodeId, expressionId, inTranslation)
	if err != nil {
		return nil, err
	}

	return translation, nil
}

func (i *NodeInteractor) DetachTranslation(nodeId *valueobject.ID, translationId *valueobject.ID) error {
	err := i.NodeRepo.DetachTranslation(nodeId, translationId)
	if err != nil {
		return err
	}

	return nil
}

func (i *NodeInteractor) AttachText(nodeId *valueobject.ID, inText app.Text) (*app.Text, error) {
	text, err := i.NodeRepo.AttachText(nodeId, inText)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func (i *NodeInteractor) DetachText(actorId *valueobject.ID, nodeId *valueobject.ID) error {
	err := i.NodeRepo.DetachText(nodeId)
	if err != nil {
		return err
	}

	return nil
}
