package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type TrainingInteractor interface {
	Create(app.Training) (*app.Training, error)
}

type trainingHandler struct {
	BaseHanlder
	trainingInteractor TrainingInteractor
}

func ConfigureTrainingHandler(ti TrainingInteractor, r *mux.Router) {
	h := &trainingHandler{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		trainingInteractor: ti,
	}

	h.router.HandleFunc("/me/trainings", h.Create()).Methods("GET")
}

func (i *trainingHandler) Create() http.HandlerFunc {
	type request struct {
		Type  app.TrainingType `json:"type"`
		Nodes []valueobject.ID `json:"nodes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		inTraining := app.Training{
			OwnerId: user.Id,
			Type:    s.Type,
			Nodes:   s.Nodes,
		}

		training, err := i.trainingInteractor.Create(inTraining)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, training, http.StatusOK)
	}
}
