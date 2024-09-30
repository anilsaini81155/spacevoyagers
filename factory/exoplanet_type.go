package factory

// Base Exoplanet interface
type Exoplanet interface {
	GetName() string
	GetDescription() string
	GetDistance() float64
	GetRadius() float64
	GetMass() float64
	GetType() string
}

// GasGiant struct
type GasGiant struct {
	Name        string
	Description string
	Distance    float64
	Radius      float64
}

func (g *GasGiant) GetName() string        { return g.Name }
func (g *GasGiant) GetDescription() string { return g.Description }
func (g *GasGiant) GetDistance() float64   { return g.Distance }
func (g *GasGiant) GetRadius() float64     { return g.Radius }
func (g *GasGiant) GetMass() float64       { return 0 } // Mass is not applicable for GasGiant
func (g *GasGiant) GetType() string        { return "GasGiant" }

// Terrestrial struct
type Terrestrial struct {
	Name        string
	Description string
	Distance    float64
	Radius      float64
	Mass        float64
}

func (t *Terrestrial) GetName() string        { return t.Name }
func (t *Terrestrial) GetDescription() string { return t.Description }
func (t *Terrestrial) GetDistance() float64   { return t.Distance }
func (t *Terrestrial) GetRadius() float64     { return t.Radius }
func (t *Terrestrial) GetMass() float64       { return t.Mass }
func (t *Terrestrial) GetType() string        { return "Terrestrial" }
