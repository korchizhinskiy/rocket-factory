package enum

type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

func (pm *PaymentMethod) String() string {
	switch *pm {
	case PaymentMethodCard:
		return "Card"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "Credit Card"
	case PaymentMethodInvestorMoney:
		return "Investor Money"
	default:
		return "Unspecified"
	}
}
