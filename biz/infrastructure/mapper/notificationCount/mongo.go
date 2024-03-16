package notification

import (
	"context"
	"errors"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
)

const (
	CollectionName       = "notificationCount"
	NotificationCountKey = "cache:NotificationCount:"
)

var _ INotificationCountMongoMapper = (*MongoMapper)(nil)

type (
	INotificationCountMongoMapper interface {
		GetNotificationCount(ctx context.Context, userId string) (int64, error)
		UpdateNotificationCount(ctx context.Context, data *NotificationCount) error
		CreateNotificationCount(ctx context.Context, data *NotificationCount) error
	}
	NotificationCount struct {
		ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Read int64              `bson:"read,omitempty" json:"read,omitempty"`
	}
	MongoMapper struct {
		conn *monc.Model
	}
)

func (m MongoMapper) CreateNotificationCount(ctx context.Context, data *NotificationCount) error {
	key := NotificationCountKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m MongoMapper) GetNotificationCount(ctx context.Context, userId string) (int64, error) {
	key := NotificationCountKey + userId
	var data *NotificationCount
	uid, _ := primitive.ObjectIDFromHex(userId)
	err := m.conn.FindOne(ctx, key, &data, bson.M{"_id": uid})
	switch {
	case errors.Is(err, monc.ErrNotFound):
		return 0, consts.ErrNotFound
	case err == nil:
		return data.Read, err
	default:
		return 0, err
	}
}

func (m MongoMapper) UpdateNotificationCount(ctx context.Context, data *NotificationCount) error {
	key := NotificationCountKey + data.ID.Hex()
	_, err := m.conn.UpdateOne(ctx, key, bson.M{consts.ID: data.ID}, bson.M{"$set": data})
	return err
}

func NewNotificationCountModel(config *config.Config) INotificationCountMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}
