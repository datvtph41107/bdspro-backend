package server

import (
	"context"
	"noti/services"
	pb_noti "pb/types/notification"
)

type ServerHttp struct {
	pb_noti.UnimplementedNotificationServiceServer
	Service  *services.NotificationService2
	Firebase *services.FirebaseService
}

func NewServerHttp(service *services.NotificationService2, firebase *services.FirebaseService) *ServerHttp {
	return &ServerHttp{
		Service:  service,
		Firebase: firebase,
	}
}

func (s *ServerHttp) SendBatch(c context.Context, dto *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error) {
	return s.Service.CreateBatch(c, dto)
	// return nil, nil
}
func (s *ServerHttp) Send(c context.Context, dto *pb_noti.NotiNewRequest) (*pb_noti.NotiNewResponse, error) {
	return s.Service.CreateNotification(c, dto)
}
func (s *ServerHttp) Count(c context.Context, dto *pb_noti.EmptyRequest) (*pb_noti.CountResponse, error) {
	total, err := s.Service.Count(c)
	return &pb_noti.CountResponse{
		Count: int32(total),
	}, err
}
func (s *ServerHttp) Read(c context.Context, dto *pb_noti.IdRequest) (*pb_noti.EmptyResponse, error) {
	_, err := s.Service.Read(c, uint64(dto.Id))
	return nil, err
}
func (s *ServerHttp) Remove(c context.Context, dto *pb_noti.IdRequest) (*pb_noti.EmptyResponse, error) {
	_, err := s.Service.Delete(c, dto)
	return nil, err
}
func (s *ServerHttp) List(c context.Context, dto *pb_noti.ListRequest) (*pb_noti.ListResponse, error) {
	return s.Service.Search(c, dto)
}
func (s *ServerHttp) PushByToken(c context.Context, dto *pb_noti.PushTokenRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}
func (s *ServerHttp) PushByTopic(c context.Context, dto *pb_noti.PushTopicRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}
