package request

import (
	"context"
	"github.com/docker/distribution/reference"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/plans"
	"github.com/google/uuid"
	"regexp"
)

var (
	// NameRegexp define the name constraint. Must start with
	// alpha numeric char, can contain _ or - chars.
	NameRegexp = regexp.MustCompile(`^[\w0-9]+[\w_\-0-9]$`)

	// TagRegexp define the tag constraint. Must start with alpha
	// char, can numericals character and contain . - _ characters.
	TagRegexp = regexp.MustCompile(`^[\w]+[\w0-9#?:.-_/]$`)

	// EnvironmentKeyRegexp define the shell posix standard to accept environment
	// key variable.
	EnvironmentKeyRegexp = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]$`)
)

type CreateContainer struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	PlanID      string            `json:"plan_id"`
	Tags        []string          `json:"tags"`
	Environment map[string]string `json:"environment"`
}

func (s *CreateContainer) Validate(ctx context.Context, namespaceId string) (errors map[string]string) {
	errors = make(map[string]string)

	if !NameRegexp.MatchString(s.Name) {
		errors["name"] = "Name must start with alphanumerical character and can contain dash or underscore."
	}
	if len(s.Name) < 3 || len(s.Name) > 32 {
		errors["name"] = "Name must be between 3 and 32 characters."
	}
	if ctr, err := containers.FindContainerByNameAndNamespaceID(ctx, s.Name, namespaceId); err != nil || ctr != nil {
		errors["name"] = "Name must be unique."
	}

	if !reference.ReferenceRegexp.MatchString(s.Image) {
		errors["image"] = "Invalid image name."
	}

	if len(s.Tags) > 100 {
		errors["tags"] = "You can't have more than 100 tags."
	} else {
		for _, tag := range s.Tags {
			if !TagRegexp.MatchString(tag) {
				errors["tags"] = "Tags must start with alpha characters, can contain numbers and this specials " +
					"characters set '#?:.-_/'."
			}
			if len(tag) == 0 || len(tag) > 253 {
				errors["tags"] = "Tags must be between 2 and 253 characters."
			}
		}
	}

	if len(s.Environment) > 100 {
		errors["environment"] = "You can't have more than 100 environments variables."
	} else {
		for key, value := range s.Environment {
			if !EnvironmentKeyRegexp.MatchString(key) {
				errors["environment"] = "Environment key must start with alpha characters or underscore and can " +
					"contain numericals values."
			}
			if len(key) == 0 || len(key) > 1024 {
				errors["environment"] = "Environment key must be between 1 and 1024 characters."
			}
			if len(value) > 32768 {
				errors["environment"] = "Environment value must be lesser than 32768 characters."
			}
		}
	}

	if _, err := uuid.Parse(s.PlanID); err != nil {
		errors["plan_id"] = "Invalid plan id."
	} else {
		plan, _ := plans.FindPlanByID(ctx, s.PlanID)
		if plan == nil {
			errors["plan_id"] = "Plan not found."
		}
	}
	return errors
}
