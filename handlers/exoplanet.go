package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/anilsaini81155/spacevoyagers/models"
	"github.com/gorilla/mux"
)

// var exoplanets []models.Exoplanet

// var idCounter = 1

// CreateExoplanet handles adding a new exoplanet
func CreateExoplanet(w http.ResponseWriter, r *http.Request) {
	var exoplanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// exoplanet.ID = idCounter
	// idCounter++

	if err := exoplanet.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.AddExoplanet(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// exoplanets = append(exoplanets, exoplanet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exoplanet)
}

// ListExoplanets handles listing all exoplanets
/*
	//sample input query params
	GET /exoplanets?sort=distance
	GET /exoplanets?sort=name
	GET /exoplanets?type=GasGiant

	GET /exoplanets?min_distance=1000&max_distance=5000
	GET /exoplanets?type=Terrestrial&sort=distance
*/

func ListExoplanets(w http.ResponseWriter, r *http.Request) {
	/*
		//below code to get for all listing
		exoplanets, err := models.GetAllExoplanets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(exoplanets)
	*/

	// Parse query parameters for filtering and sorting
	sortField := r.URL.Query().Get("sort")
	filterType := r.URL.Query().Get("type")
	filterMinDistance := r.URL.Query().Get("min_distance")
	filterMaxDistance := r.URL.Query().Get("max_distance")

	// Construct SQL query
	query := "SELECT id, name, description, distance, radius, mass, type FROM exoplanets WHERE 1=1"
	args := []interface{}{}

	// Apply filters based on the query parameters
	if filterType != "" {
		query += " AND type = ?"
		args = append(args, filterType)
	}

	if filterMinDistance != "" {
		minDistance, err := strconv.ParseFloat(filterMinDistance, 64)
		if err == nil {
			query += " AND distance >= ?"
			args = append(args, minDistance)
		}
	}

	if filterMaxDistance != "" {
		maxDistance, err := strconv.ParseFloat(filterMaxDistance, 64)
		if err == nil {
			query += " AND distance <= ?"
			args = append(args, maxDistance)
		}
	}

	// Apply sorting based on the query parameters
	if sortField != "" {
		switch sortField {
		case "name":
			query += " ORDER BY name"
		case "distance":
			query += " ORDER BY distance"
		case "radius":
			query += " ORDER BY radius"
		case "type":
			query += " ORDER BY type"
		default:
			query += " ORDER BY id"
		}
	} else {
		query += " ORDER BY id"
	}

	// Execute the query
	rows, err := models.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying exoplanets: %v", err)
		http.Error(w, "Error retrieving exoplanets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse the results into a slice of exoplanets
	var exoplanets []models.Exoplanet
	for rows.Next() {
		var exoplanet models.Exoplanet
		err := rows.Scan(&exoplanet.ID, &exoplanet.Name, &exoplanet.Description, &exoplanet.Distance, &exoplanet.Radius, &exoplanet.Mass, &exoplanet.Type)
		if err != nil {
			log.Printf("Error scanning exoplanet: %v", err)
			http.Error(w, "Error retrieving exoplanets", http.StatusInternalServerError)
			return
		}
		exoplanets = append(exoplanets, exoplanet)
	}

	// Send the response back as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exoplanets)
}

// GetExoplanetByID handles fetching an exoplanet by its ID
func GetExoplanetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	/*
		for _, exoplanet := range exoplanets {
			if exoplanet.ID == id {
				json.NewEncoder(w).Encode(exoplanet)
				return
			}
		}

		http.Error(w, "Exoplanet not found", http.StatusNotFound)
	*/
	exoplanet, err := models.GetExoplanetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(exoplanet)
}

// UpdateExoplanet handles updating an exoplanet by its ID
func UpdateExoplanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	/*
		for i, exoplanet := range exoplanets {
			if exoplanet.ID == id {
				var updatedExoplanet models.Exoplanet
				if err := json.NewDecoder(r.Body).Decode(&updatedExoplanet); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				updatedExoplanet.ID = exoplanet.ID

				if err := updatedExoplanet.Validate(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				exoplanets[i] = updatedExoplanet
				json.NewEncoder(w).Encode(updatedExoplanet)
				return
			}
		}

		http.Error(w, "Exoplanet not found", http.StatusNotFound)
	*/

	var updatedExoplanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&updatedExoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatedExoplanet.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedExoplanet.ID = id

	if err := models.UpdateExoplanet(&updatedExoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedExoplanet)
}

// DeleteExoplanet handles removing an exoplanet by its ID
func DeleteExoplanet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	/*
		for i, exoplanet := range exoplanets {
			if exoplanet.ID == id {
				exoplanets = append(exoplanets[:i], exoplanets[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		http.Error(w, "Exoplanet not found", http.StatusNotFound)
	*/
	if err := models.DeleteExoplanet(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// FuelEstimation calculates the fuel required for a trip to the exoplanet
func FuelEstimation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	crewCapacity, _ := strconv.Atoi(r.URL.Query().Get("crewCapacity"))

	exoplanet, err := models.GetExoplanetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fuel, err := exoplanet.FuelEstimation(crewCapacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"fuel": fuel})

}
