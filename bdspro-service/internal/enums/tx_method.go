package enums

type TxMethod uint32

const (
	// TxMethodWithdraw TxMethod = 20
	TxMethodDeposit    TxMethod = 10
	TxMethodSign       TxMethod = 20
	TxMethodApprove    TxMethod = 30
	TxMethodCost       TxMethod = 40
	TxMethodInvestment TxMethod = 50
	TxMethodRevenue    TxMethod = 60
	TxMethodRent       TxMethod = 70
	TxMethodSale       TxMethod = 80
	TxMethodBuy        TxMethod = 90
	TxMethodTransfer   TxMethod = 100
)

var TxMethodNames = map[TxMethod]string{
	TxMethodSign:       "Ký hợp đồng", // sign (Ký, chuyển nhượng, bàn giao, nhận bàn giao), approve/reject
	TxMethodCost:       "Chi phí",     // transfer, approve/reject (Chi phí)
	TxMethodRevenue:    "Thu nhập",    // approve (Thu nhập)
	TxMethodInvestment: "Góp vốn",     // transfer, approve/reject
	TxMethodDeposit:    "Đặt cọc",     // transfer, approve/reject
	TxMethodRent:       "Cho thuê",    // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
	TxMethodSale:       "Bán",         // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
	TxMethodBuy:        "Mua",         // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
	TxMethodTransfer:   "Chuyển tiền", // transfer (Chuyển tiền)
}

// IsValidTxMethod checks if the given method is defined in TxMethodNames
func IsValidTxMethod(method TxMethod) bool {
	_, ok := TxMethodNames[method]
	return ok
}
