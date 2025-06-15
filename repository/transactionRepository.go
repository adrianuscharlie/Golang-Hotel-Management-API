package repository

import (
	"hms-backend/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(t *model.Transaction) error
	GetByID(i string) (model.Transaction, error)
	Update(t *model.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(t *model.Transaction) error {
	return r.db.Table("transactions").Create(t).Error
}

func (r *transactionRepository) GetByID(i string) (model.Transaction, error) {
	var t model.Transaction
	err := r.db.Preload("bookings").Table("transactions").Where("id = ?", i).First(&t).Error
	return t, err
}

func (r *transactionRepository) Update(t *model.Transaction) error {
	return r.db.Model(&model.Transaction{}).Where("id = ?", t.Id).Updates(t).Error
}
