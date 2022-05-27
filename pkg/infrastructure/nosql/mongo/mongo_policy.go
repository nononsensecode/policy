package mongo

import (
	log "github.com/sirupsen/logrus"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nononsensecode.com/policy/pkg/components"
	"nononsensecode.com/policy/pkg/domain/model/policy"
)

type mongoPolicy struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Version    string             `bson:"version,omitempty"`
	Statements []mongoStatement   `bson:"statements,omitempty"`
}

type mongoStatement struct {
	Sid      string      `bson:"sid,omitempty"`
	Effect   string      `bson:"effect,omitempty"`
	Resource string      `bson:"resource,omitempty"`
	Actions  []string    `bson:"actions,omitempty"`
	Filter   mongoFilter `bson:"filter,omitempty"`
}

type mongoFilter struct {
	Like   []string `bson:"like,omitempty"`
	Equals []string `bson:"equals,omitempty"`
}

func (m *mongoPolicy) Marshal(p policy.Policy) (err error) {
	if p.ID() != "" {
		m.Id, err = primitive.ObjectIDFromHex(p.ID())
		if err != nil {
			return
		}
	}

	m.Name = p.Name()
	m.Version = p.Version()
	for _, s := range p.Statements() {
		statement := mongoStatement{}
		statement.Sid = s.Sid()
		statement.Effect = s.Effect()
		statement.Resource = s.Resource()
		statement.Actions = s.Actions()
		statement.Filter.Like = s.Filter().Like()
		statement.Filter.Equals = s.Filter().Equals()
		m.Statements = append(m.Statements, statement)
	}
	return
}

func (m mongoPolicy) Unmarshal() (p policy.Policy, err error) {
	var statements []policy.Statement
	for _, s := range m.Statements {
		var statement policy.Statement
		statement, err = policy.NewStatement(s.Sid, s.Effect, s.Resource, s.Actions, s.Filter.Like, s.Filter.Equals)
		if err != nil {
			return
		}
		statements = append(statements, statement)
	}

	p, err = policy.UnmarshalFromPersistence(m.Id.Hex(), m.Name, m.Version, statements...)
	return
}

type mongoPolicyD bson.D

func (m mongoPolicyD) Unmarshal() (mp mongoPolicy, err error) {
	var b []byte
	b, err = bson.Marshal(m)
	if err != nil {
		log.Errorf("cannot marshal bson.D: %v", err)
		err = errors.NewUnknownError(MongoPolicyFindByNameMarshalByteError, components.Policy, err)
		return
	}

	if err = bson.Unmarshal(b, &mp); err != nil {
		log.Errorf("cannot unmarshal to mongopolicy: %v", err)
		err = errors.NewUnknownError(MongoPolicyFindByNameUnmarshalMongoPolicyError, components.Policy, err)
		return
	}
	return
}
