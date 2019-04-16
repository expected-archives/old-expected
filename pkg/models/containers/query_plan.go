package containers

import (
	"context"
	"database/sql"
	"github.com/expectedsh/expected/pkg/services"
)

func planFromRows(rows *sql.Rows) (*Plan, error) {
	plan := &Plan{}
	if err := rows.Scan(&plan.ID, &plan.Name, &plan.Price, &plan.CPU, &plan.Memory,
		&plan.Available); err != nil {
		return nil, err
	}
	return plan, nil
}

func FindPlans(ctx context.Context) ([]*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, price, cpu, memory, available FROM container_plans
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func FindPlanByID(ctx context.Context, id string) (*Plan, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, price, cpu, memory, available FROM container_plans
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return planFromRows(rows)
	}

	return nil, nil
}
