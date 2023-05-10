package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/prayogatriady/sawer-transaction/model"
	"github.com/prayogatriady/sawer-transaction/repository"
)

type TransactionServInterface interface {
	Sawer(ctx context.Context, trRequest *model.TransactionRequest) (*model.TransactionResponse, error)
	PaymentNotification(ctx context.Context, notifPayload *model.NotificationPayload) error
}

type TransactionService struct {
	TransactionRepository repository.TransactionRepoInterface
	s                     snap.Client
}

func NewTransactionService(trRepository repository.TransactionRepoInterface, s snap.Client) TransactionServInterface {
	return &TransactionService{
		TransactionRepository: trRepository,
		s:                     s,
	}
}

func (ts *TransactionService) PaymentNotification(ctx context.Context, notifPayload *model.NotificationPayload) error {

	trxId, _ := strconv.Atoi(strings.Split(notifPayload.OrderId, "_")[0])

	trEntity, err := ts.TransactionRepository.GetTrx(ctx, trxId)
	if err != nil {
		return err
	}

	trEntity.TransactionStatus = notifPayload.TransactionStatus

	if err := ts.TransactionRepository.UpdatePayment(ctx, trEntity); err != nil {
		return err
	}

	return nil
}

func (ts *TransactionService) Sawer(ctx context.Context, trRequest *model.TransactionRequest) (*model.TransactionResponse, error) {

	trEntity := &model.TransactionEntity{
		UserId:      trRequest.UserId,
		Amount:      trRequest.Amount,
		SawerUserId: trRequest.SawerUserId,
	}

	trEntity, err := ts.TransactionRepository.Sawer(ctx, trEntity)
	if err != nil {
		return &model.TransactionResponse{}, err
	}

	// Initiate client for Midtrans Snap
	ts.s.New("SB-Mid-server-CbaeF5Q5pwRT3p-uAsIiQHTG", midtrans.Sandbox)
	// Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(trEntity.ID) + "_" + trEntity.CreatedAt.String(),
			GrossAmt: int64(trEntity.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}
	// Request create Snap transaction to Midtrans
	snapResp, _ := ts.s.CreateTransaction(req)

	trResponse := &model.TransactionResponse{
		ID:          trEntity.ID,
		UserId:      trEntity.UserId,
		Amount:      trEntity.Amount,
		SawerUserId: trEntity.SawerUserId,
		CreatedAt:   trEntity.CreatedAt,
		UpdatedAt:   trEntity.UpdatedAt,
		Midtrans:    snapResp,
	}

	return trResponse, nil
}
