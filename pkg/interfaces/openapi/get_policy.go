package openapi

import (
	"net/http"

	"github.com/go-chi/render"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/errors"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/interfaces/httpsrvr/httperr"
	"nononsensecode.com/policy/pkg/components"
	"nononsensecode.com/policy/pkg/domain/model/policy"
)

const (
	OpenapiGetByNameRenderError errors.Code = iota + 30
)

func (s Server) GetPolicy(w http.ResponseWriter, r *http.Request, policyName string) {
	p, err := s.policyService.FindByName(r.Context(), policyName)
	if err != nil {
		httperr.RespondWithApiError(err, w, r)
		return
	}

	var pr = new(Policy)
	pr.Marshal(p)

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, pr); err != nil {
		err = errors.NewUnknownError(OpenapiGetByNameRenderError, components.Policy, err)
		httperr.RespondWithApiError(err, w, r)
		return
	}
}

func (p *Policy) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (pr *Policy) Marshal(p policy.Policy) {
	pr.Name = p.Name()
	pr.Version = p.Version()

	var statements []Statement
	for _, s := range p.Statements() {
		sid := s.Sid()
		effect := s.Effect()
		resource := s.Resource()
		actions := s.Actions()
		like := s.Filter().Like()
		equals := s.Filter().Equals()

		stmt := Statement{
			Sid:      &sid,
			Effect:   &effect,
			Resource: &resource,
			Actions:  &actions,
			Filter: &Filter{
				Like:   &like,
				Equals: &equals,
			},
		}
		statements = append(statements, stmt)
	}

	pr.Statements = statements
}
