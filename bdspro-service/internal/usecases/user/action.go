package user_usecase

// type ActionService struct {
// 	common.CrudUsecase[domain.Action]
// }

// func NewActionService(repo *ActionRepo) *ActionService {
// 	return &ActionService{
// 		BaseService: crud2.BaseService[domain.Action, *ActionRepo]{
// 			Repo: repo,
// 		},
// 	}
// }

// func (s *ActionService) CreateAction(c *gin.Context, dto ActionRequest) (*domain.Action, error) {
// 	entity := &domain.Action{
// 		ActionType: dto.Action,
// 		TargetID:   dto.TargetID,
// 	}

// 	err := _db.SaveWithAudit(c, &entity)

// 	return entity, err
// }

// func (s *ActionService) DashboardPost(c *gin.Context, id uint64) (*DashboardPostResponse, error) {
// 	result, err := s.Repo.DashboardPost(c, id)
// 	if err != nil {
// 		return nil, &_routes.Except{
// 			Code:    500,
// 			Message: err.Error(),
// 		}
// 	}

// 	return result, nil
// }

// func (s *ActionService) DashboardPostBar(c *gin.Context, dto PostBarRequest) (*DashboardPostResponse, error) {
// 	result, err := s.Repo.PostBar(c, dto)
// 	if err != nil {
// 		return nil, &_routes.Except{
// 			Code:    500,
// 			Message: err.Error(),
// 		}
// 	}

// 	return result, nil
// }
