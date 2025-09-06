package billing

import (
	"database/sql"
	"time"
)

type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
	Interval    string `json:"interval"`
	CreatedAt   string `json:"created_at"`
}

type Subscription struct {
	ID                string    `json:"id"`
	OrgID            string    `json:"org_id"`
	PlanID           string    `json:"plan_id"`
	Status           string    `json:"status"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`
	CreatedAt        string    `json:"created_at"`
}

type Invoice struct {
	ID             string     `json:"id"`
	SubscriptionID string     `json:"subscription_id"`
	AmountCents    int        `json:"amount_cents"`
	Status         string     `json:"status"`
	DueDate        time.Time  `json:"due_date"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedAt      string     `json:"created_at"`
}

type BillingService struct {
	db *sql.DB
}

func NewBillingService(db *sql.DB) *BillingService {
	return &BillingService{db: db}
}

func (s *BillingService) CreatePlan(name, description string, priceCents int, interval string) (*Plan, error) {
	var plan Plan
	err := s.db.QueryRow(`
		INSERT INTO plans (name, description, price_cents, interval)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price_cents, interval, created_at
	`, name, description, priceCents, interval).Scan(
		&plan.ID, &plan.Name, &plan.Description,
		&plan.PriceCents, &plan.Interval, &plan.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *BillingService) GetPlans() ([]Plan, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, price_cents, interval, created_at
		FROM plans
		ORDER BY price_cents ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var plan Plan
		if err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description,
			&plan.PriceCents, &plan.Interval, &plan.CreatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func (s *BillingService) CreateSubscription(orgID, planID string) (*Subscription, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get plan details
	var interval string
	var priceCents int
	err = tx.QueryRow(`
		SELECT interval, price_cents FROM plans WHERE id = $1
	`, planID).Scan(&interval, &priceCents)

	if err != nil {
		return nil, err
	}

	// Calculate period end based on interval
	periodEnd := time.Now()
	if interval == "month" {
		periodEnd = periodEnd.AddDate(0, 1, 0)
	} else {
		periodEnd = periodEnd.AddDate(1, 0, 0)
	}

	var sub Subscription
	err = tx.QueryRow(`
		INSERT INTO subscriptions (org_id, plan_id, status, current_period_start, current_period_end)
		VALUES ($1, $2, 'active', NOW(), $3)
		RETURNING id, org_id, plan_id, status, current_period_end, created_at
	`, orgID, planID, periodEnd).Scan(
		&sub.ID, &sub.OrgID, &sub.PlanID,
		&sub.Status, &sub.CurrentPeriodEnd, &sub.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Create first invoice
	_, err = tx.Exec(`
		INSERT INTO invoices (subscription_id, amount_cents, status, due_date)
		VALUES ($1, $2, 'unpaid', NOW())
	`, sub.ID, priceCents)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *BillingService) GetOrgSubscription(orgID string) (*Subscription, error) {
	var sub Subscription
	err := s.db.QueryRow(`
		SELECT id, org_id, plan_id, status, current_period_end, created_at
		FROM subscriptions
		WHERE org_id = $1 AND status = 'active'
	`, orgID).Scan(
		&sub.ID, &sub.OrgID, &sub.PlanID,
		&sub.Status, &sub.CurrentPeriodEnd, &sub.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *BillingService) GetInvoices(orgID string) ([]Invoice, error) {
	rows, err := s.db.Query(`
		SELECT i.id, i.subscription_id, i.amount_cents, i.status,
			   i.due_date, i.paid_at, i.created_at
		FROM invoices i
		JOIN subscriptions s ON s.id = i.subscription_id
		WHERE s.org_id = $1
		ORDER BY i.created_at DESC
	`, orgID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(
			&inv.ID, &inv.SubscriptionID, &inv.AmountCents,
			&inv.Status, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}

	return invoices, nil
}
