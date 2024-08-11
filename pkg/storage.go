package pkg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	CreateAccount(*Account) error
	DeleteAccount(int) error
	Transfer(int, int, int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id serial primary key,
    first_name varchar(50),
    last_name varchar(50),
    number serial,
    balance serial,
    created_at timestamp
  )`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account (first_name, last_name, number, balance, created_at)
  VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) addFundsToAccountID(tx *sql.Tx, id int, amount int) error {
	query := `
    UPDATE account
    SET balance = balance + $1
    WHERE id = $2
  `

	_, err := tx.Exec(query, amount, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) Transfer(from int, to int, amount int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = s.addFundsToAccountID(tx, from, -amount)
	if err != nil {
		return err
	}
	err = s.addFundsToAccountID(tx, to, amount)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE ID = $1`

	_, err := s.db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM account WHERE ID = $1`

	row := s.db.QueryRow(query, id)

	account := &Account{}
	err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM account`
	res, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for res.Next() {
		account := &Account{}
		err := res.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
