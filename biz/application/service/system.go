package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/convertor"
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	"github.com/CloudStriver/go-pkg/utils/pagination/mongop"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type SystemService interface {
	ReadNotification(ctx context.Context, req *gensystem.ReadNotificationReq) (resp *gensystem.ReadNotificationResp, err error)
	GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error)
	CreateNotification(ctx context.Context, req *gensystem.CreateNotificationReq) (resp *gensystem.CreateNotificationResp, err error)
	ReadNotifications(ctx context.Context, req *gensystem.ReadNotificationsReq) (resp *gensystem.ReadNotificationsResp, err error)
	CleanNotification(ctx context.Context, req *gensystem.CleanNotificationReq) (resp *gensystem.CleanNotificationResp, err error)
}

type SystemServiceImpl struct {
	MongoMapper notificationmapper.INotificationMongoMapper
}

func (s *SystemServiceImpl) ReadNotification(ctx context.Context, req *gensystem.ReadNotificationReq) (resp *gensystem.ReadNotificationResp, err error) {
	resp = new(gensystem.ReadNotificationResp)
	if err = s.MongoMapper.ReadNotification(ctx, req.NotificationId); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error) {
	resp = new(gensystem.GetNotificationsResp)
	p := pconvertor.PaginationOptionsToModelPaginationOptions(req.PaginationOptions)
	notifications, total, err := s.MongoMapper.GetNotificationsAndCount(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions), p, mongop.IdCursorType)
	if err != nil {
		return resp, err
	}
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Notifications = lo.Map[*notificationmapper.Notification, *gensystem.Notification](notifications,
		func(item *notificationmapper.Notification, index int) *gensystem.Notification {
			return convertor.NotificationMapperToNotification(item)
		})
	resp.Total = total
	return resp, nil
}

func (s *SystemServiceImpl) CleanNotification(ctx context.Context, req *gensystem.CleanNotificationReq) (resp *gensystem.CleanNotificationResp, err error) {
	resp = new(gensystem.CleanNotificationResp)
	if err = s.MongoMapper.CleanNotification(ctx, req.UserId); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error) {
	resp = new(gensystem.GetNotificationCountResp)
	if resp.Total, err = s.MongoMapper.Count(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateNotification(ctx context.Context, req *gensystem.CreateNotificationReq) (resp *gensystem.CreateNotificationResp, err error) {
	resp = new(gensystem.CreateNotificationResp)
	if err = s.MongoMapper.InsertOne(ctx, convertor.NotificationToNotificationMapper(req.Notification)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) ReadNotifications(ctx context.Context, req *gensystem.ReadNotificationsReq) (resp *gensystem.ReadNotificationsResp, err error) {
	resp = new(gensystem.ReadNotificationsResp)
	if err = s.MongoMapper.ReadNotifications(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions)); err != nil {
		return resp, err
	}
	return resp, nil
}

var SystemSet = wire.NewSet(
	wire.Struct(new(SystemServiceImpl), "*"),
	wire.Bind(new(SystemService), new(*SystemServiceImpl)),
)
