package pending_data

import (
	"context"
	"fmt"

	_err "common/domain/err"
)

// IDispatcher chịu trách nhiệm thực thi lại lời gọi tới service tương ứng.
type IDispatcher interface {
	Dispatch(ctx context.Context, pending *PendingData) error
}

// IPendingDataUsecase định nghĩa usecase quản lý pending data.
type IPendingDataUsecase interface {
	// Enqueue lưu lại yêu cầu cần retry.
	Enqueue(ctx context.Context, serviceName string, port uint32, path string, data string) (*PendingData, error)
	// DequeueAll trả về toàn bộ pending của service để xử lý.
	DequeueAll(ctx context.Context, serviceName string) ([]PendingData, error)
	// Flush cố gắng gửi lại các pending tới service đích thông qua dispatcher.
	Flush(ctx context.Context, serviceName string, dispatcher IDispatcher) *ErrorCollection
}

type pendingDataUsecase struct {
	repo IPendingDataRepo
}

func NewPendingDataUsecase(repo IPendingDataRepo) IPendingDataUsecase {
	return &pendingDataUsecase{
		repo: repo,
	}
}

func (uc *pendingDataUsecase) Enqueue(ctx context.Context, serviceName string, port uint32, path string, data string) (*PendingData, error) {
	entity := &PendingData{
		ServiceName: serviceName,
		Port:        port,
		Path:        path,
		Data:        data,
	}

	if err := uc.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: fmt.Sprintf("không thể lưu pending data: %v", err),
		}
	}

	return entity, nil
}

func (uc *pendingDataUsecase) DequeueAll(ctx context.Context, serviceName string) ([]PendingData, error) {
	result, err := uc.repo.GetByServiceName(ctx, serviceName, 0)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: fmt.Sprintf("không thể lấy pending data của service %s: %v", serviceName, err),
		}
	}
	return result, nil
}

func (uc *pendingDataUsecase) Flush(ctx context.Context, serviceName string, dispatcher IDispatcher) *ErrorCollection {
	pendingList, err := uc.repo.GetByServiceName(ctx, serviceName, 0)
	if err != nil {
		collection := &ErrorCollection{}
		collection.Append(fmt.Errorf("không thể lấy pending data: %w", err))
		return collection
	}

	var processedIDs []uint64
	errCollection := &ErrorCollection{}

	for i := range pendingList {
		item := pendingList[i]
		if dispatchErr := dispatcher.Dispatch(ctx, &item); dispatchErr != nil {
			errCollection.Append(fmt.Errorf("dispatch pending id=%d thất bại: %w", item.ID, dispatchErr))
			continue
		}
		processedIDs = append(processedIDs, item.ID)
	}

	if len(processedIDs) > 0 {
		if deleteErr := uc.repo.DeleteByIDs(ctx, processedIDs); deleteErr != nil {
			errCollection.Append(fmt.Errorf("không thể xóa pending đã xử lý: %w", deleteErr))
		}
	}

	if errCollection.IsEmpty() {
		return nil
	}
	return errCollection
}

// ErrorCollection gom các lỗi phát sinh khi flush.
type ErrorCollection struct {
	Errors []error
}

func (ec *ErrorCollection) Append(err error) {
	if err == nil {
		return
	}
	ec.Errors = append(ec.Errors, err)
}

func (ec *ErrorCollection) IsEmpty() bool {
	return ec == nil || len(ec.Errors) == 0
}

func (ec *ErrorCollection) Error() string {
	if ec.IsEmpty() {
		return ""
	}
	return ec.Errors[0].Error()
}
