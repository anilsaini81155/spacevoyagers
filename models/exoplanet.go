package models

import (
	"database/sql"
	"errors"
	_ "fmt"

	"github.com/anilsaini81155/spacevoyagers/db"
)

type ExoplanetType string

const (
	GasGiant    ExoplanetType = "GasGiant"
	Terrestrial ExoplanetType = "Terrestrial"
)

type Exoplanet struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Distance    float64       `json:"distance"`
	Radius      float64       `json:"radius"`
	Mass        float64       `json:"mass,omitempty"`
	Type        ExoplanetType `json:"type"`
}

// AddExoplanet inserts a new exoplanet into the database
func AddExoplanet(exoplanet *Exoplanet) error {

	DB, err := db.GetDB() // Get singleton DB instance
	if err != nil {
		return err
	}

	query := `INSERT INTO exoplanets (name, description, distance, radius, mass, type) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, exoplanet.Name, exoplanet.Description, exoplanet.Distance, exoplanet.Radius, exoplanet.Mass, exoplanet.Type)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	exoplanet.ID = int(id)
	return nil
}

// GetAllExoplanets retrieves all exoplanets from the database
func GetAllExoplanets() ([]Exoplanet, error) {

	DB, dberr := db.GetDB() // Get singleton DB instance
	if dberr != nil {
		return nil, dberr
	}

	query := `SELECT id, name, description, distance, radius, mass, type FROM exoplanets`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exoplanets []Exoplanet
	for rows.Next() {
		var exoplanet Exoplanet
		if err := rows.Scan(&exoplanet.ID, &exoplanet.Name, &exoplanet.Description, &exoplanet.Distance, &exoplanet.Radius, &exoplanet.Mass, &exoplanet.Type); err != nil {
			return nil, err
		}
		exoplanets = append(exoplanets, exoplanet)
	}
	return exoplanets, nil
}

// GetExoplanetByID retrieves a specific exoplanet by its ID
func GetExoplanetByID(id int) (*Exoplanet, error) {

	DB, dberr := db.GetDB() // Get singleton DB instance
	if dberr != nil {
		return nil, dberr
	}

	query := `SELECT id, name, description, distance, radius, mass, type FROM exoplanets WHERE id = ?`
	var exoplanet Exoplanet
	err := DB.QueryRow(query, id).Scan(&exoplanet.ID, &exoplanet.Name, &exoplanet.Description, &exoplanet.Distance, &exoplanet.Radius, &exoplanet.Mass, &exoplanet.Type)
	if err == sql.ErrNoRows {
		return nil, errors.New("exoplanet not found")
	} else if err != nil {
		return nil, err
	}
	return &exoplanet, nil
}

// DeleteExoplanet removes an exoplanet from the database
func DeleteExoplanet(id int) error {

	DB, dberr := db.GetDB() // Get singleton DB instance
	if dberr != nil {
		return dberr
	}

	query := `DELETE FROM exoplanets WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

// UpdateExoplanet updates an existing exoplanet in the database
func UpdateExoplanet(exoplanet *Exoplanet) error {

	DB, dberr := db.GetDB() // Get singleton DB instance
	if dberr != nil {
		return dberr
	}

	query := `UPDATE exoplanets SET name = ?, description = ?, distance = ?, radius = ?, mass = ?, type = ? WHERE id = ?`
	_, err := DB.Exec(query, exoplanet.Name, exoplanet.Description, exoplanet.Distance, exoplanet.Radius, exoplanet.Mass, exoplanet.Type, exoplanet.ID)
	return err
}

// Validate ensures that the planet details are correct
func (p *Exoplanet) Validate() error {
	if p.Name == "" || p.Description == "" || p.Distance <= 0 || p.Radius <= 0 {
		return errors.New("invalid exoplanet data")
	}
	if p.Type == Terrestrial && p.Mass <= 0 {
		return errors.New("mass required for terrestrial exoplanets")
	}
	return nil
}

// CalculateGravity returns the gravity of the exoplanet based on type
func (p *Exoplanet) CalculateGravity() float64 {
	if p.Type == GasGiant {
		return 0.5 / (p.Radius * p.Radius)
	}
	return p.Mass / (p.Radius * p.Radius)
}

// FuelEstimation calculates the fuel based on distance, gravity, and crew capacity
func (p *Exoplanet) FuelEstimation(crewCapacity int) (float64, error) {
	if crewCapacity <= 0 {
		return 0, errors.New("invalid crew capacity")
	}
	gravity := p.CalculateGravity()
	fuel := (p.Distance / (gravity * gravity)) * float64(crewCapacity)
	return fuel, nil
}
