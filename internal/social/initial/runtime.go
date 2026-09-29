package initial

import (
	"social/infra/service"
	"social/internal/usecase"
)

type InitialApp struct {
	NewsFeedService *service.NewsFeedService
	ReportService   *service.ReportService
	CommentService  *service.CommentService
	LikeService     *service.LikeService
	NewsFeedUsecase *usecase.NewsFeedUsecase
}

func NewApp(newsFeedService *service.NewsFeedService,
	reportService *service.ReportService,
	commentService *service.CommentService,
	likeService *service.LikeService,
	newsFeedUsecase *usecase.NewsFeedUsecase,
) *InitialApp {
	return &InitialApp{
		NewsFeedService: newsFeedService,
		ReportService:   reportService,
		CommentService:  commentService,
		LikeService:     likeService,
		NewsFeedUsecase: newsFeedUsecase,
	}
}
