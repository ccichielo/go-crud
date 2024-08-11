package mocks

import (
	"github.com/ccichielo/gobank/pkg"
)

type MockStorage struct {
	GetAccountsFunc    func() ([]*pkg.Account, error)
	GetAccountByIDFunc func(int) (*pkg.Account, error)
	CreateAccountFunc  func(*pkg.Account) error
	DeleteAccountFunc  func(int) error
	TransferFunc       func(int, int, int) error
}

func (m *MockStorage) GetAccounts() ([]*pkg.Account, error) {
	if m.GetAccountsFunc != nil {
		return m.GetAccountsFunc()
	}

	return nil, nil
}

func (m *MockStorage) GetAccountByID(id int) (*pkg.Account, error) {
	if m.GetAccountByIDFunc != nil {
		return m.GetAccountByIDFunc(id)
	}
	return nil, nil
}

func (m *MockStorage) CreateAccount(account *pkg.Account) error {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(account)
	}
	return nil
}

func (m *MockStorage) DeleteAccount(id int) error {
	if m.DeleteAccountFunc != nil {
		return m.DeleteAccountFunc(id)
	}
	return nil
}

func (m *MockStorage) Transfer(from int, to int, amount int) error {
	if m.TransferFunc != nil {
		return m.Transfer(from, to, amount)
	}
	return nil
}
