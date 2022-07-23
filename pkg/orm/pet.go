package orm

import (
	"easycoding/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	ID   int32 `gorm:"primaryKey;autoIncrement"`
	Name string
	// TODO(qujiabao): replace int32 to pet_pb.PetType, because of `sqlize`
	Type      int32
	Age       int32
	CreatedAt time.Time `gorm:"default:now()"`
}

var _ Model = (*Pet)(nil)

func (Pet) TableName() string {
	return "pet"
}

func (pet *Pet) GetPet(db *gorm.DB, id int32) error {
	if err := db.Take(pet, "id = ?", id).Error; err != nil {
		if errors.ErrorIs(err, gorm.ErrRecordNotFound) {
			return errors.ErrNotFound(err)
		}
		return errors.ErrInternal(err)
	}
	return nil
}

func (pet *Pet) DeletePet(db *gorm.DB, id int32) error {
	if err := db.Where("id = ?", id).Delete(pet).Error; err != nil {
		if errors.ErrorIs(err, gorm.ErrRecordNotFound) {
			return errors.ErrNotFound(err)
		}
		return errors.ErrInternal(err)
	}
	return nil
}

func (pet *Pet) PutPet(db *gorm.DB) error {
	if err := db.Create(pet).Error; err != nil {
		return errors.ErrInvalid(err)
	}
	return nil
}
