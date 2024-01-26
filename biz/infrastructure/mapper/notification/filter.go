package notification

import (
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterOptions struct {
	OnlyUserId     *string
	OnlyType       *int64
	OnlyTargetType *int64
	OnlyFirstId    *string
	OnlyLastId     *string
	OnlyIsRead     *bool
}

type MongoFilter struct {
	m bson.M
	*FilterOptions
}

func MakeBsonFilter(options *FilterOptions) bson.M {
	return (&MongoFilter{
		m:             bson.M{},
		FilterOptions: options,
	}).toBson()
}

func (f *MongoFilter) toBson() bson.M {
	f.CheckOnlyUserId()
	f.CheckOnlyType()
	f.CheckOnlyTargetType()
	f.CheckRange()
	f.CheckOnlyIsRead()
	return f.m
}

func (f *MongoFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.m[consts.TargetUserId] = *f.OnlyUserId
	}
}

func (f *MongoFilter) CheckOnlyType() {
	if f.OnlyType != nil {
		f.m[consts.Type] = *f.OnlyType
	}
}

func (f *MongoFilter) CheckOnlyTargetType() {
	if f.OnlyTargetType != nil {
		f.m[consts.TargetType] = *f.OnlyTargetType
	}
}

func (f *MongoFilter) CheckRange() {
	if f.OnlyLastId != nil && f.OnlyFirstId != nil {
		firstId, _ := primitive.ObjectIDFromHex(*f.OnlyFirstId)
		lastId, _ := primitive.ObjectIDFromHex(*f.OnlyLastId)
		f.m[consts.ID] = bson.M{"$gte": firstId, "$lte": lastId}
	}
}

func (f *MongoFilter) CheckOnlyIsRead() {
	if f.OnlyIsRead != nil {
		f.m[consts.IsRead] = bson.M{"$exists": *f.OnlyIsRead}
	}
}
