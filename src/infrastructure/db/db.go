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

func GetDB(ctx context.Context) (*gorm.DB, error) {
	if tx, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		// コンテキストにトランザクションがある場合、それを返す
		return tx, nil
	}
	deadline := time.After(5 * time.Second) // タイムアウトの設定

	for {
		connValue := connection.Load()
		if connValue != nil {
			conn := connValue.(*DBConnection)

			select {
			case <-conn.ctx.Done():
				// 接続のコンテキストがキャンセルされている場合、少し待ってから再試行
				time.Sleep(time.Second)
			case <-ctx.Done():
				// 呼び出し元のコンテキストがキャンセルされた場合、エラーを返す
				return nil, ctx.Err()
			case <-deadline:
				// タイムアウトした場合、エラーを返す
				return nil, fmt.Errorf("GetDB: データベース接続待機中にタイムアウト")
			default:
				// どちらのコンテキストもキャンセルされていない場合、データベース接続を返す
				return conn.db, nil
			}
		} else {
			// 接続がない場合、少し待ってから再試行
			time.Sleep(time.Second)
		}
	}
}

func Find(ctx context.Context, query string, args []interface{}, out interface{}) error {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	if err := db.WithContext(ctx).Where(query, args...).First(out).Error; err != nil {
		return err
	}

	return nil
}

func FindAll(ctx context.Context, query string, args []interface{}, out interface{}) error {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	if err := db.WithContext(ctx).Where(query, args...).Find(out).Error; err != nil {
		return err
	}

	return nil
}

func FindWithRawJoinQuery(ctx context.Context, sqlQuery string, out interface{}, params ...interface{}) error {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	if err := db.Raw(sqlQuery, params...).Find(out).Error; err != nil {
		return err
	}

	return nil
}

func Create(ctx context.Context, value interface{}) error {
	// データベース接続を取得
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	// データベースにレコードを作成
	if err := db.WithContext(ctx).Create(value).Error; err != nil {
		return fmt.Errorf("データベースにレコードを作成できませんでした: %w", err)
	}

	return nil
}

func Delete(ctx context.Context, query string, args []interface{}, out interface{}) error {
	// データベース接続を取得
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	// データベースにレコードを作成
	if err := db.WithContext(ctx).Where(query, args...).Delete(out).Error; err != nil {
		return fmt.Errorf("データベースのレコードを削除できませんでした: %w", err)
	}

	return nil
}

func WithTransaction(ctx context.Context, fn TransactionFunc) error {
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	tx := db.WithContext(ctx).Begin()
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

func Update(ctx context.Context, value interface{}) error {
	// データベース接続を取得
	db, err := GetDB(ctx)
	if err != nil {
		return fmt.Errorf("データベース接続の取得に失敗: %w", err)
	}

	// レコード更新
	if err := db.WithContext(ctx).Save(value).Error; err != nil {
		return fmt.Errorf("データベースのレコードを更新できませんでした: %w", err)
	}

	return nil
}
