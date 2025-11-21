package billing

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	portalsession "github.com/stripe/stripe-go/v78/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/subscription"
	"github.com/stripe/stripe-go/v78/webhook"
	"github.com/thebushidocollective/kasoku/server/internal/db"
	"github.com/thebushidocollective/kasoku/server/internal/models"
)

type Service struct {
	db                *db.DB
	webhookSecret     string
	individualPriceID string
	teamPriceID       string
}

func NewService(database *db.DB) *Service {
	// Initialize Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	return &Service{
		db:                database,
		webhookSecret:     os.Getenv("STRIPE_WEBHOOK_SECRET"),
		individualPriceID: os.Getenv("STRIPE_INDIVIDUAL_PRICE_ID"),
		teamPriceID:       os.Getenv("STRIPE_TEAM_PRICE_ID"),
	}
}

// CreateCheckoutSession creates a Stripe checkout session for a plan
func (s *Service) CreateCheckoutSession(userID uuid.UUID, email, plan string) (*stripe.CheckoutSession, error) {
	var priceID string
	switch plan {
	case "individual":
		priceID = s.individualPriceID
	case "team":
		priceID = s.teamPriceID
	default:
		return nil, fmt.Errorf("invalid plan: %s", plan)
	}

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(email),
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
		Metadata: map[string]string{
			"user_id": userID.String(),
			"plan":    plan,
		},
	}

	return checkoutsession.New(params)
}

// CreatePortalSession creates a Stripe billing portal session
func (s *Service) CreatePortalSession(customerID string) (string, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(os.Getenv("STRIPE_RETURN_URL")),
	}

	sess, err := portalsession.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}

// HandleWebhook processes Stripe webhook events
func (s *Service) HandleWebhook(payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(event)

	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event)

	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event)

	case "invoice.payment_succeeded":
		// Payment succeeded - subscription is active
		return nil

	case "invoice.payment_failed":
		return s.handlePaymentFailed(event)
	}

	return nil
}

func (s *Service) handleCheckoutCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return err
	}

	userIDStr := session.Metadata["user_id"]
	plan := session.Metadata["plan"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	// Get user
	user, err := s.db.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update user plan and store customer/subscription IDs
	user.Plan = plan
	user.StripeCustomerID = &session.Customer.ID
	user.StripeSubscriptionID = &session.Subscription.ID

	if err := s.db.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update storage quota based on plan
	var limitBytes int64
	switch plan {
	case "individual":
		limitBytes = 2 * 1024 * 1024 * 1024 // 2GB
	case "team":
		limitBytes = 50 * 1024 * 1024 * 1024 // 50GB
	default:
		limitBytes = 0
	}

	quota := &models.StorageQuota{
		UserID:     userID,
		UsedBytes:  0,
		LimitBytes: limitBytes,
	}

	if err := s.db.CreateStorageQuota(quota); err != nil {
		return fmt.Errorf("failed to create storage quota: %w", err)
	}

	return nil
}

func (s *Service) handleSubscriptionUpdated(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return err
	}

	// Get user by subscription ID
	_, err := s.db.GetUserByStripeSubscriptionID(sub.ID)
	if err != nil {
		return fmt.Errorf("user not found for subscription: %w", err)
	}

	// Update subscription status
	if sub.Status == stripe.SubscriptionStatusActive {
		// Subscription is active - ensure user has correct plan
		return nil
	}

	// Handle other statuses if needed (paused, past_due, etc.)
	return nil
}

func (s *Service) handleSubscriptionDeleted(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return err
	}

	// Get user by subscription ID
	user, err := s.db.GetUserByStripeSubscriptionID(sub.ID)
	if err != nil {
		return fmt.Errorf("user not found for subscription: %w", err)
	}

	// Downgrade to free plan
	user.Plan = "free"
	user.StripeSubscriptionID = nil

	if err := s.db.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update storage quota to 0
	quota, err := s.db.GetStorageQuota(user.ID)
	if err == nil && quota != nil {
		quota.LimitBytes = 0
		s.db.UpdateStorageQuota(quota)
	}

	return nil
}

func (s *Service) handlePaymentFailed(event stripe.Event) error {
	// Handle failed payment - could send notification email, etc.
	return nil
}

// CancelSubscription cancels a user's subscription
func (s *Service) CancelSubscription(subscriptionID string) error {
	params := &stripe.SubscriptionCancelParams{}
	_, err := subscription.Cancel(subscriptionID, params)
	return err
}

// GetSubscription retrieves subscription details
func (s *Service) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}
