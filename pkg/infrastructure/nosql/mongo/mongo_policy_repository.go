package mongo

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	base "gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/errors"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nononsensecode.com/policy/pkg/components"
	"nononsensecode.com/policy/pkg/domain/model/policy"
)

const (
	MongoPolicyConnectError errors.Code = iota + 300
	MongoPolicyDisconnectError
	MongoPolicyInsertOneError
	MongoPolicyInsertUnmarshalError
	MongoPolicyInvalidIdError
	MongoPolicyFindByNameNotFoundError
	MongoPolicyFindByNameMarshalByteError
	MongoPolicyFindByNameUnmarshalMongoPolicyError
	MongoPolicyFindByNameUnmarshalPolicyError
)

type MongoPolicyRepository struct {
	builder mongodb.ClientProvider
}

func NewPolicyRepository(b base.MongoDbClientBuilder) MongoPolicyRepository {
	if b == nil {
		panic("mongo client provider is nil")
	}

	return MongoPolicyRepository{
		builder: b,
	}
}

func (m MongoPolicyRepository) Save(ctx context.Context, p policy.Policy) (saved policy.Policy, err error) {
	var (
		client *mongo.Client
		res    *mongo.InsertOneResult
		id     primitive.ObjectID
		ok     bool
	)

	log.Debug("retrieving client")
	client, err = m.builder.GetMongoDbClient(ctx)
	if err != nil {
		log.Errorf("cannot retrieve client: %v", err)
		err = errors.NewUnknownError(MongoPolicyConnectError, components.Policy, err)
		return
	}
	defer func() {
		log.Debug("disconnecting client")
		disConnErr := client.Disconnect(ctx)
		if disConnErr != nil {
			err = errors.NewUnknownError(MongoPolicyDisconnectError, components.Policy, disConnErr)
		}
	}()

	mp := new(mongoPolicy)
	if err = mp.Marshal(p); err != nil {
		err = errors.NewUnknownError(MongoPolicyInsertUnmarshalError, components.Policy, err)
		return
	}

	log.Debugf("policy is %+v", mp)

	coll := client.Database("dhp_policy").Collection("policies")
	log.Info("Saving policy to disk")
	if res, err = coll.InsertOne(ctx, mp); err != nil {
		log.Errorf("inserting policy to disk failed: %v", err)
		err = errors.NewUnknownError(MongoPolicyInsertOneError, components.Policy, err)
		return
	}

	if id, ok = res.InsertedID.(primitive.ObjectID); !ok {
		idErr := fmt.Errorf("invalid id type: %v", res.InsertedID)
		err = errors.NewUnknownError(MongoPolicyInvalidIdError, components.Policy, idErr)
		return
	}
	log.Debugf("Inserted id is %v", res.InsertedID)

	saved, err = policy.UnmarshalFromPersistence(id.Hex(), p.Name(), p.Version(), p.Statements()...)
	return
}

func (m MongoPolicyRepository) FindByName(ctx context.Context, name string) (p policy.Policy, err error) {
	var (
		client *mongo.Client
	)

	client, err = m.builder.GetMongoDbClient(ctx)
	if err != nil {
		log.Errorf("cannot retrieve client: %v", err)
		err = errors.NewUnknownError(MongoPolicyConnectError, components.Policy, err)
		return
	}

	coll := client.Database("dhp_policy").Collection("policies")
	log.Debugf("Finding policy named %s", name)
	filter := bson.D{{Key: "name", Value: name}}
	var r mongoPolicyD
	if err = coll.FindOne(ctx, filter).Decode(&r); err != nil {
		log.Error("cannot find the document: %+v", err)
		err = errors.NewNotFoundError(MongoPolicyFindByNameNotFoundError, components.Policy, err)
		return
	}

	var mp mongoPolicy
	if mp, err = r.Unmarshal(); err != nil {
		return
	}

	p, err = mp.Unmarshal()
	return
}
