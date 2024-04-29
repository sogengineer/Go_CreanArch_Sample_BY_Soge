package user

import (
	"context"
	"fmt"

	"github.com/Go_CleanArch/common/crypto"
	"github.com/Go_CleanArch/infrastructure/db"
	dbConnect "github.com/Go_CleanArch/infrastructure/db"
	"github.com/Go_CleanArch/interface_adapter/gateway/entity"
	repository "github.com/Go_CleanArch/usecase/repository_interface"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	ctx context.Context
}

// コンストラクタ
func NewUserRepository(ctx context.Context) (repository.UserRepository, error) {
	dbConnectionResult, err := db.GetDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB接続に失敗しました")
	}

	result := userRepository{
		db:  dbConnectionResult,
		ctx: ctx,
	}

	return &result, nil
}

// ユーザーレコード作成
func (userRepository *userRepository) CreateUser(ctx context.Context, userJson []byte) (*entity.User, error) {
	var user entity.User
	if err := crypto.CopyBeans(userJson, &user); err != nil {
		return nil, err
	}

	// 新しいCreateメソッドを使用してデータベースにユーザーを作成
	if err := dbConnect.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Userの存在チェック
func (userRepository *userRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	// dbConnectのFind関数を使用
	err := dbConnect.Find(ctx, "email = ?", []interface{}{email}, &user)
	if err == gorm.ErrRecordNotFound {
		// レコードが見つからなかったエラー
		return nil, fmt.Errorf("条件に一致するレコードが見つかりません: %w", err)
	} else if err != nil {
		// その他のエラー
		fmt.Print(err)
		return nil, fmt.Errorf("DB検索に失敗しました: %w", err)
	}

	return &user, nil
}
