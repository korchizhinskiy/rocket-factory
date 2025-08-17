package repository

type PaymentRepository interface {
	Pay(ctx context.Cotext, payment model.Payment)
}

