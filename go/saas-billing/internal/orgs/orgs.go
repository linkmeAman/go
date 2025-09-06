package orgs

import (
	"database/sql"
	"errors"
)

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
}

type Member struct {
	UserID    string `json:"user_id"`
	OrgID     string `json:"org_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type OrganizationService struct {
	db *sql.DB
}

func NewOrganizationService(db *sql.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

func (s *OrganizationService) Create(name string, ownerID string) (*Organization, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var org Organization
	err = tx.QueryRow(`
		INSERT INTO organizations (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`, name).Scan(&org.ID, &org.Name, &org.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Add owner as member
	_, err = tx.Exec(`
		INSERT INTO memberships (user_id, org_id, role)
		VALUES ($1, $2, $3)
	`, ownerID, org.ID, "owner")

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &org, nil
}

func (s *OrganizationService) GetUserOrgs(userID string) ([]Organization, error) {
	rows, err := s.db.Query(`
		SELECT o.id, o.name, o.created_at
		FROM organizations o
		JOIN memberships m ON m.org_id = o.id
		WHERE m.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (s *OrganizationService) AddMember(orgID, userID, role string) error {
	_, err := s.db.Exec(`
		INSERT INTO memberships (user_id, org_id, role)
		VALUES ($1, $2, $3)
	`, userID, orgID, role)

	return err
}

func (s *OrganizationService) CheckUserRole(userID, orgID string) (string, error) {
	var role string
	err := s.db.QueryRow(`
		SELECT role FROM memberships
		WHERE user_id = $1 AND org_id = $2
	`, userID, orgID).Scan(&role)

	if err == sql.ErrNoRows {
		return "", errors.New("user is not a member of this organization")
	}

	return role, err
}
