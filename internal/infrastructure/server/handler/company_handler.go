package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/innoglobe/xmgo/internal/entity"
	"github.com/innoglobe/xmgo/internal/middleware"
	"github.com/innoglobe/xmgo/internal/usecase"
	"net/http"
)

// CompanyHandler is a struct that contains the usecase for company
type CompanyHandler struct {
	companyUsecase usecase.CompanyUsecaseInterface
}

// NewCompanyHandler is a function that returns a new CompanyHandler
func NewCompanyHandler(companyUsecase usecase.CompanyUsecaseInterface) *CompanyHandler {
	return &CompanyHandler{
		companyUsecase: companyUsecase,
	}
}

// RegisterRoutes is a function that registers the routes for the company handler
func (h *CompanyHandler) RegisterRoutes(r *gin.RouterGroup, secretKey string) {
	companyRoutes := r.Group("/companies")
	companyRoutes.Use(middleware.JWTAuthMiddleware(secretKey))
	{
		companyRoutes.POST("/", h.CreateCompany)
		companyRoutes.PATCH("/:id", h.PatchCompany)   // PATCH /companies/:id
		companyRoutes.DELETE("/:id", h.DeleteCompany) // DELETE /companies/:id
		companyRoutes.GET("/:id", h.GetCompany)       // GET /companies/:id
	}
}

// CreateCompany godoc
// @Summary Create a new company
// @Description Create a new company with the provided details
// @Tags companies
// @Accept json
// @Produce json
// @Param company body entity.Company true "Company details"
// @Success 201 {object} entity.Company
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req entity.Company

	// Validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	res, err := h.companyUsecase.CreateCompany(c.Request.Context(), &req)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GetCompany godoc
// @Summary Get a company by ID
// @Description Get details of a company by its ID
// @Tags companies
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} entity.Company
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	company, err := h.companyUsecase.GetCompany(c.Request.Context(), cid)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

// PatchCompany godoc
// @Summary Update an existing company
// @Description Update the details of an existing company by its ID
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company body entity.Company true "Updated company details"
// @Success 200 {object} entity.Company
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/companies/{id} [patch]
func (h *CompanyHandler) PatchCompany(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var req entity.Company
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	res, err := h.companyUsecase.UpdateCompany(c.Request.Context(), cid, &req)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteCompany godoc
// @Summary Delete a company by ID
// @Description Delete a company by its ID
// @Tags companies
// @Produce json
// @Param id path string true "Company ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = h.companyUsecase.DeleteCompany(c.Request.Context(), cid)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Company deleted successfully"})
}

func (h *CompanyHandler) handleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]map[string]string, len(ve))
		for i, fe := range ve {
			out[i] = map[string]string{
				"field":   fe.Field(),
				"message": fe.Error(),
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": out,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
}

func (h *CompanyHandler) getStatusCode(err interface{}) int {
	if customErr, ok := err.(interface{ StatusCode() int }); ok {
		return customErr.StatusCode()
	}

	return http.StatusInternalServerError
}
