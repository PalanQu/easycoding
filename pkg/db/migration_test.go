package db

import (
	"easycoding/common/workspace"
	"fmt"
	"path"
	"testing"
)

type testData struct {
	From            int
	To              int
	ExpectDirection MigrationDirection
	ExceptStep      int
	ExceptErr       error
}

var dirname = path.Join(workspace.GetWorkspace(), "pkg", "db", "test", "migrations")

func TestMigration(t *testing.T) {
	dataList := []testData{
		{
			From:            20220914020936,
			To:              20220914055720,
			ExpectDirection: MigrationDirectionUP,
			ExceptStep:      1,
			ExceptErr:       nil,
		},
		{
			From:            20220914020936,
			To:              20221120094056,
			ExpectDirection: MigrationDirectionUP,
			ExceptStep:      4,
			ExceptErr:       nil,
		},
		{
			From:            20221120094056,
			To:              20220914020936,
			ExpectDirection: MigrationDirectionDown,
			ExceptStep:      4,
			ExceptErr:       nil,
		},
	}
	for _, data := range dataList {
		t.Run(fmt.Sprintf("from: %v, to: %v", data.From, data.To), func(t *testing.T) {
			d, step, err := MigrationGenerate(dirname, data.From, data.To)
			if d != data.ExpectDirection {
				t.Errorf("from: %v, to: %v, error direction %s, expected %s", data.From, data.To, d, data.ExpectDirection)
			}
			if step != data.ExceptStep {
				t.Errorf("from: %v, to: %v, error step %v, expected %v", data.From, data.To, step, data.ExceptStep)
			}
			if err != data.ExceptErr {
				t.Errorf("from: %v, to: %v, error %v, expected %v", data.From, data.To, err, data.ExceptErr)
			}
		})
	}
}
