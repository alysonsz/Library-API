package web

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project-go/internal/services"
	"strconv"
	"time"
)

type LoanHandlers struct {
	service *services.LoanService
}

func NewLoanHandlers(service *services.LoanService) *LoanHandlers {
	return &LoanHandlers{service: service}
}

// GetLoans godoc
// @Summary      List loans
// @Description  Get all loans
// @Tags         loans
// @Accept       json
// @Produce      json
// @Success      200  {array}   services.Loan
// @Failure      500  {object}  map[string]interface{}
// @Router       /loans [get]
func (h *LoanHandlers) GetLoans(w http.ResponseWriter, r *http.Request) {
	loans, err := h.service.GetLoans()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve loans")
		return
	}

	RespondWithJSON(w, http.StatusOK, loans)
}

// CreateLoan godoc
// @Summary      Create a loan
// @Description  Register a new book loan
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        loan  body      object  true  "Loan data (bookId, borrowerName, dueDate RFC3339)"
// @Success      201   {object}  services.Loan
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /loans [post]
func (h *LoanHandlers) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookID       int    `json:"bookId"`
		BorrowerName string `json:"borrowerName"`
		DueDate      string `json:"dueDate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_date", "Invalid due date format (use RFC3339)")
		return
	}

	loan := services.Loan{
		BookID:       req.BookID,
		BorrowerName: req.BorrowerName,
		DueDate:      dueDate,
		Status:       "active",
	}

	if ve := ValidateLoan(&loan); !ve.IsEmpty() {
		RespondWithJSON(w, http.StatusBadRequest, ve)
		return
	}

	if err := h.service.CreateLoan(&loan); err != nil {
		if err.Error() == "book "+strconv.Itoa(req.BookID)+" is already loaned" {
			RespondWithError(w, http.StatusConflict, "already_loaned", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to create loan")
		return
	}
	RespondWithJSON(w, http.StatusCreated, loan)
}

// GetLoanByID godoc
// @Summary      Get a loan
// @Description  Get loan by ID
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Loan ID"
// @Success      200  {object}  services.Loan
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /loans/{id} [get]
func (h *LoanHandlers) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid loan ID")
		return
	}

	loan, err := h.service.GetLoanByID(id)

	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "not_found", "Loan not found")
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve loan")
		return
	}
	RespondWithJSON(w, http.StatusOK, loan)
}

// ReturnLoan godoc
// @Summary      Return a loan
// @Description  Mark a loan as returned
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Loan ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /loans/{id}/return [post]
func (h *LoanHandlers) ReturnLoan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid loan ID")
		return
	}

	if err := h.service.ReturnLoan(id); err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "not_found", "Loan not found or already returned")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to return loan")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetOverdueLoans godoc
// @Summary      List overdue loans
// @Description  Get all overdue active loans
// @Tags         loans
// @Accept       json
// @Produce      json
// @Success      200  {array}   services.Loan
// @Failure      500  {object}  map[string]interface{}
// @Router       /loans/overdue [get]
func (h *LoanHandlers) GetOverdueLoans(w http.ResponseWriter, r *http.Request) {
	loans, err := h.service.GetOverdueLoans()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve overdue loans")
		return
	}
	RespondWithJSON(w, http.StatusOK, loans)
}
