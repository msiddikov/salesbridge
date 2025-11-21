package zenotiv1

type (
	ItemType      string
	ZenotiStatus  int
	InvoiceStatus int
	PaymentType   int
	SaleType      int
)

const (
	NoShowed  ZenotiStatus = -2
	Canceled  ZenotiStatus = -1
	Booked    ZenotiStatus = 0
	Closed    ZenotiStatus = 1
	CheckedIn ZenotiStatus = 2
	Confirmed ZenotiStatus = 4
)

const (
	InvoiceOpen            InvoiceStatus = 0
	InvoiceProcessed       InvoiceStatus = 1
	InvoiceCampaignApplied InvoiceStatus = 2
	InvoiceCouponApplied   InvoiceStatus = 3
	InvoiceClosed          InvoiceStatus = 4
	InvoiceNotSpecified    InvoiceStatus = 11
	InvoiceVoided          InvoiceStatus = 99
)

const (
	PaymentTypeCash               PaymentType = 0
	PaymentTypeCard               PaymentType = 1
	PaymentTypeCheck              PaymentType = 2
	PaymentTypeCustomFinancial    PaymentType = 3
	PaymentTypeCustomNonFinancial PaymentType = 4
	PaymentTypeMembershipCredits  PaymentType = 5
	PaymentTypeMembershipBenefits PaymentType = 6
	PaymentTypePackages           PaymentType = 7
	PaymentTypeGiftCards          PaymentType = 8
	PaymentTypePrepaidCards       PaymentType = 9
	PaymentTypeLoyaltyPoint       PaymentType = 10
	PaymentTypeCashBack           PaymentType = 16
	PaymentTypeNoPayment          PaymentType = 32
)

var PaymentTypeMap = map[PaymentType]string{
	PaymentTypeCash:               "Cash",
	PaymentTypeCard:               "Card",
	PaymentTypeCheck:              "Check",
	PaymentTypeCustomFinancial:    "CustomFinancial",
	PaymentTypeCustomNonFinancial: "CustomNonFinancial",
	PaymentTypeMembershipCredits:  "MembershipCredits",
	PaymentTypeMembershipBenefits: "MembershipBenefits",
	PaymentTypePackages:           "Packages",
	PaymentTypeGiftCards:          "GiftCards",
	PaymentTypePrepaidCards:       "PrepaidCards",
	PaymentTypeLoyaltyPoint:       "LoyaltyPoint",
	PaymentTypeCashBack:           "CashBack",
	PaymentTypeNoPayment:          "NoPayment",
}

func GetAllPaymentTypes() []string {
	paymentTypes := []string{}
	for _, pt := range PaymentTypeMap {
		paymentTypes = append(paymentTypes, pt)
	}
	return paymentTypes
}

func GetPaymentType(str string) (PaymentType, bool) {
	for k, v := range PaymentTypeMap {
		if v == str {
			return k, true
		}
	}
	return -1, false
}

const (
	SaleTypeSale      SaleType = 0
	SaleTypeRefund    SaleType = 1
	SaleTypeRecurring SaleType = 2
	SaleTypeCharges   SaleType = 3
)

var SaleTypeMap = map[SaleType]string{
	SaleTypeSale:      "Sale",
	SaleTypeRefund:    "Refund",
	SaleTypeRecurring: "Recurring",
	SaleTypeCharges:   "Charges",
}

func GetAllSaleTypes() []string {
	saleTypes := []string{}
	for _, st := range SaleTypeMap {
		saleTypes = append(saleTypes, st)
	}
	return saleTypes
}

func GetSaleType(str string) (SaleType, bool) {
	for k, v := range SaleTypeMap {
		if v == str {
			return k, true
		}
	}
	return -1, false
}

const (
	ItemTypeService         ItemType = "Service"
	ItemTypeProduct         ItemType = "Product"
	ItemTypeMembership      ItemType = "Membership"
	ItemTypeDayPromoPackage ItemType = "DayPromoPackage"
	ItemTypeGiftCard        ItemType = "GiftCard"
	ItemTypePrepaidCard     ItemType = "PrepaidCard"
	ItemTypePackage         ItemType = "Package"
	ItemTypeOther           ItemType = "OtherFees"
)

func GetAllItemTypes() []string {
	return []string{
		string(ItemTypeService),
		string(ItemTypeProduct),
		string(ItemTypeMembership),
		string(ItemTypeDayPromoPackage),
		string(ItemTypeGiftCard),
		string(ItemTypePrepaidCard),
		string(ItemTypePackage),
		string(ItemTypeOther),
	}
}
