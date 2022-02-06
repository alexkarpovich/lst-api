package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type TranslationInteractor interface {
	AttachTranscription(*valueobject.ID, *valueobject.ID, *valueobject.ID) error
	DetachTranscription(*valueobject.ID, *valueobject.ID, *valueobject.ID) error
}

type translationHanlder struct {
	BaseHanlder
	translationInteractor TranslationInteractor
}

func ConfigureTranslationHandler(ti TranslationInteractor, r *mux.Router) {
	h := &translationHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		translationInteractor: ti,
	}

	h.router.HandleFunc("/translations/{translation_id}/transcriptions/{transcription_id}", h.AttachTranscription()).Methods("POST")
	h.router.HandleFunc("/translations/{translation_id}/transcriptions/{transcription_id}", h.DetachTranscription()).Methods("DELETE")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
}

func (i *translationHanlder) AttachTranscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		translationIdArg, err := strconv.Atoi(vars["translation_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid translation id", http.StatusBadRequest)
			return
		}
		translationId := valueobject.ID(translationIdArg)

		transcriptionIdArg, err := strconv.Atoi(vars["transcription_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid transcription id", http.StatusBadRequest)
			return
		}
		transcriptionId := valueobject.ID(transcriptionIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		err = i.translationInteractor.AttachTranscription(user.Id, &translationId, &transcriptionId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *translationHanlder) DetachTranscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		translationIdArg, err := strconv.Atoi(vars["translation_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid translation id", http.StatusBadRequest)
			return
		}
		translationId := valueobject.ID(translationIdArg)

		transcriptionIdArg, err := strconv.Atoi(vars["transcription_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid transcription id", http.StatusBadRequest)
			return
		}
		transcriptionId := valueobject.ID(transcriptionIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		err = i.translationInteractor.DetachTranscription(user.Id, &translationId, &transcriptionId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
