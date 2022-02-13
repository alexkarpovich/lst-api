package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type TrainingInteractor interface {
	GetOrCreate(app.Training) (*app.Training, error)
	Get(*valueobject.ID, *valueobject.ID) (*app.Training, error)
	List(*valueobject.ID) ([]*app.Training, error)
	Reset(*valueobject.ID, *valueobject.ID) error
	Next(*valueobject.ID, *valueobject.ID) (*app.TrainingItem, error)
	GetItem(*valueobject.ID, *valueobject.ID) (*app.TrainingItem, error)
	ItemAnswers(*valueobject.ID, *valueobject.ID) ([]*app.TrainingAnswer, error)
	MarkItemAsComplete(*valueobject.ID, *valueobject.ID) error
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

	h.router.HandleFunc("/me/trainings", h.List()).Methods("GET")
	h.router.HandleFunc("/me/trainings", h.Create()).Methods("POST")
	h.router.HandleFunc("/me/trainings/{training_id}", h.Get()).Methods("GET")
	h.router.HandleFunc("/me/trainings/{training_id}/next", h.Next()).Methods("GET")
	h.router.HandleFunc("/me/trainings/{training_id}/reset", h.Reset()).Methods("POST")
	h.router.HandleFunc("/me/training-items/{item_id}", h.GetItem()).Methods("GET")
	h.router.HandleFunc("/me/training-items/{item_id}/answers", h.ItemAnswers()).Methods("GET")
	h.router.HandleFunc("/me/training-items/{item_id}/complete", h.Complete()).Methods("POST")
}

func (i *trainingHandler) Create() http.HandlerFunc {
	type request struct {
		Type   app.TrainingType `json:"type"`
		Slices []valueobject.ID `json:"slices"`
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
			Slices:  s.Slices,
		}

		training, err := i.trainingInteractor.GetOrCreate(inTraining)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, training, http.StatusOK)
	}
}

func (i *trainingHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		trainingIdArg, err := strconv.Atoi(vars["training_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid training id", http.StatusBadRequest)
			return
		}
		trainingId := valueobject.ID(trainingIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		training, err := i.trainingInteractor.Get(user.Id, &trainingId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, training, http.StatusOK)
	}
}

func (i *trainingHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		trainings, err := i.trainingInteractor.List(user.Id)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, trainings, http.StatusOK)
	}
}

func (i *trainingHandler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		trainingIdArg, err := strconv.Atoi(vars["training_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid training id", http.StatusBadRequest)
			return
		}
		trainingId := valueobject.ID(trainingIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		err = i.trainingInteractor.Reset(user.Id, &trainingId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *trainingHandler) Next() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		trainingIdArg, err := strconv.Atoi(vars["training_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid training id", http.StatusBadRequest)
			return
		}
		trainingId := valueobject.ID(trainingIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		trainingItem, err := i.trainingInteractor.Next(user.Id, &trainingId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, trainingItem, http.StatusOK)
	}
}

func (i *trainingHandler) GetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemIdArg, err := strconv.Atoi(vars["item_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid training id", http.StatusBadRequest)
			return
		}
		itemId := valueobject.ID(itemIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		trainingItem, err := i.trainingInteractor.GetItem(user.Id, &itemId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, trainingItem, http.StatusOK)
	}
}

func (i *trainingHandler) ItemAnswers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemIdArg, err := strconv.Atoi(vars["item_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid item id", http.StatusBadRequest)
			return
		}
		itemId := valueobject.ID(itemIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		answers, err := i.trainingInteractor.ItemAnswers(user.Id, &itemId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, answers, http.StatusOK)
	}
}

func (i *trainingHandler) Complete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemIdArg, err := strconv.Atoi(vars["item_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid training id", http.StatusBadRequest)
			return
		}
		itemId := valueobject.ID(itemIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user context")
			return
		}

		err = i.trainingInteractor.MarkItemAsComplete(user.Id, &itemId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
