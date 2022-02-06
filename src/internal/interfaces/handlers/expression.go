package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type ExpressionInteractor interface {
	Search(string, string) ([]*domain.Expression, error)
	CreateTranscription(*valueobject.ID, domain.Transcription) (*domain.Transcription, error)
	GetTranscriptionMap(*valueobject.ID, *valueobject.ID) (map[string][]*domain.TranscriptionItem, error)
}

type expressionHanlder struct {
	BaseHanlder
	expressionInteractor ExpressionInteractor
}

func ConfigureExpressionHandler(ei ExpressionInteractor, r *mux.Router) {
	h := &expressionHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		expressionInteractor: ei,
	}

	h.router.
		HandleFunc("/x", h.Search()).
		Queries("lang", "{[a-z]{2}}").
		Queries("search", "{.+}").
		Methods("GET")
	h.router.
		HandleFunc("/x/{expression_id}/transcription-parts", h.GetTranscriptionMap()).
		Queries("type", "{\\d+}").
		Methods("GET")
	h.router.HandleFunc("/x/{expression_id}/transcriptions", h.CreateTranscription()).Methods("POST")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
}

func (i *expressionHanlder) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := r.FormValue("lang")
		searchValue := r.FormValue("search")

		expressions, err := i.expressionInteractor.Search(langCode, searchValue)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, expressions, http.StatusOK)
	}
}

func (i *expressionHanlder) CreateTranscription() http.HandlerFunc {
	type request struct {
		Type  *valueobject.ID `json:"type" db:"type"`
		Value string          `json:"value" db:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

		vars := mux.Vars(r)
		expressionIdArg, err := strconv.Atoi(vars["expression_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		expressionId := valueobject.ID(expressionIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		inTranscription := domain.Transcription{
			Type:  s.Type,
			Value: s.Value,
		}

		transcription, err := i.expressionInteractor.CreateTranscription(&expressionId, inTranscription)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, transcription, http.StatusOK)
	}
}

func (i *expressionHanlder) GetTranscriptionMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeIdArg, err := strconv.Atoi(r.FormValue("type"))
		if err != nil {
			utils.SendJsonError(w, "Invalid type id", http.StatusBadRequest)
			return
		}
		typeId := valueobject.ID(typeIdArg)

		vars := mux.Vars(r)
		expressionIdArg, err := strconv.Atoi(vars["expression_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		expressionId := valueobject.ID(expressionIdArg)

		transcriptionMap, err := i.expressionInteractor.GetTranscriptionMap(&expressionId, &typeId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, transcriptionMap, http.StatusOK)
	}
}
