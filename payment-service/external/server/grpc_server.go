package grpc_server

import (
	_jwt "common/jwt"
	_utils "common/utils"
	"context"
	"fmt"
	"payment/internal/wallet/domain"
	"payment/internal/wallet/usecases"
	pb_payment "pb/types/payment"
	pb_spec "pb/types/spec"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb_payment.UnimplementedPaymentServiceServer
	accountUc     *usecases.AcountUsecase
	transactionUc *usecases.TransactionUsecase
}

func NewGrpcServer(accountUc *usecases.AcountUsecase,
	transactionUc *usecases.TransactionUsecase,
) *GrpcServer {
	return &GrpcServer{
		accountUc:     accountUc,
		transactionUc: transactionUc,
	}
}

func (s *GrpcServer) Register(c context.Context, data *pb_payment.RegisterRequest) (*pb_payment.RegisterResponse, error) {
	dto := &domain.Account{}
	if err := _utils.MapWithValidation(data, dto); err != nil {
		// st := status.New(codes.Code(err.Code), err.Message)
		// stWithDetails, _ := st.WithDetails(&pb_spec.FieldErrors{
		// 	Errors: err.Errors,
		// })
		// return nil, stWithDetails.Err()
		return nil, _utils.ErrorField(err, &pb_spec.FieldErrors{
			Errors: err.Errors,
		})
	}

	_, err := s.accountUc.Register(dto)
	return &pb_payment.RegisterResponse{Success: true}, err
}
func (s *GrpcServer) MakeTransaction(c context.Context, data *pb_payment.TransactionRequest) (*pb_payment.TransactionResponse, error) {
	dto := &domain.Transaction{}
	if err := _utils.MapWithValidation(data, dto); err != nil {
		st := status.New(codes.Code(err.Code), err.Message)
		stWithDetails, _ := st.WithDetails(&pb_spec.FieldErrors{
			Errors: err.Errors,
		})
		return nil, stWithDetails.Err()
	}

	metadata := _jwt.MetadataFromContext(c)
	if metadata.ProfileID != nil {
		dto.AccountId = *metadata.ProfileID
	} else {
		return nil, fmt.Errorf("không xác minh người dùng")
	}

	_, err := s.transactionUc.MakeTransaction(dto)
	return &pb_payment.TransactionResponse{
		Data: 1,
	}, err
	// return nil, status.Errorf(codes.Unimplemented, "method MakeTransaction not implemented")
}
func (s *GrpcServer) Balance(c context.Context, data *pb_payment.Empty) (*pb_payment.BalanceRes, error) {
	metadata := _jwt.MetadataFromContext(c)

	balance, err := s.transactionUc.Balance(metadata.ProfileID)

	return &pb_payment.BalanceRes{
		Balance: balance,
	}, err
}
