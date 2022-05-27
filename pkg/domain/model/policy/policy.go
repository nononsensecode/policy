package policy

import (
	"fmt"

	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/errors"
	"nononsensecode.com/policy/pkg/components"
)

const (
	PolicyNameEmptyError errors.Code = iota + 1
	PolicyIdEmptyError
	PolicyVersionEmptyError
	PolicyStatementsEmptyError
)

type Policy struct {
	id         string
	name       string
	version    string
	statements []Statement
}

func (p Policy) ID() string {
	return p.id
}

func (p Policy) Name() string {
	return p.name
}

func (p Policy) Version() string {
	return p.version
}

func (p Policy) Statements() []Statement {
	return p.statements
}

func New(name, version string, statements ...Statement) (p Policy, err error) {
	if name == "" {
		err = errors.NewIncorrectInputError(PolicyNameEmptyError, components.Policy,
			fmt.Errorf("policy name is empty"))
		return
	}

	if version == "" {
		err = errors.NewIncorrectInputError(PolicyVersionEmptyError, components.Policy,
			fmt.Errorf("policy version is empty"))
	}

	if len(statements) == 0 {
		err = errors.NewIncorrectInputError(PolicyStatementsEmptyError, components.Policy,
			fmt.Errorf("policy statements are empty"))
		return
	}

	p = Policy{
		name:       name,
		version:    version,
		statements: statements,
	}
	return
}

func UnmarshalFromPersistence(id, name, version string, statements ...Statement) (p Policy, err error) {
	p, err = New(name, version, statements...)
	if id == "" {
		err = errors.NewUnknownError(PolicyStatementsEmptyError, components.Policy,
			fmt.Errorf("policy id is empty"))
	}
	p.id = id
	return
}
