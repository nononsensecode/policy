package policy

import "context"

type Repository interface {
	Save(context.Context, Policy) (Policy, error)
	FindByName(context.Context, string) (Policy, error)
}
