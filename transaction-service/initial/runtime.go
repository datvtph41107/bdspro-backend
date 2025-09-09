package initial

import (
	"tx/infra/handler"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	TransactionHandler         *handler.TransactionHandler
	InternalTransactionHandler *handler.InternalTransactionHandler
	ContractDealHandler        *handler.ContractDealHandler
	DealCostHandler            *handler.DealCostHandler
	CostTypeHandler            *handler.CostTypeHandler
	Logger                     *flogging.FabricLogger
}

func NewInitalApp(
	transactionHandler *handler.TransactionHandler,
	internalTransactionHandler *handler.InternalTransactionHandler,
	contractDealHandler *handler.ContractDealHandler,
	dealCostHandler *handler.DealCostHandler,
	costTypeHandler *handler.CostTypeHandler,
) *InitialApp {
	return &InitialApp{
		TransactionHandler:         transactionHandler,
		InternalTransactionHandler: internalTransactionHandler,
		ContractDealHandler:        contractDealHandler,
		DealCostHandler:            dealCostHandler,
		CostTypeHandler:            costTypeHandler,
		Logger:                     flogging.MustGetLogger(viper.GetString("server.name")),
	}
}
