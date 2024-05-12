package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	db  *gorm.DB
	ctx context.Context
}

func NewDBConnection(ctx context.Context) (*DBConnection, error) {
	// データベース接続を取得
	dbInstance, err := getDB(ctx)
	if err != nil {
		return nil, err
	} else {
		dbConnect := &DBConnection{
			db:  dbInstance,
			ctx: ctx,
		}
		return dbConnect, nil
	}
}

type contextKey string
type TransactionFunc func(tx *gorm.DB) error

const CancelKey = "cancel"
const dbKey contextKey = "transactionDB"

var (
	connection atomic.Value // *DBConnection
)

func Init() {
	// プログラム起動時に最初のDB接続を確立します
	err := establishConnection()
	if err != nil {
		log.Fatalf("データベース接続の初期化エラー: %v", err)
	}
}

// establishConnection はデータベースの接続を設定します
func establishConnection() error {
	dsn := fmt.Sprintf("host=db port=5432 user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))

	// 新しいデータベース接続を開きます
	newDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	// データベース接続を保存します
	conn := &DBConnection{
		db:  newDb,
		ctx: context.Background(),
	}
	connection.Store(conn)

	return nil
}

func getDB(ctx context.Context) (*gorm.DB, error) {
	if tx, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		// コンテキストにトランザクションがある場合、それを返す
		return tx, nil
	}

	const (
		maxRetries = 3
		retryDelay = time.Second
	)

	deadline, ok := ctx.Deadline()
	if !ok {
		// デッドラインが設定されていない場合、デフォルトのタイムアウトを設定
		deadline = time.Now().Add(5 * time.Second)
	}

	for i := 0; i < maxRetries; i++ {
		connValue := connection.Load()
		if connValue != nil {
			conn := connValue.(*DBConnection)

			select {
			case <-conn.ctx.Done():
				// 接続のコンテキストがキャンセルされている場合、少し待ってから再試行
				time.Sleep(retryDelay)
			case <-ctx.Done():
				// 呼び出し元のコンテキストがキャンセルされた場合、エラーを返す
				return nil, ctx.Err()
			default:
				// どちらのコンテキストもキャンセルされていない場合、データベース接続を返す
				return conn.db, nil
			}
		} else {
			// 接続がない場合、少し待ってから再試行
			select {
			case <-ctx.Done():
				// 呼び出し元のコンテキストがキャンセルされた場合、エラーを返す
				return nil, ctx.Err()
			case <-time.After(retryDelay):
				// 少し待ってから次の試行へ
			}
		}

		// タイムアウトチェック
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("getDB: データベース接続の取得がタイムアウトしました")
		}
	}

	return nil, fmt.Errorf("getDB: データベース接続の取得に失敗しました（リトライ回数: %d）", maxRetries)
}

func (dbConnect DBConnection) Find(ctx context.Context, query string, args []interface{}, out interface{}) error {

	if err := dbConnect.db.WithContext(ctx).Where(query, args...).First(out).Error; err != nil {
		return err
	}

	return nil
}

func (dbConnect DBConnection) FindAll(ctx context.Context, query string, args []interface{}, out interface{}) error {
	if err := dbConnect.db.WithContext(ctx).Where(query, args...).Find(out).Error; err != nil {
		return err
	}

	return nil
}

func (dbConnect DBConnection) FindWithRawJoinQuery(ctx context.Context, sqlQuery string, out interface{}, params ...interface{}) error {
	if err := dbConnect.db.Raw(sqlQuery, params...).Find(out).Error; err != nil {
		return err
	}

	return nil
}

func (dbConnect DBConnection) Create(ctx context.Context, value interface{}) error {
	// データベースにレコードを作成
	if err := dbConnect.db.WithContext(ctx).Create(value).Error; err != nil {
		return fmt.Errorf("データベースにレコードを作成できませんでした: %w", err)
	}

	return nil
}

func (dbConnect DBConnection) Delete(ctx context.Context, query string, args []interface{}, out interface{}) error {
	// データベース接続を取得
	if err := dbConnect.db.WithContext(ctx).Where(query, args...).Delete(out).Error; err != nil {
		return fmt.Errorf("データベースのレコードを削除できませんでした: %w", err)
	}

	return nil
}

func (dbConnect DBConnection) WithTransaction(ctx context.Context, fn TransactionFunc) error {

	tx := dbConnect.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("トランザクションの開始に失敗: %w", tx.Error)
	}

	if err := fn(tx); err != nil { // この行を元に戻します
		tx.Rollback()
		return fmt.Errorf("トランザクション中のエラー: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗: %w", err)
	}

	return nil
}

func (dbConnect DBConnection) Update(ctx context.Context, value interface{}) error {
	// レコード更新
	if err := dbConnect.db.WithContext(ctx).Save(value).Error; err != nil {
		return fmt.Errorf("データベースのレコードを更新できませんでした: %w", err)
	}

	return nil
}
