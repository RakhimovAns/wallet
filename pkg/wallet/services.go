package wallet

import (
	"errors"
	"fmt"
	"github.com/RakhimovAns/wallet/pkg/types"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strconv"
)

type Service struct {
	nextAccountID int64
	Accounts      []*types.Account
	payments      []*types.Payment
	favorite      []*types.Favorite
}
type testAccount struct {
	phone    types.Phone
	Balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}
type Error string

func (e Error) Error() string {
	return string(e)
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

var ErrPhoneRegistered = errors.New("phone already registerd")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPhoneNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("account not found")
var ErrPaymentNotFound = errors.New("account not found")

func (s *Service) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccounts(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register account error = %v", err)
	}
	err = s.Deposit(account.ID, data.Balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account , error =%v", err)
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make pamynet,error  = %v", err)
		}
	}
	return account, payments, nil
}

func (s *Service) RegisterAccounts(phone types.Phone) (*types.Account, error) {
	for _, account := range s.Accounts {
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
	s.Accounts = append(s.Accounts, account)
	return account, nil
}
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.Accounts {
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

	for _, acc := range s.Accounts {
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
	for _, acc := range s.Accounts {
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

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccounts(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account ,error =%v", err)
	}
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account ,error=%v", err)
	}
	return account, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
			break
		}
	}
	if payment == nil {
		return nil, fmt.Errorf("Repeat(): payment can't be found ,error :=%v", ErrPaymentNotFound)
	}
	var pay types.Payment
	pay = *payment
	pay.ID = uuid.New().String()
	payment = &pay
	return payment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	favorite := types.Favorite{
		ID:        payment.ID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}
	var Favorite *types.Favorite
	Favorite = &favorite
	return Favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var Favorite *types.Favorite
	for _, fav := range s.favorite {
		if fav.ID == favoriteID {
			Favorite = fav
			break
		}
	}
	if Favorite == nil {
		return nil, fmt.Errorf("Favorite wasn't found")
	}
	return s.Pay(Favorite.AccountID, Favorite.Amount, Favorite.Category)
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(err)
		}
	}()
	for _, acc := range s.Accounts {
		_, err := file.Write([]byte(strconv.FormatInt(int64(acc.ID), 10) + ";" + string(acc.Phone) + "|"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(err)
		}
	}()
	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := string(content)
	temp := ""
	for _, i := range data {
		temp += string(i)
		if string(i) == "|" {
			temporary := ""
			for _, j := range temp {
				temporary += string(j)
				if string(j) == ";" {
					temporary = ""
				}
			}
			runes := []rune(temporary)
			runes = runes[:len(runes)-1]
			temporary = string(runes)
			s.RegisterAccounts(types.Phone(temporary))
			temp = ""
		}
	}
	return nil
}
