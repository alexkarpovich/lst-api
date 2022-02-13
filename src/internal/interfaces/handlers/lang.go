package handlers

import (
	"net/http"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type LanguageInteractor interface {
	List() ([]*domain.Language, error)
	ListTranscriptionTypes(string) ([]*domain.TranscriptionType, error)
}

type languageHandler struct {
	BaseHanlder
	languageInteractor LanguageInteractor
}

func ConfigureLangHandler(li LanguageInteractor, r *mux.Router) {
	h := &languageHandler{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		languageInteractor: li,
	}

	h.router.HandleFunc("/langs", h.List()).Methods("GET")
	h.router.HandleFunc("/langs/{code:[a-z]{2}}/transcription-types", h.ListTranscriptionTypes()).Methods("GET")
}

func (i *languageHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		languages, err := i.languageInteractor.List()
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, languages, http.StatusOK)
	}
}

func (i *languageHandler) ListTranscriptionTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		transcriptionTypes, err := i.languageInteractor.ListTranscriptionTypes(vars["code"])
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, transcriptionTypes, http.StatusOK)
	}
}
