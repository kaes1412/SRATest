package domain

type Loan struct {
	ID            string
	Principal     int64
	InterestRate  float64
	TotalWeeks    int
	WeeklyPayment int64
	Payments      []Payment
}

type Payment struct {
	Week int
	Paid bool
}
