package testing

import (
	"database/sql/driver"
	"reflect"
	ori_testing "testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type login struct {
	ID               int32     `gorm:"primaryKey;autoIncrement"`
	UserID           string    `gorm:"index:logins_user_id_idx"`
	Email            string    `gorm:"not null;index:logins_email_idx"`
	PasswordSalt     string    `gorm:"not null"`
	PasswordHash     string    `gorm:"not null"`
	RequireResetPass bool      `gorm:"default:false"`
	ModifiedTime     time.Time `gorm:"default:now()"`
}

type user struct {
	ID           string `gorm:"primaryKey"`
	FullName     string
	Email        string `gorm:"not null;index:users_email_idx"`
	Phone        string
	Avatar       string
	Certificate  string
	CreatedTime  time.Time `gorm:"default:now()"`
	ModifiedTime time.Time `gorm:"default:now()"`
}

func TestStructFields(t *ori_testing.T) {
	u := login{}
	fields := getStructFields(u)
	if ok := reflect.DeepEqual(fields, []string{
		"ID",
		"UserID",
		"Email",
		"PasswordSalt",
		"PasswordHash",
		"RequireResetPass",
		"ModifiedTime",
	}); !ok {
		t.Fatal("expect to equals")
	}
}

func TestStructValue(t *ori_testing.T) {
	var id int32 = 34
	userID := "qujiabao"
	email := "qujiabao@example.com"
	salt := "FCRz"
	hash := "d03f5a65bbc703908dbcb84be6161b0023aaeb1d4bbce69a4f61a733a0b83de6"
	requireResetPass := false
	modifiedTime := time.Now()
	testUser := &login{
		ID:               id,
		UserID:           userID,
		Email:            email,
		PasswordSalt:     salt,
		PasswordHash:     hash,
		RequireResetPass: requireResetPass,
		ModifiedTime:     modifiedTime,
	}
	res := GetStructValues(testUser)
	expectRes := []driver.Value{
		id,
		userID,
		email,
		salt,
		hash,
		requireResetPass,
		modifiedTime,
	}
	assert.Equal(t, res, expectRes)
}

func TestJoinFields(t *ori_testing.T) {
	l := login{}
	u := user{}
	cols := getJoinedCols(&l, &u)
	if ok := reflect.DeepEqual(cols, []string{
		"ID",
		"UserID",
		"Email",
		"PasswordSalt",
		"PasswordHash",
		"RequireResetPass",
		"ModifiedTime",
		"user__ID",
		"user__FullName",
		"user__Email",
		"user__Phone",
		"user__Avatar",
		"user__Certificate",
		"user__CreatedTime",
		"user__ModifiedTime",
	}); !ok {
		t.Fatal("expect to equals")
	}
}
