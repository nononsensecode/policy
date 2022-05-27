package application

import (
	"context"

	log "github.com/sirupsen/logrus"
	"nononsensecode.com/policy/pkg/domain/model/policy"
)

type PolicyService struct {
	policyRepo policy.Repository
}

func NewPolicyService(p policy.Repository) PolicyService {
	if p == nil {
		panic("policy repository is nil")
	}

	return PolicyService{p}
}

func (ps PolicyService) SavePolicy(ctx context.Context, p policy.Policy) (pid string, err error) {
	log.Infof("saving policy \"%s\"", p.Name())
	p, err = ps.policyRepo.Save(ctx, p)
	if err != nil {
		log.Errorf("Saving policy failed: %v", err)
		return
	}

	log.Infof("Policy saved successfully. Id: %s", p.ID())
	pid = p.ID()
	return
}

func (ps PolicyService) FindByName(ctx context.Context, name string) (p policy.Policy, err error) {
	log.Infof("Finding policy with name %s", name)
	p, err = ps.policyRepo.FindByName(ctx, name)
	return
}
