package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/innoglobe/xmgo/internal/entity"
	postgresrepository "github.com/innoglobe/xmgo/internal/infrastructure/db/postgres"
	"github.com/innoglobe/xmgo/internal/infrastructure/server/handler"
	eventservice "github.com/innoglobe/xmgo/internal/service"
	"github.com/innoglobe/xmgo/internal/usecase"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"xmgo", "xmgopass", "127.0.0.1", 25432, "xmgo_db")
	companyName = generateRandomCompanyName()
)

func setupRouter(dsn string) *gin.Engine {
	// Initialize db conn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect to test database: %v\n", err)
	}

	repo := postgresrepository.NewPostgresRepository(db)
	usecase := usecase.NewCompanyUsecase(repo, &eventservice.NoOpProducer{})
	companyHandler := handler.NewCompanyHandler(usecase)

	r := gin.New()
	gin.SetMode(gin.TestMode)
	api := r.Group("/api")
	companyHandler.RegisterRoutes(api, "test-secret")
	return r
}

func TestCompanyHandler_CreateCompanyUnauthorized(t *testing.T) {
	router := setupRouter(dsn)

	company := entity.Company{
		Name:              "My test company",
		AmountOfEmployees: 5,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCompanyHandler_CreateCompany(t *testing.T) {
	router := setupRouter(dsn)

	company := entity.Company{
		Name:              companyName,
		Description:       "random description",
		AmountOfEmployees: 5,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)

	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCompanyHandler_CreateCompanyDuplicate(t *testing.T) {
	router := setupRouter(dsn)
	company := entity.Company{
		Name:              companyName,
		Description:       "random description",
		AmountOfEmployees: 5,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)

	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, fmt.Sprintf("{\"error\":\"Company with name %s already exists\"}", companyName), w.Body.String())
}

func TestCompanyHandler_CreateCompanyMissingFields(t *testing.T) {
	router := setupRouter(dsn)

	company := entity.Company{
		Name:        generateRandomCompanyName(),
		Description: "random description",
		//AmountOfEmployees: 5,
		//Registered:        true,
		Type: entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)

	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"details\":[{\"field\":\"AmountOfEmployees\",\"message\":\"Key: 'Company.AmountOfEmployees' Error:Field validation for 'AmountOfEmployees' failed on the 'required' tag\"},{\"field\":\"Registered\",\"message\":\"Key: 'Company.Registered' Error:Field validation for 'Registered' failed on the 'required' tag\"}],\"error\":\"Invalid input data\"}", w.Body.String())
}

func TestCompanyHandler_UpdateCompany(t *testing.T) {
	router := setupRouter(dsn)

	// Create a company
	company := entity.Company{
		Name:              generateRandomCompanyName(),
		Description:       "random description",
		AmountOfEmployees: 6,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)
	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Update the company
	var createdCompany entity.Company
	err := json.Unmarshal(w.Body.Bytes(), &createdCompany)
	if err != nil {
		log.Println(err)
		return
	}
	createdCompany.Description = "updated description"
	createdCompany.AmountOfEmployees = 10
	jsonValue, _ = json.Marshal(createdCompany)

	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/api/companies/%s", createdCompany.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCompanyHandler_GetCompany(t *testing.T) {
	router := setupRouter(dsn)

	// Create a company
	company := entity.Company{
		Name:              "Get company test name",
		Description:       "random description",
		AmountOfEmployees: 7,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)
	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get the company
	var createdCompany entity.Company
	err := json.Unmarshal(w.Body.Bytes(), &createdCompany)
	if err != nil {
		log.Println(err)
		return
	}

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/companies/%s", createdCompany.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedCompany entity.Company
	err = json.Unmarshal(w.Body.Bytes(), &retrievedCompany)
	if err != nil {
		log.Println(err)
		return
	}
	assert.Equal(t, createdCompany.ID, retrievedCompany.ID)
	assert.Equal(t, createdCompany.Name, retrievedCompany.Name)
}

func TestCompanyHandler_GetCompanyNotFoundAndInvalid(t *testing.T) {
	router := setupRouter(dsn)

	token := getToken()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/companies/%s", "safasfa-asfas-afasfas-fas"), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/companies/%s", "7b0f69f0-ed63-4d5b-a4d0-45908f698666"), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestCompanyHandler_DeleteCompany(t *testing.T) {
	router := setupRouter(dsn)

	// Create a company
	company := entity.Company{
		Name:              generateRandomCompanyName(),
		Description:       "random description",
		AmountOfEmployees: 8,
		Registered:        true,
		Type:              entity.Corporation,
	}
	jsonValue, _ := json.Marshal(company)
	token := getToken()

	req, _ := http.NewRequest("POST", "/api/companies/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Delete the company
	var createdCompany entity.Company
	err := json.Unmarshal(w.Body.Bytes(), &createdCompany)
	if err != nil {
		log.Println(err)
		return
	}

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/companies/%s", createdCompany.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func getToken() string {
	secretKey := []byte("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"sub": "test-user",
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return tokenString
}

func generateRandomCompanyName() string {
	gofakeit.Seed(0) // Seed to ensure randomness
	return gofakeit.Company()
}
