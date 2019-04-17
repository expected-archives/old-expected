package plans

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func planFromRows(rows *sql.Rows) (*Plan, error) {
	var met Metadata
	var planType, planMetadata string

	plan := &Plan{}
	err := rows.Scan(&plan.ID, &plan.Name, &planType, &plan.Price,
		&planMetadata, &plan.Public, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	plan.Type = Type(planType)
	err = json.Unmarshal([]byte(planMetadata), &met)
	if err != nil {
		return nil, err
	}
	plan.Metadata = met
	return plan, nil
}

func CreatePlan(ctx context.Context, name string, planType Type, price float32, met Metadata, public bool) (*Plan, error) {
	strMet, err := json.Marshal(met)
	if err != nil {
		return nil, err
	}
	plan := Plan{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      planType,
		Price:     price,
		Metadata:  met,
		Public:    public,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO plans VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
	`, plan.ID, plan.Name, plan.Type, plan.Price, string(strMet), plan.Public, plan.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func DeletePlanByID(ctx context.Context, id string) error {
	_, e := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM plans WHERE id = $1
	`, id)
	return e
}

func UpdatePlan(ctx context.Context, plan Plan) error {
	strMet, err := json.Marshal(plan.Metadata)
	if err != nil {
		return err
	}
	_, e := services.Postgres().Client().ExecContext(ctx, `
		UPDATE plans 
		SET 
		    name = $2, 
		    type = $3, 
		    price = $4, 
		    metadata = $5, 
		    public = $6, 
		    updated_at = now() 
		WHERE id = $1
	`, plan.ID, plan.Name, plan.Type, plan.Price, string(strMet), plan.Public)
	return e
}

func FindPlanByID(ctx context.Context, id string) (*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT * FROM  plans WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return planFromRows(rows)
	}
	return nil, errors.New("plan not found")
}

func FindPlansByType(ctx context.Context, planType string, public bool) ([]*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT * FROM  plans WHERE type = $1 AND public = $2
	`, planType, public)
	if err != nil {
		return nil, err
	}
	var plans []*Plan
	for rows.Next() {
		plan, err := planFromRows(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func SetCustomPlan(ctx context.Context, namespaceId string, planId string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
			INSERT INTO custom_plans VALUES ($1, $2)
		`, planId, namespaceId)
	return err
}
