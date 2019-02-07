package models

import "time"

type ContainersModel struct{}

type Container struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Endpoint    string            `json:"endpoint"`
	Memory      int               `json:"memory"`
	Status      string            `json:"status"`
	Command     string            `json:"command"`
	Environment map[string]string `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (m *ContainersModel) GetAccountContainers(id string) ([]*Container, error) {
	rows, err := db.QueryContext(`
		SELECT id, name, image, endpoint, memory, created_at
		FROM containers WHERE account_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []*Container
	for rows.Next() {
		container := &Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.Image,
			&container.Endpoint, &container.Memory, &container.CreatedAt); err != nil {
			return nil, err
		}
	}
	return containers, nil
}
