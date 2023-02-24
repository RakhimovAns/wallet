package wallet

import (
	"errors"
	"github.com/RakhimovAns/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}
type Error string

func (e Error) Error() string {
	return string(e)
}

var ErrPhoneRegistered = errors.New("phone already registerd")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPhoneNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("account not found")
var ErrPaymentNotFound = errors.New("account not found")

func (s *Service) RegisterAccounts(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrPhoneNotFound
	}
	account.Balance += amount
	return nil
}
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrPhoneNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentsStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrPhoneNotFound
	}
	return account, nil
}

func (s *Service) Reject(paymentID string) error {
	for _, acc := range s.payments {
		if acc.ID == paymentID {
			acc.Status = types.PaymentsStatusFail
			account, err := s.FindAccountByID(acc.AccountID)
			account.Balance += acc.Amount
			return err
		}
	}
	return ErrPaymentNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, acc := range s.payments {
		if acc.ID == paymentID {
			payment = acc
			break
		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}
