package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/convertor"
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
	"github.com/CloudStriver/go-pkg/utils/pagination/mongop"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type SystemService interface {
	GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error)
	CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error)
	ReadNotifications(ctx context.Context, req *gensystem.ReadNotificationsReq) (resp *gensystem.ReadNotificationsResp, err error)
	CleanNotification(ctx context.Context, req *gensystem.CleanNotificationReq) (resp *gensystem.CleanNotificationResp, err error)
	DeleteSlider(ctx context.Context, req *gensystem.DeleteSliderReq) (resp *gensystem.DeleteSliderResp, err error)
	UpdateSlider(ctx context.Context, req *gensystem.UpdateSliderReq) (resp *gensystem.UpdateSliderResp, err error)
	CreateSlider(ctx context.Context, req *gensystem.CreateSliderReq) (resp *gensystem.CreateSliderResp, err error)
	GetSliders(ctx context.Context, req *gensystem.GetSlidersReq) (resp *gensystem.GetSlidersResp, err error)
}

type SystemServiceImpl struct {
	NotificationMongoMapper notificationmapper.INotificationMongoMapper
	SliderMongoMapper       slidermapper.ISliderMongoMapper
}

func (s *SystemServiceImpl) DeleteSlider(ctx context.Context, req *gensystem.DeleteSliderReq) (resp *gensystem.DeleteSliderResp, err error) {
	resp = new(gensystem.DeleteSliderResp)
	if err = s.SliderMongoMapper.DeleteOne(ctx, req.SliderId); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) UpdateSlider(ctx context.Context, req *gensystem.UpdateSliderReq) (resp *gensystem.UpdateSliderResp, err error) {
	resp = new(gensystem.UpdateSliderResp)
	if err = s.SliderMongoMapper.UpdateOne(ctx, convertor.SliderToSliderMapper(req.Slider)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateSlider(ctx context.Context, req *gensystem.CreateSliderReq) (resp *gensystem.CreateSliderResp, err error) {
	resp = new(gensystem.CreateSliderResp)
	if err = s.SliderMongoMapper.InsertOne(ctx, convertor.SliderToSliderMapper(req.Slider)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetSliders(ctx context.Context, req *gensystem.GetSlidersReq) (resp *gensystem.GetSlidersResp, err error) {
	resp = new(gensystem.GetSlidersResp)
	p := pconvertor.PaginationOptionsToModelPaginationOptions(req.PaginationOptions)
	sliders, total, err := s.SliderMongoMapper.GetSlidersAndCount(ctx, convertor.SliderFilterOptionsToMapper(req.FilterOptions), p, mongop.IdCursorType)
	if err != nil {
		return resp, err
	}
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Sliders = lo.Map[*slidermapper.Slider, *gensystem.Slider](sliders,
		func(item *slidermapper.Slider, _ int) *gensystem.Slider {
			return convertor.SliderMapperToSlider(item)
		})
	resp.Total = total
	return resp, nil
}

func (s *SystemServiceImpl) GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error) {
	resp = new(gensystem.GetNotificationsResp)
	p := pconvertor.PaginationOptionsToModelPaginationOptions(req.PaginationOptions)
	notifications, total, err := s.NotificationMongoMapper.GetNotificationsAndCount(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions), p, mongop.IdCursorType)
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
	if err = s.NotificationMongoMapper.CleanNotification(ctx, req.UserId); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error) {
	resp = new(gensystem.GetNotificationCountResp)
	if resp.Total, err = s.NotificationMongoMapper.Count(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error) {
	resp = new(gensystem.CreateNotificationsResp)
	notifications := lo.Map[*gensystem.Notification, *notificationmapper.Notification](req.Notifications, func(item *gensystem.Notification, _ int) *notificationmapper.Notification {
		return convertor.NotificationToNotificationMapper(item)
	})
	if err = s.NotificationMongoMapper.InsertMany(ctx, notifications); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) ReadNotifications(ctx context.Context, req *gensystem.ReadNotificationsReq) (resp *gensystem.ReadNotificationsResp, err error) {
	resp = new(gensystem.ReadNotificationsResp)
	if err = s.NotificationMongoMapper.ReadNotifications(ctx, convertor.NotificationFilterOptionsToMapper(req.FilterOptions)); err != nil {
		return resp, err
	}
	return resp, nil
}

var SystemSet = wire.NewSet(
	wire.Struct(new(SystemServiceImpl), "*"),
	wire.Bind(new(SystemService), new(*SystemServiceImpl)),
)
