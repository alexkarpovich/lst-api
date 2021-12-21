package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type SliceInteractor interface {
	Get(*valueobject.ID) (*app.Slice, error)
	AttachExpression(*valueobject.ID, *app.Expression) (*app.Expression, error)
	DetachExpression(*valueobject.ID, *valueobject.ID) error
	AttachTranslation(*valueobject.ID, *valueobject.ID, *app.Translation) (*app.Translation, error)
	DetachTranslation(*valueobject.ID, *valueobject.ID) error
}

type sliceHanlder struct {
	BaseHanlder
	sliceInteractor SliceInteractor
}

func ConfigureSliceHandler(fi SliceInteractor, r *mux.Router) {
	h := &sliceHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		sliceInteractor: fi,
	}

	h.router.HandleFunc("/me/slice/{sliceId}", h.Get()).Methods("GET")
	h.router.HandleFunc("/me/slice/{sliceId}/attach-expression", h.AttachExpression()).Methods("POST")
	h.router.HandleFunc("/me/slice/{sliceId}/detach-expression/{expressionId}", h.DetachExpression()).Methods("POST")
	h.router.HandleFunc("/me/slice/{sliceId}/attach-translation", h.AttachTranslation()).Methods("POST")
	h.router.HandleFunc("/me/slice/{sliceId}/detach-translation/{translationId}", h.DetachTranslation()).Methods("POST")
	// h.router.HandleFunc("/me/profile", h.ProfileList()).Methods("GET")
	// i.router.HandleFunc("/login", i.Login()).Methods("POST")
}

func (i *sliceHanlder) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		sliceIdArg, err := strconv.Atoi(vars["sliceId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}
		sliceId := valueobject.ID(sliceIdArg)

		slice, err := i.sliceInteractor.Get(&sliceId)
		if err != nil {
			utils.SendJsonError(w, "Attach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, slice, http.StatusOK)
	}
}

func (i *sliceHanlder) AttachExpression() http.HandlerFunc {
	type request struct {
		Id    *valueobject.ID `json:"id"`
		Value string          `json:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		sliceIdArg, err := strconv.Atoi(vars["sliceId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		sliceId := valueobject.ID(sliceIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		expression := &app.Expression{
			Id:    s.Id,
			Value: s.Value,
		}

		expression, err = i.sliceInteractor.AttachExpression(&sliceId, expression)
		if err != nil {
			utils.SendJsonError(w, "Attach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, expression, http.StatusOK)
	}
}

func (i *sliceHanlder) DetachExpression() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		sliceIdArg, err := strconv.Atoi(vars["sliceId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		expressionIdArg, err := strconv.Atoi(vars["expressionId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid expression id", http.StatusBadRequest)
			return
		}
		sliceId := valueobject.ID(sliceIdArg)
		expressionId := valueobject.ID(expressionIdArg)

		err = i.sliceInteractor.DetachExpression(&sliceId, &expressionId)
		if err != nil {
			utils.SendJsonError(w, "Detach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *sliceHanlder) AttachTranslation() http.HandlerFunc {
	type translation struct {
		Id             *valueobject.ID `json:"id"`
		Value          string          `json:"value"`
		Transcriptions []string        `json:"transcriptions"`
		Comment        string          `json:"comment"`
	}
	type request struct {
		ExpressionId *valueobject.ID `json:"expressionId"`
		Translation  translation     `json:"translation"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		sliceIdArg, err := strconv.Atoi(vars["sliceId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid slice id", http.StatusBadRequest)
			return
		}

		sliceId := valueobject.ID(sliceIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		translation := &app.Translation{
			Id:             s.Translation.Id,
			Transcriptions: s.Translation.Transcriptions,
			Comment:        s.Translation.Comment,
			Value:          s.Translation.Value,
		}

		translation, err = i.sliceInteractor.AttachTranslation(&sliceId, s.ExpressionId, translation)
		if err != nil {
			utils.SendJsonError(w, "Attach translation error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, translation, http.StatusOK)
	}
}

func (i *sliceHanlder) DetachTranslation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		sliceIdArg, err := strconv.Atoi(vars["sliceId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid folder id", http.StatusBadRequest)
			return
		}
		translationIdArg, err := strconv.Atoi(vars["translationId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid translation id", http.StatusBadRequest)
			return
		}
		sliceId := valueobject.ID(sliceIdArg)
		translationId := valueobject.ID(translationIdArg)

		err = i.sliceInteractor.DetachTranslation(&sliceId, &translationId)
		if err != nil {
			utils.SendJsonError(w, "Detach expression error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
