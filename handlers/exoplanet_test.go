package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/anilsaini81155/spacevoyagers/db"
	"github.com/anilsaini81155/spacevoyagers/models"
	_ "github.com/go-sql-driver/mysql" // for MySQL driver
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("GO_ENV", "TEST") // Set the environment to test
	loadEnvForTests()           // Load .envtest or test environment variables
}

func loadEnvForTests() {

	envPath, _ := filepath.Abs("../.envtest")
	log.Printf("Loading environment file from: %s", envPath)

	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("Warning: .env.test file not found or failed to load, skipping...")
	}
}

func TestMain(m *testing.M) {
	// Load environment variables for the test
	err := godotenv.Load("../.envtest")
	if err != nil {
		log.Println("Warning: .envtest file not found, skipping...")
	}

	// Initialize the database connection
	DB, dberr := db.GetDB()
	if dberr != nil {
		log.Fatalf("Error connecting to the database: %v", dberr)
	}

	models.SetDB(DB)
	// Run migrations before tests
	models.RunMigrations()

	// Run the tests
	code := m.Run()

	// Exit with the proper code after tests
	os.Exit(code)
}

// TestCreateExoplanet tests the CreateExoplanet handler.
func TestCreateExoplanet(t *testing.T) {
	loadEnvForTests()
	// Mock request body
	requestBody := `{
        "name": "Planet X",
        "description": "A mysterious planet.",
        "distance": 4500,
        "radius": 50,
        "mass": 5,
        "type": "Terrestrial"
    }`

	// loadEnv()

	req, err := http.NewRequest("POST", "/exoplanets", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateExoplanet)

	// Call the handler
	handler.ServeHTTP(rr, req)
	t.Logf("Response Body: %s", rr.Body.String())

	// Check the response code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the fields of the created exoplanet
	assert.Equal(t, "Planet X", response["name"])
	assert.Equal(t, "A mysterious planet.", response["description"])
	assert.Equal(t, 4500, int(response["distance"].(float64)))
	assert.Equal(t, 50, int(response["radius"].(float64)))
	assert.Equal(t, 5, int(response["mass"].(float64)))
	assert.Equal(t, "Terrestrial", response["type"])

	// assert.Equal(t, "Exoplanet added successfully", response["message"])
}

// TestListExoplanets tests the ListExoplanets handler.
func TestListExoplanets(t *testing.T) {

	loadEnvForTests()

	req, err := http.NewRequest("GET", "/exoplanets", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListExoplanets)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body (assuming the response should contain a list)
	var response []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Assuming at least one exoplanet exists
	assert.GreaterOrEqual(t, len(response), 1)
}

// TestListExoplanetsWithSortingAndFiltering tests the ListExoplanets handler with sorting and filtering.
func TestListExoplanetsWithSortingAndFiltering(t *testing.T) {
	// Test with sorting by distance and filtering by type

	loadEnvForTests()

	req, err := http.NewRequest("GET", "/exoplanets?sort=distance&type=Terrestrial", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListExoplanets)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Assuming the database has planets of type 'Terrestrial'
	for _, planet := range response {
		assert.Equal(t, "Terrestrial", planet["type"])
	}
}

// TestGetExoplanetByID tests the GetExoplanetByID handler.
func TestGetExoplanetByID(t *testing.T) {

	loadEnvForTests()

	req, err := http.NewRequest("GET", "/exoplanets/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a path parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExoplanetByID)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Planet X", response["name"])
}

// TestUpdateExoplanet tests the UpdateExoplanet handler.
func TestUpdateExoplanet(t *testing.T) {

	loadEnvForTests()

	// Mock request body
	requestBody := `{
        "name": "Updated Planet",
        "description": "Updated description",
        "distance": 5000,
        "radius": 100,
        "mass": 10,
        "type": "Terrestrial"
    }`
	req, err := http.NewRequest("PUT", "/exoplanets/1", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add a path parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateExoplanet)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Updated Planet", response["name"])
	assert.Equal(t, "Updated description", response["description"])
	assert.Equal(t, 5000, int(response["distance"].(float64)))
	assert.Equal(t, 100, int(response["radius"].(float64)))
	assert.Equal(t, 10, int(response["mass"].(float64)))
	assert.Equal(t, "Terrestrial", response["type"])

	// assert.Equal(t, "Exoplanet updated successfully", response["message"])
}

// TestDeleteExoplanet tests the DeleteExoplanet handler.
func TestDeleteExoplanet(t *testing.T) {

	loadEnvForTests()

	req, err := http.NewRequest("DELETE", "/exoplanets/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a path parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteExoplanet)

	// Call the handler
	handler.ServeHTTP(rr, req)
	t.Logf("Response Body: %s", rr.Body.String())
	// Check the response code
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Check the response body
	// var response map[string]interface{}
	// err = json.Unmarshal(rr.Body.Bytes(), &response)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// assert.Equal(t, "Exoplanet deleted successfully", response["message"])
}

// TestFuelEstimation tests the FuelEstimation handler.
func TestFuelEstimation(t *testing.T) {

	loadEnvForTests()

	// Mock request for fuel estimation
	req, err := http.NewRequest("GET", "/exoplanets/1/fuel?crewCapacity=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a path parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FuelEstimation)

	// Call the handler
	handler.ServeHTTP(rr, req)

	t.Logf("Response Body: %s", rr.Body.String())

	// Check the response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Assuming fuel calculation logic is correct, assert response
	assert.NotNil(t, response["fuel"])
}
