package enums

type EMethod uint32

const (
	// EMethodWithdraw EMethod = 20
	EMethodDeposit    EMethod = 10
	EMethodSign       EMethod = 20
	EMethodApprove    EMethod = 30
	EMethodCost       EMethod = 40
	EMethodInvestment EMethod = 50
	EMethodRevenue    EMethod = 60
	EMethodRent       EMethod = 70
	EMethodSale       EMethod = 80
	EMethodBuy        EMethod = 90
)

var MethodNames = map[EMethod]string{
	EMethodSign:       "Ký hợp đồng", // sign (Ký, chuyển nhượng, bàn giao, nhận bàn giao), approve/reject
	EMethodCost:       "Chi phí",     // transfer, approve/reject (Chi phí)
	EMethodRevenue:    "Thu nhập",    // approve (Thu nhập)
	EMethodInvestment: "Góp vốn",     // transfer, approve/reject
	EMethodDeposit:    "Đặt cọc",     // transfer, approve/reject
	EMethodRent:       "Cho thuê",    // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
	EMethodSale:       "Bán",         // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
	EMethodBuy:        "Mua",         // transfer(cọc), transfer(chuyển tiền), sign, approve/reject
}

// IsValidMethod checks if the given method is defined in MethodNames
func IsValidMethod(method EMethod) bool {
	_, ok := MethodNames[method]
	return ok
}
