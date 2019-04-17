package plans

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func planFromRows(rows *sql.Rows) (*Plan, error) {
	var metadata Metadata
	var planType, planMetadata string

	plan := &Plan{}
	err := rows.Scan(&plan.ID, &plan.Name, &planType, &plan.Price,
		&planMetadata, &plan.Public, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	plan.Type = Type(planType)
	err = json.Unmarshal([]byte(planMetadata), &metadata)
	if err != nil {
		return nil, err
	}
	plan.Metadata = metadata
	return plan, nil
}

func CreatePlan(ctx context.Context, name string, planType Type, price float32, metadata Metadata, public bool) (*Plan, error) {
	strMet, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	plan := Plan{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      planType,
		Price:     price,
		Metadata:  metadata,
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

func DeletePlanByID(ctx context.Context, id string) error {
	_, e := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM plans WHERE id = $1
	`, id)
	return e
}

func FindPlanByID(ctx context.Context, id string) (*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, type, price, metadata, public, created_at, updated_at
		FROM plans WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return planFromRows(rows)
	}
	return nil, nil
}

func FindPlansByType(ctx context.Context, planType string) ([]*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, type, price, metadata, public, created_at, updated_at
		FROM plans WHERE type = $1
	`, planType)
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
