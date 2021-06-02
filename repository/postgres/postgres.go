package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/junereycasuga/gokit-grpc-demo/repository"
	_ "github.com/lib/pq"
)

type ConnParam struct {
	Host        string
	Port        string
	DBName      string
	User        string
	Pass        string
	Options     string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

type postgres struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewPostgresSql(p *ConnParam) (repository.Repository, error) {
	// psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	p.Host, p.Port, p.User, p.Pass, p.DBName)

	psqlconn := "postgres://postgres:123@localhost:5432/demo_cpq2?sslmode=disable&search_path=public"
	print(psqlconn)
	db, err := sql.Open("postgres", psqlconn)
	print("\n db: ", db)
	print("\n dberr: ", err)
	if err != nil {
		print("\n here: ", err)
		return nil, fmt.Errorf("db open: %v", err)
	}
	if err := db.Ping(); err != nil {
		print("\n PONG22")
		return nil, err
	}
	print("\n PONG")
	db.SetMaxOpenConns(p.MaxOpenConn)
	db.SetMaxIdleConns(p.MaxIdleConn)
	db.SetConnMaxLifetime(p.MaxLifetime)
	return &postgres{db: db}, nil
}

// Close ...
func (p *postgres) Close() error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return err
		}
		p.db = nil
	}
	return nil
}

const (
	getTransaction = "select pt.id,pt.identifier ,pt.fare ,pt.cv_number ,pt.created_at ,pt.sap_sent_at,pt.order_id ,pt.itop_id ,pt.paid_amount ,pt.discount_amt ,a.account_code from payment_transactions pt"
)

func (p *postgres) GetTransaction(ctx context.Context, accountNo, postingDateStart, postingDateEnd string) (result []repository.Transaction, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	sql := fmt.Sprintf("%s left join accounts a on pt.account_id = a.id where a.account_code =$1 and pt.created_at between $2 and $3", getTransaction)

	row, err := p.db.Query(sql, accountNo, postingDateStart, postingDateEnd)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := row.Close(); e != nil {
			err = e
		}
	}()

	var data repository.Transaction
	for row.Next() {
		err := row.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("row scan: %v", err)
		}
		result = append(result, data)
	}
	return result, nil
}

func (p *postgres) GetTransaction2(ctx context.Context) (result []repository.Transaction, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	sql := "select id,offer_id,quote_id from lul where id IN (18,19)"

	row, err := p.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := row.Close(); e != nil {
			err = e
		}
	}()

	var data repository.Transaction
	for row.Next() {
		err := row.Scan(&data.Id, &data.Offer_id, &data.Quote_id)
		if err != nil {
			return nil, fmt.Errorf("row scan: %v", err)
		}
		result = append(result, data)
	}
	return result, nil
}
