package order_payment_model

type Order struct {
	Order_ID
	Order_Cart
	Ordered_At
	Price
	Discount
	Payment_Method
}

type Payment struct {
	Digital bool
	Cash    bool
}
