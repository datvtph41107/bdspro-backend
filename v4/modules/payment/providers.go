package payment

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"userv4/modules/payment/controller"
	paymentpostgres "userv4/modules/payment/postgres"
	"userv4/modules/payment/repository"
	"userv4/modules/payment/services"
)

// ProviderSet wires modules/payment on PaymentDB.
var ProviderSet = wire.NewSet(
	ProvideBankRepository,
	services.NewLocalBankService,
	wire.Bind(new(services.BankService), new(*services.LocalBankService)),
	controller.NewHandler,
)

func ProvideBankRepository(db *gorm.DB) repository.BankRepository {
	return paymentpostgres.NewBankPostgres(db)
}
