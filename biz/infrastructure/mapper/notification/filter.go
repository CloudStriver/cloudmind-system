package notification

import (
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"go.mongodb.org/mongo-driver/bson"
)

type FilterOptions struct {
	OnlyUserId  *string
	OnlyUserIds []string
	OnlyType    *int64
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
	f.CheckOnlyUserIds()
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

func (f *MongoFilter) CheckOnlyUserIds() {
	if f.OnlyUserIds != nil {
		f.m[consts.TargetUserId] = bson.M{
			"$in": f.OnlyUserIds,
		}
	}
}
