package policy

import (
	"fmt"

	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/errors"
	"nononsensecode.com/policy/pkg/components"
)

const (
	PolicyStatementSidEmptyError errors.Code = iota + 11
	PolicyStatementEffectEmptyError
	PolicyStatementResourceEmptyError
	PolicyStatementActionsEmptyError
)

type Statement struct {
	sid      string
	effect   string
	actions  []string
	resource string
	filter   Filter
}

func (s Statement) Sid() string {
	return s.sid
}

func (s Statement) Effect() string {
	return s.effect
}

func (s Statement) Actions() []string {
	return s.actions
}

func (s Statement) Resource() string {
	return s.resource
}

func (s Statement) Filter() Filter {
	return s.filter
}

func NewStatement(sid, effect, resource string, actions, like, equals []string) (s Statement, err error) {
	if sid == "" {
		err = errors.NewIncorrectInputError(PolicyStatementSidEmptyError, components.Policy,
			fmt.Errorf("statement sid is empty"))
		return
	}

	if effect == "" {
		err = errors.NewIncorrectInputError(PolicyStatementEffectEmptyError, components.Policy,
			fmt.Errorf("statement effect is empty"))
		return
	}

	if resource == "" {
		err = errors.NewIncorrectInputError(PolicyStatementResourceEmptyError, components.Policy,
			fmt.Errorf("statement resource is empty"))
		return
	}

	if len(actions) == 0 {
		err = errors.NewIncorrectInputError(PolicyStatementActionsEmptyError, components.Policy,
			fmt.Errorf("statement actions are empty"))
		return
	}

	s = Statement{
		sid:      sid,
		effect:   effect,
		resource: resource,
		actions:  actions,
		filter: Filter{
			like:   like,
			equals: equals,
		},
	}

	return
}
