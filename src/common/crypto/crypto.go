package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 暗号(Hash)化
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// 第1引数にJson,第2引数に空の構造体を格納することで、空の第2引数に値を入力する。
func CopyBeans(target []byte, variable interface{}) error {
	if err := json.Unmarshal(target, &variable); err != nil {
		return err
	}
	return nil
}

// 構造体をjson化させる関数
func ConvertStructIntoJson(source interface{}) ([]byte, error) {
	result, err := json.Marshal(&source)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ConvertJsonAndCopyBean(source interface{}, target interface{}) error {
	json, err := ConvertStructIntoJson(source)
	if err != nil {
		return err
	}
	if err := CopyBeans(json, &target); err != nil {
		return err
	}

	return nil
}

// 新しいUUIDを生成。
func GenerateUUIDWithDate() string {
	newUUID := uuid.New()
	uuidStr := newUUID.String() // UUIDを文字列に変換
	shortHash := uuidStr[:20]
	return shortHash
}

// ユーザーIDの暗号化
func GenerateUserId(email string) string {
	// SHA-256ハッシュ関数を使用してハッシュ値を生成
	hasher := sha256.New()
	hasher.Write([]byte(email))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// ハッシュ値から先頭6桁を取得
	shortHash := hash[:12]
	// ユーザーIDを生成
	return shortHash
}
