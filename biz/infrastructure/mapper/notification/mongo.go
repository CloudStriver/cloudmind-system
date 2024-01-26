package notification

import (
	"context"
	"github.com/CloudStriver/go-pkg/utils/pagination"
	"github.com/CloudStriver/go-pkg/utils/pagination/mongop"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/zeromicro/go-zero/core/stores/monc"

	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
)

const (
	CollectionName = "notification"
)

const prefixNotificationCacheKey = "cache:notification:"

var _ INotificationMongoMapper = (*MongoMapper)(nil)

type (
	INotificationMongoMapper interface {
		GetNotifications(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Notification, int64, error)
		CleanNotification(ctx context.Context, userId string) error
		ReadNotification(ctx context.Context, id string) error
		Count(ctx context.Context, fopts *FilterOptions) (int64, error)
		ReadNotifications(ctx context.Context, fopts *FilterOptions) error
		InsertOne(ctx context.Context, data *Notification) error
		GetNotificationsAndCount(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Notification, int64, error)
	}
	Notification struct {
		ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		TargetUserId    string             `bson:"targetUserId,omitempty" json:"targetUserId,omitempty"`
		SourceUserId    string             `bson:"sourceUserId,omitempty" json:"sourceUserId,omitempty"`
		SourceContentId string             `bson:"sourceContentId,omitempty" json:"sourceContentId,omitempty"`
		Type            int64              `bson:"type,omitempty" json:"type,omitempty"`
		TargetType      int64              `bson:"targetType,omitempty" json:"targetType,omitempty"`
		Text            string             `bson:"text,omitempty" json:"text,omitempty"`
		CreateAt        time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
		UpdateAt        time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
		IsRead          bool               `bson:"isRead,omitempty" json:"isRead,omitempty"`
	}
	MongoMapper struct {
		conn *monc.Model
	}
)

func (m *MongoMapper) GetNotificationsAndCount(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Notification, int64, error) {
	var (
		data       []*Notification
		count      int64
		err1, err2 error
	)
	p := mongop.NewMongoPaginator(pagination.NewRawStore(sorter), popts)

	filter := MakeBsonFilter(fopts)
	sort, err := p.MakeSortOptions(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if err = mr.Finish(func() error {
		count, err1 = m.conn.CountDocuments(ctx, filter)
		if err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		if err2 = m.conn.Find(ctx, &data, filter, &options.FindOptions{
			Sort:  sort,
			Limit: popts.Limit,
			Skip:  popts.Offset,
		}); err2 != nil {
			return err2
		}
		// 如果是反向查询，反转数据
		if *popts.Backward {
			lo.Reverse(data)
		}
		if len(data) > 0 {
			err2 = p.StoreCursor(ctx, data[0], data[len(data)-1])
			if err2 != nil {
				return err2
			}
		}
		return nil
	}); err != nil {
		return nil, 0, err
	}

	return data, count, nil
}
func (m *MongoMapper) InsertOne(ctx context.Context, data *Notification) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
	}
	data.CreateAt = time.Now()
	data.UpdateAt = time.Now()
	key := prefixNotificationCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}
func (m *MongoMapper) GetNotifications(ctx context.Context, fopts *FilterOptions, popts *pagination.PaginationOptions, sorter mongop.MongoCursor) ([]*Notification, int64, error) {
	var data []*Notification
	p := mongop.NewMongoPaginator(pagination.NewRawStore(sorter), popts)

	filter := MakeBsonFilter(fopts)
	sort, err := p.MakeSortOptions(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if err = m.conn.Find(ctx, &data, filter, &options.FindOptions{
		Sort:  sort,
		Limit: popts.Limit,
		Skip:  popts.Offset,
	}); err != nil {
		return nil, 0, err
	}

	// 如果是反向查询，反转数据
	if *popts.Backward {
		for i := 0; i < len(data)/2; i++ {
			data[i], data[len(data)-i-1] = data[len(data)-i-1], data[i]
		}
	}
	if len(data) > 0 {
		err = p.StoreCursor(ctx, data[0], data[len(data)-1])
		if err != nil {
			return nil, 0, err
		}
	}

	count, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (m *MongoMapper) ReadNotifications(ctx context.Context, fopts *FilterOptions) error {
	filter := MakeBsonFilter(fopts)
	if _, err := m.conn.UpdateManyNoCache(ctx, filter, bson.M{"$set": bson.M{consts.IsRead: true, consts.UpdateAt: time.Now()}}); err != nil {
		return err
	}
	return nil
}

// CleanNotification 清除未读消息
func (m *MongoMapper) CleanNotification(ctx context.Context, userId string) error {
	filter := bson.M{
		consts.TargetUserId: userId,
		consts.IsRead:       bson.M{"$exists": false},
	}
	_, err := m.conn.UpdateManyNoCache(ctx, filter, bson.M{"$set": bson.M{consts.IsRead: true, consts.UpdateAt: time.Now()}})
	return err
}

func (m *MongoMapper) ReadNotification(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}

	key := prefixNotificationCacheKey + id
	_, err = m.conn.UpdateByID(ctx, key, oid, bson.M{"$set": bson.M{consts.IsRead: true, consts.UpdateAt: time.Now()}})
	return err
}

func (m *MongoMapper) Count(ctx context.Context, fopts *FilterOptions) (int64, error) {
	f := MakeBsonFilter(fopts)
	return m.conn.CountDocuments(ctx, f)
}

func NewNotificationModel(config *config.Config) INotificationMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}
