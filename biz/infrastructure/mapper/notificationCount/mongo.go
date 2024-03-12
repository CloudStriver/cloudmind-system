package notification

import (
	"context"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionName       = "notificationCount"
	NotificationCountKey = "NotificationCount"
)

var _ INotificationCountMongoMapper = (*MongoMapper)(nil)

type (
	INotificationCountMongoMapper interface {
		AddCount(ctx context.Context, notification *NotificationCount) error
		GetCount(ctx context.Context, userId string) (*NotificationCount, error)
		InsertCount(ctx context.Context, notification *NotificationCount) error
	}
	NotificationCount struct {
		ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Sum  int64              `bson:"sum" json:"sum"`
		Read int64              `bson:"read" json:"read"`
	}
	MongoMapper struct {
		conn *monc.Model
	}
)

func (m MongoMapper) GetCount(ctx context.Context, userId string) (*NotificationCount, error) {
	var data *NotificationCount
	uid, _ := primitive.ObjectIDFromHex(userId)
	key := NotificationCountKey + userId
	if err := m.conn.FindOne(ctx, key, &data, bson.M{consts.ID: uid}); err != nil {
		return nil, err
	}
	return data, nil
}

func (m MongoMapper) AddCount(ctx context.Context, notification *NotificationCount) error {
	key := NotificationCountKey + notification.ID.Hex()
	if _, err := m.conn.UpdateOne(ctx, key, bson.M{consts.ID: notification.ID}, bson.M{"$inc": bson.M{
		consts.Sum:  notification.Sum,
		consts.Read: notification.Read,
	}}); err != nil {
		return err
	}
	return nil
}

func (m MongoMapper) InsertCount(ctx context.Context, notification *NotificationCount) error {
	key := NotificationCountKey + notification.ID.Hex()
	if _, err := m.conn.InsertOne(ctx, key, notification); err != nil {
		return err
	}
	return nil
}
func NewNotificationCountModel(config *config.Config) INotificationCountMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}
