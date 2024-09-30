package factory

import (
	"errors"

	"github.com/anilsaini81155/spacevoyagers/models"
)

// CreateExoplanet creates an exoplanet based on its type.
func CreateExoplanet(exoplanetType, name, description string, distance, radius, mass float64) (models.Exoplanet, error) {
	switch exoplanetType {
	case "Terrestrial":
		return models.Exoplanet{
			Name:        name,
			Description: description,
			Distance:    distance,
			Radius:      radius,
			Mass:        mass,
			Type:        "Terrestrial",
		}, nil
	case "GasGiant":
		return models.Exoplanet{
			Name:        name,
			Description: description,
			Distance:    distance,
			Radius:      radius,
			Mass:        mass,
			Type:        "GasGiant",
		}, nil
	default:
		return models.Exoplanet{}, errors.New("unknown exoplanet type")
	}
}
