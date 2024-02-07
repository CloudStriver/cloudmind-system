package notification

import (
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterOptions struct {
	OnlyUserId          *string
	OnlyType            *int64
	OnlyIsRead          *bool
	OnlyNotificationIds []string
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
	f.CheckOnlyIsRead()
	f.CheckOnlyNotificationIds()
	return f.m
}

func (f *MongoFilter) CheckOnlyNotificationIds() {
	if f.OnlyNotificationIds != nil {
		f.m[consts.ID] = bson.M{"$in": lo.Map[string, primitive.ObjectID](f.OnlyNotificationIds, func(item string, index int) primitive.ObjectID {
			oid, _ := primitive.ObjectIDFromHex(item)
			return oid
		}),
		}
	}
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

func (f *MongoFilter) CheckOnlyIsRead() {
	if f.OnlyIsRead != nil {
		f.m[consts.IsRead] = bson.M{"$exists": *f.OnlyIsRead}
	}
}
