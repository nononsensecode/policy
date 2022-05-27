package openapi

import "nononsensecode.com/policy/pkg/application"

type Server struct {
	policyService application.PolicyService
}

func NewServer(ps application.PolicyService) Server {
	return Server{ps}
}
