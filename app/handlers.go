package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/feastM/HatParty/config"
	"github.com/feastM/HatParty/database"
	"github.com/feastM/HatParty/models"
	"github.com/google/uuid"
)

func StartParty(w http.ResponseWriter, r *http.Request) {
	hatsRequested := r.URL.Query().Get("hatsRequested")
	intVar, err := strconv.Atoi(hatsRequested)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "No number of hats specified"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}
	party := models.Party{
		Id:            uuid.New().String(),
		Status:        int(models.Active),
		HatsRequested: intVar,
		Hats:          nil,
	}

	if party.HatsRequested > config.Cfg.MaxHatsPerParty {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Number of hats requested is too big"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	err = database.AddParty(database.DB, party)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot add party"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["partyId"] = party.Id
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func StopParty(w http.ResponseWriter, r *http.Request) {
	partyId := r.URL.Query().Get("partyId")
	if partyId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "No partyId specified"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	err := database.StopParty(database.DB, partyId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot stop party"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Party stopped successfully"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func Init(w http.ResponseWriter, r *http.Request) {
	hats := generateHats(config.Cfg.InitNumberOfHats)
	database.InsertHats(database.DB, "Hats", hats)
}

func HandleRequests() {
	http.HandleFunc("/start", StartParty)
	http.HandleFunc("/stop", StopParty)
	http.HandleFunc("/init", Init)
}

func generateHats(numberOfHats int) []models.Hat {
	var hats []models.Hat
	for i := 0; i < numberOfHats; i++ {
		hats = append(hats, models.Hat{
			Id:       i,
			UsedBy:   nil,
			Priority: i,
			Cleaning: nil,
		})
	}

	return hats
}
