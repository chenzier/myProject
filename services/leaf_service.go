package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"product/datamodels"
	"product/repositories"
)

type LeafService struct {
	sql *sql.DB
	db  *repositories.LeafManager
}

func NewLeafDAO(db *repositories.LeafManager, sql *sql.DB) *LeafService {
	return &LeafService{
		db:  db,
		sql: sql,
	}
}

func (l *LeafService) NextSegment(ctx context.Context, bizTag string) (*datamodels.Leaf, error) {
	// 开启事务
	tx, err := l.sql.Begin()
	defer func() {
		if err != nil {
			l.rollback(tx)
		}
	}()
	if err = l.checkError(err); err != nil {
		return nil, err
	}
	err = l.db.UpdateMaxID(ctx, bizTag, tx)
	if err = l.checkError(err); err != nil {
		return nil, err
	}
	leaf, err := l.db.Get(ctx, bizTag, tx)
	if err = l.checkError(err); err != nil {
		return nil, err
	}
	// 提交事务
	err = tx.Commit()
	if err = l.checkError(err); err != nil {
		return nil, err
	}
	return leaf, nil
}

func (l *LeafService) checkError(err error) error {
	if err == nil {
		return nil
	}
	if message, ok := err.(*mysql.MySQLError); ok {
		fmt.Printf("it's sql error; str:%v", message.Message)
	}
	return errors.New("db error")
}

func (l *LeafService) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		fmt.Println("rollback error")
	}
}

func (l *LeafService) Add(ctx context.Context, leaf *datamodels.Leaf) error {
	return l.db.Create(ctx, leaf)
}

func (l *LeafService) Get(ctx context.Context, bizTag string) (*datamodels.Leaf, error) {
	return l.db.Get(ctx, bizTag, nil)
}
func (l *LeafService) UpdateMaxID(ctx context.Context, bizTag string) error {
	return l.db.UpdateMaxID(ctx, bizTag, nil)
}

func (l *LeafService) UpdateMaxIDByCustomStep(ctx context.Context, bizTag string, step int32) error {
	return l.db.UpdateMaxIdByCustomStep(ctx, step, bizTag, nil)
}

func (l *LeafService) UpdateStep(ctx context.Context, step int32, bizTag string) error {
	return l.db.UpdateStep(ctx, step, bizTag)
}
