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
	OpenApiSavePolicyIncorrectInputError errors.Code = iota + 400
	OpenApiSavePolicyRenderError
)

func (s Server) SavePolicy(w http.ResponseWriter, r *http.Request) {
	policyRequest := new(SavePolicyJSONRequestBody)
	if err := render.Bind(r, policyRequest); err != nil {
		err = errors.NewIncorrectInputError(OpenApiSavePolicyIncorrectInputError, components.Policy, err)
		httperr.RespondWithApiError(err, w, r)
		return
	}

	p, err := policyRequest.Unmarshal()
	if err != nil {
		err = errors.NewIncorrectInputError(OpenApiSavePolicyIncorrectInputError, components.Policy, err)
		httperr.RespondWithApiError(err, w, r)
		return
	}

	pid, err := s.policyService.SavePolicy(r.Context(), p)
	if err != nil {
		httperr.RespondWithApiError(err, w, r)
		return
	}

	pidRes := PolicyId{Id: &pid}
	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, pidRes); err != nil {
		err = errors.NewUnknownError(OpenApiSavePolicyRenderError, components.Policy, err)
		httperr.RespondWithApiError(err, w, r)
		return
	}
}

func (p *SavePolicyJSONRequestBody) Bind(r *http.Request) error {
	return nil
}

func (pr *SavePolicyJSONRequestBody) Unmarshal() (p policy.Policy, err error) {
	var statements []policy.Statement
	for _, s := range pr.Statements {
		var statement policy.Statement
		statement, err = policy.NewStatement(*s.Sid, *s.Effect, *s.Resource, *s.Actions, *s.Filter.Like, *s.Filter.Equals)
		if err != nil {
			return
		}
		statements = append(statements, statement)
	}

	p, err = policy.New(pr.Name, pr.Version, statements...)
	return
}

func (p PolicyId) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Statement) Marshal(ps policy.Statement) {
	sid := ps.Sid()
	effect := ps.Effect()
	resource := ps.Resource()
	actions := ps.Actions()
	equals := ps.Filter().Equals()
	like := ps.Filter().Like()

	s.Sid = &sid
	s.Effect = &effect
	s.Resource = &resource
	s.Actions = &actions
	s.Filter.Equals = &equals
	s.Filter.Like = &like
}
