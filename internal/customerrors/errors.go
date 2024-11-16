package customerrors

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

//type APIError struct {
//	Status int
//	Msg    string
//}
//
//func (e APIError) Error() string {
//	return e.Msg
//}

//type GenericError struct {
//	Status int
//	Msg    string
//}
//
//func (e GenericError) Error() string {
//	return e.Msg
//}

type DBConnectionError struct{}

func (e *DBConnectionError) Error() string {
	return "database connection error"
}

func (e *DBConnectionError) StatusCode() int {
	return http.StatusServiceUnavailable
}

type GenericTxError struct {
	Msg string
}

func (e GenericTxError) Error() string {
	return e.Msg
}

func (e GenericTxError) StatusCode() int {
	return http.StatusInternalServerError
}

type RecordNotFoundError struct {
	ID uuid.UUID
}

func (e RecordNotFoundError) Error() string {
	return fmt.Sprintf("Record with ID %s not found", e.ID)
}

func (e RecordNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

type CompanyExistsError struct {
	Name string
}

func (e CompanyExistsError) Error() string {
	return fmt.Sprintf("Company with name %s already exists", e.Name)
}

func (e CompanyExistsError) StatusCode() int {
	return http.StatusBadRequest
}

type IDUpdateError struct {
	ID uuid.UUID
}

func (e IDUpdateError) Error() string {
	return fmt.Sprintf("ID %s cannot be updated", e.ID)
}

func (e IDUpdateError) StatusCode() int {
	return http.StatusBadRequest
}

type InvalidIDError struct {
	ID uuid.UUID
}

func (e InvalidIDError) Error() string {
	return fmt.Sprintf("Invalid ID %s", e.ID)
}

func (e InvalidIDError) StatusCode() int {
	return http.StatusBadRequest
}
