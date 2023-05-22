package midtransHandler

import "gorm.io/gorm"

type Repository interface {
	Save(paymentRequest PaymentRequest) (PaymentRequest, error)
	Update(paymentRequest PaymentRequest) (PaymentRequest, error)
	FindByUUID(id string) (PaymentRequest, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db,
	}
}

func (r *repository) Save(paymentRequest PaymentRequest) (PaymentRequest, error) {
	err := r.db.Create(&paymentRequest).Error

	if err != nil {
		return paymentRequest, err
	}

	return paymentRequest, nil
}

func (r *repository) Update(paymentRequest PaymentRequest) (PaymentRequest, error) {
	err := r.db.Save(&paymentRequest).Error

	if err != nil {
		return paymentRequest, err
	}

	return paymentRequest, nil
}

func (r *repository) FindByUUID(id string) (PaymentRequest, error) {

	var paymentRequest PaymentRequest
	err := r.db.Where("uuid = ?", id).Find(&paymentRequest).Error

	if err != nil {
		return paymentRequest, err
	}

	return paymentRequest, nil
}
