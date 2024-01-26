package slider

import (
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"go.mongodb.org/mongo-driver/bson"
)

type FilterOptions struct {
	OnlyType     *int64
	OnlyIsPublic *int64
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
	f.CheckOnlyType()
	f.CheckOnlyIsPublic()
	return f.m
}

func (f *MongoFilter) CheckOnlyType() {
	if f.OnlyType != nil {
		f.m[consts.Type] = *f.OnlyType
	}
}

func (f *MongoFilter) CheckOnlyIsPublic() {
	if f.OnlyIsPublic != nil {
		f.m[consts.IsPublic] = *f.OnlyIsPublic
	}
}
