package repository

import (
	"context"
	"time"

	"github.com/prayogatriady/sawer-transaction/model"
	"gorm.io/gorm"
)

type TransactionRepoInterface interface {
	Sawer(ctx context.Context, tr *model.TransactionEntity) (*model.TransactionEntity, error)
	UpdatePayment(ctx context.Context, tr *model.TransactionEntity) error
	GetTrx(ctx context.Context, trxId int) (*model.TransactionEntity, error)
}

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepoInterface {
	return &TransactionRepository{
		DB: db,
	}
}

func (ur *TransactionRepository) GetTrx(ctx context.Context, trxId int) (*model.TransactionEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var tr *model.TransactionEntity
	if err := ur.DB.WithContext(ctx).Table("transactions").Where("id =?", trxId).Find(&tr).Error; err != nil {
		return tr, err
	}
	return tr, nil
}

func (ur *TransactionRepository) UpdatePayment(ctx context.Context, tr *model.TransactionEntity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := ur.DB.WithContext(ctx).Table("transactions").Where("id =?", tr.ID).Update("transaction_status", tr.TransactionStatus).Error; err != nil {
		return err
	}

	if err := ur.DB.WithContext(ctx).Table("users").Where("id =?", tr.UserId).Update("balance", gorm.Expr("balance - ?", tr.Amount)).Error; err != nil {
		return err
	}

	if err := ur.DB.WithContext(ctx).Table("users").Where("id =?", tr.SawerUserId).Update("balance", gorm.Expr("balance + ?", tr.Amount)).Error; err != nil {
		return err
	}

	return nil
}

func (ur *TransactionRepository) Sawer(ctx context.Context, tr *model.TransactionEntity) (*model.TransactionEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := ur.DB.WithContext(ctx).Table("transactions").Create(&tr).Error; err != nil {
		return tr, err
	}

	// if err := ur.DB.WithContext(ctx).Table("users").Where("id =?", tr.UserId).Update("balance", gorm.Expr("balance - ?", tr.Amount)).Error; err != nil {
	// 	return nil, err
	// }

	// if err := ur.DB.WithContext(ctx).Table("users").Where("id =?", tr.SawerUserId).Update("balance", gorm.Expr("balance + ?", tr.Amount)).Error; err != nil {
	// 	return nil, err
	// }

	return tr, nil
}
