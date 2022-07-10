package testing

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	db "easycoding/pkg/db"
)

func SetupMockDB() (sqlmock.Sqlmock, *gorm.DB, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gdb, err := db.CreateTestingGdb(mockDB)
	if err != nil {
		return nil, nil, err
	}
	return mock, gdb, nil
}

func getStructFields(ormStruct interface{}) []string {
	fields := []string{}
	v := reflect.Indirect(reflect.ValueOf(ormStruct)).Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Type.Kind() {
		case reflect.Ptr, reflect.Slice:
			continue
		default:
			fields = append(fields, field.Name)
		}
	}
	return fields
}

// GetTableCols generate some sqlmock rows from an orm struct.
// For example
// 		GetTableCols(&orm.Login{}) =>
//			*sqlmock.Rows {
//				"ID",
//				"UserID",
//				"Email",
//				"PasswordSalt",
//				"PasswordHash",
//				"RequireResetPass",
//				"ModifiedTime",
//			}
func GetTableCols(s interface{}) *sqlmock.Rows {
	fields := getStructFields(s)
	return sqlmock.NewRows(fields)
}

func GetRawTableColsWithPrefix(s interface{}, prefix string) []string {
	fields := getStructFields(s)
	if prefix == "" {
		return fields
	}
	for i := 0; i < len(fields); i++ {
		fieldName := fmt.Sprintf("%s__%s", prefix, fields[i])
		fields[i] = fieldName
	}
	return fields
}

func GetJoinedTableCols(
	mainStruct interface{}, joinedStruct ...interface{}) *sqlmock.Rows {
	fields := getJoinedCols(mainStruct, joinedStruct...)
	return sqlmock.NewRows(fields)
}

func getJoinedCols(mainStruct interface{}, joinedStructs ...interface{}) []string {
	mainFields := getStructFields(mainStruct)
	joinedFields := []string{}
	for _, s := range joinedStructs {
		v := reflect.Indirect(reflect.ValueOf(s)).Type()
		prefix := v.Name()
		fields := getStructFields(s)
		for i := 0; i < len(fields); i++ {
			fieldName := fmt.Sprintf("%s__%s", prefix, fields[i])
			joinedFields = append(joinedFields, fieldName)
		}
	}
	return append(mainFields, joinedFields...)
}

// GetStructValues get struct values from the given struct.
// Read the test case for more information.
func GetStructValues(ormStruct interface{}) []driver.Value {
	v := reflect.ValueOf(ormStruct).Elem()
	values := []driver.Value{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Ptr, reflect.Slice:
			continue
		default:
			values = append(values, field.Interface())
		}
	}
	return values
}
