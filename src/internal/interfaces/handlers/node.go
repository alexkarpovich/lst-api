package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type NodeInteractor interface {
	Get(*valueobject.ID) (*app.Node, error)
	View([]valueobject.ID) (*app.NodeView, error)
	Update(*valueobject.ID, app.FlatNode) error
	AttachExpression(*valueobject.ID, app.Expression) (*app.Expression, error)
	DetachExpression(*valueobject.ID, *valueobject.ID) error
	AvailableTranslations(*valueobject.ID, *valueobject.ID) ([]*app.Translation, error)
	AttachTranslation(*valueobject.ID, *valueobject.ID, app.Translation) (*app.Translation, error)
	DetachTranslation(*valueobject.ID, *valueobject.ID) error
	AttachText(*valueobject.ID, app.Text) (*app.Text, error)
	DetachText(*valueobject.ID, *valueobject.ID) error
}

type nodeHandler struct {
	BaseHanlder
	NodeInteractor NodeInteractor
}

func ConfigureNodeHandler(fi NodeInteractor, r *mux.Router) {
	h := &nodeHandler{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		NodeInteractor: fi,
	}

	h.router.HandleFunc("/me/nodes", h.View()).
		Queries("ids", "{[0-9]+}").
		Methods("GET")
	h.router.HandleFunc("/me/nodes/{node_id}", h.Get()).Methods("GET")
	h.router.HandleFunc("/me/nodes/{node_id}", h.Update()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/translations", h.AvailableTranslations()).
		Queries("expression_id", "{[0-9]+}").
		Methods("GET")
	h.router.HandleFunc("/me/nodes/{node_id}/attach-expression", h.AttachExpression()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/detach-expression/{expression_id}", h.DetachExpression()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/attach-translation", h.AttachTranslation()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/detach-translation/{translation_id}", h.DetachTranslation()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/attach-text", h.AttachText()).Methods("POST")
	h.router.HandleFunc("/me/nodes/{node_id}/detach-text", h.DetachText()).Methods("POST")
}

func (i *nodeHandler) View() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		queryParams := r.URL.Query()
		ids := []valueobject.ID{}

		for _, idStr := range queryParams["ids"] {
			id, err := strconv.Atoi(string(idStr))
			if err != nil {
				utils.SendJsonError(w, err, http.StatusBadRequest)
				return
			}
			ids = append(ids, valueobject.ID(id))
		}

		if len(ids) == 0 {
			utils.SendJsonError(w, errors.New("No node ids specified."), http.StatusBadRequest)
			return
		}

		slice, err := i.NodeInteractor.View(ids)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, slice, http.StatusOK)
	}
}

func (i *nodeHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		slice, err := i.NodeInteractor.Get(&nodeId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, slice, http.StatusOK)
	}
}

func (i *nodeHandler) Update() http.HandlerFunc {
	type request struct {
		Name       string             `json:"name,omitempty"`
		Path       string             `json:"path,omitempty"`
		Visibility app.NodeVisibility `json:"visibility,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error slice create context")
			return
		}

		inNode := app.FlatNode{
			Id:         &nodeId,
			Name:       s.Name,
			Path:       s.Path,
			Visibility: s.Visibility,
		}

		err = i.NodeInteractor.Update(user.Id, inNode)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *nodeHandler) AttachExpression() http.HandlerFunc {
	type request struct {
		Id    *valueobject.ID `json:"id"`
		Value string          `json:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		inExpr := app.Expression{
			Id:    s.Id,
			Value: s.Value,
		}

		expression, err := i.NodeInteractor.AttachExpression(&nodeId, inExpr)
		if err != nil {
			utils.SendJsonError(w, "Attach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, expression, http.StatusOK)
	}
}

func (i *nodeHandler) DetachExpression() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		expressionIdArg, err := strconv.Atoi(vars["expression_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid expression id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)
		expressionId := valueobject.ID(expressionIdArg)

		err = i.NodeInteractor.DetachExpression(&nodeId, &expressionId)
		if err != nil {
			utils.SendJsonError(w, "Detach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *nodeHandler) AvailableTranslations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid node id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		queryParams := r.URL.Query()
		expressionIdArg, err := strconv.Atoi(queryParams.Get("expression_id"))
		if err != nil {
			utils.SendJsonError(w, "Invalid expression id", http.StatusBadRequest)
			return
		}
		expressionId := valueobject.ID(expressionIdArg)

		translations, err := i.NodeInteractor.AvailableTranslations(&nodeId, &expressionId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, translations, http.StatusOK)
	}
}

func (i *nodeHandler) AttachTranslation() http.HandlerFunc {
	type translation struct {
		Id      *valueobject.ID `json:"id"`
		Value   string          `json:"value"`
		Comment string          `json:"comment"`
	}
	type request struct {
		ExpressionId *valueobject.ID `json:"expressionId"`
		Translation  translation     `json:"translation"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}

		nodeId := valueobject.ID(nodeIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		inTranslation := app.Translation{
			Id:      s.Translation.Id,
			Comment: s.Translation.Comment,
			Value:   s.Translation.Value,
		}

		translation, err := i.NodeInteractor.AttachTranslation(&nodeId, s.ExpressionId, inTranslation)
		if err != nil {
			utils.SendJsonError(w, "Attach translation error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, translation, http.StatusOK)
	}
}

func (i *nodeHandler) DetachTranslation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		translationIdArg, err := strconv.Atoi(vars["translation_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid translation id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)
		translationId := valueobject.ID(translationIdArg)

		err = i.NodeInteractor.DetachTranslation(&nodeId, &translationId)
		if err != nil {
			utils.SendJsonError(w, "Detach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *nodeHandler) AttachText() http.HandlerFunc {
	type request struct {
		Id      *valueobject.ID `json:"id"`
		Title   string          `json:"title"`
		Content string          `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error slice create context")
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		inText := app.Text{
			Id:       s.Id,
			AuthorId: user.Id,
			Content:  s.Content,
		}

		text, err := i.NodeInteractor.AttachText(&nodeId, inText)
		if err != nil {
			utils.SendJsonError(w, "Attach text error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, text, http.StatusOK)
	}
}

func (i *nodeHandler) DetachText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		nodeIdArg, err := strconv.Atoi(vars["node_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}
		nodeId := valueobject.ID(nodeIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error slice create context")
			return
		}

		err = i.NodeInteractor.DetachText(user.Id, &nodeId)
		if err != nil {
			utils.SendJsonError(w, "Detach text error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
