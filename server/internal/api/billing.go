package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/billing"
)

type BillingHandler struct {
	billing *billing.Service
}

func NewBillingHandler(billingService *billing.Service) *BillingHandler {
	return &BillingHandler{
		billing: billingService,
	}
}

type CreateCheckoutRequest struct {
	Plan string `json:"plan"` // "individual" or "team"
}

// CreateCheckout handles POST /billing/checkout
func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Plan != "individual" && req.Plan != "team" {
		respondError(w, http.StatusBadRequest, "invalid plan")
		return
	}

	// Get user email from context or database
	// For simplicity, we'll use a placeholder - in reality you'd fetch the user
	email := "user@example.com" // TODO: Get from database

	session, err := h.billing.CreateCheckoutSession(userID, email, req.Plan)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"url": session.URL,
	})
}

// CreatePortal handles POST /billing/portal
func (h *BillingHandler) CreatePortal(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get user's Stripe customer ID
	// TODO: Fetch from database using userID
	customerID := "cus_example" // Placeholder

	url, err := h.billing.CreatePortalSession(customerID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"url": url,
	})
}

// HandleWebhook handles POST /billing/webhook
func (h *BillingHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		respondError(w, http.StatusBadRequest, "missing signature")
		return
	}

	if err := h.billing.HandleWebhook(payload, signature); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
