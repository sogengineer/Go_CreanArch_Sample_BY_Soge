package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// 暗号(Hash)化
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Failed to generate password hash")
		return "", err
	}
	log.Info("Password encrypted successfully")
	return string(hash), nil
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.WithError(err).Error("Failed to compare password hash")
		return err
	}
	log.Info("Password hash compared successfully")
	return nil
}

// 第1引数にJson,第2引数に空の構造体を格納することで、空の第2引数に値を入力する。
func CopyBeans(target []byte, variable interface{}) error {
	if err := json.Unmarshal(target, &variable); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		return err
	}
	log.Info("JSON unmarshalled successfully")
	return nil
}

// 構造体をjson化させる関数
func ConvertStructIntoJson(source interface{}) ([]byte, error) {
	result, err := json.Marshal(&source)
	if err != nil {
		log.WithError(err).Error("Failed to marshal struct into JSON")
		return nil, err
	}
	log.Info("Struct marshalled into JSON successfully")
	return result, nil
}

func ConvertJsonAndCopyBean(source interface{}, target interface{}) error {
	json, err := ConvertStructIntoJson(source)
	if err != nil {
		log.WithError(err).Error("Failed to convert struct into JSON")
		return err
	}
	if err := CopyBeans(json, &target); err != nil {
		log.WithError(err).Error("Failed to copy JSON into target")
		return err
	}
	log.Info("JSON converted and copied into target successfully")
	return nil
}

// 新しいUUIDを生成。
func GenerateUUIDWithDate() string {
	newUUID := uuid.New()
	uuidStr := newUUID.String() // UUIDを文字列に変換
	shortHash := uuidStr[:20]
	log.WithField("uuid", shortHash).Info("UUID generated successfully")
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
	log.WithField("userId", shortHash).Info("User ID generated successfully")
	return shortHash
}
