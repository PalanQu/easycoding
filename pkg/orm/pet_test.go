package orm

import (
	"testing"

	pet_pb "easycoding/api/pet"
	"easycoding/pkg/errors"
	test_utils "easycoding/pkg/testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetPet(t *testing.T) {
	mock, db, _ := test_utils.SetupMockDB()
	pet1 := &Pet{
		ID:   1,
		Name: "cat1",
		Type: int32(pet_pb.PetType_PET_TYPE_CAT),
	}
	t.Run("test get pet", func(t *testing.T) {
		mock.ExpectQuery("^SELECT (.+) FROM \\`pet\\` WHERE id = (.+) LIMIT 1$").
			WithArgs(pet1.ID).WillReturnRows(
			test_utils.GetJoinedTableCols(&Pet{}).
				AddRow(test_utils.GetStructValues(pet1)...),
		)
		pet := &Pet{}
		if err := pet.GetPet(db, 1); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, pet, pet1)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("Failed to meet expectations, got error: %v", err)
		}
	})
	t.Run("test get non-exist pet", func(t *testing.T) {
		mock.ExpectQuery("^SELECT (.+) FROM \\`pet\\` WHERE id = (.+) LIMIT 1$").
			WithArgs(pet1.ID).WillReturnError(gorm.ErrRecordNotFound)
		pet := &Pet{}
		if err := pet.GetPet(db, 1); err != nil {
			if !errors.ErrorIs(err, errors.NotFoundError) {
				t.Fatal(err)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("Failed to meet expectations, got error: %v", err)
		}
	})
}

func TestPutPet(t *testing.T) {
	mock, db, _ := test_utils.SetupMockDB()
	pet1 := &Pet{
		ID:   1,
		Name: "cat1",
		Type: int32(pet_pb.PetType_PET_TYPE_CAT),
	}
	t.Run("test put pet", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `pet` (.+) VALUES \\(\\?,\\?,\\?\\)").
			WithArgs(pet1.Name, pet1.Type, pet1.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		if err := pet1.PutPet(db); err != nil {
			t.Fatal(err)
		}
	})
}
