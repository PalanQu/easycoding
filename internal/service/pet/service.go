package pet

import (
	"context"
	pet_pb "easycoding/api/pet"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	Logger *logrus.Logger
	DB     *gorm.DB
}

var _ pet_pb.PetStoreSvcServer = (*service)(nil)

func New(logger *logrus.Logger, db *gorm.DB) *service {
	return &service{
		Logger: logger,
		DB:     db,
	}
}

func (s *service) GetPet(
	ctx context.Context,
	req *pet_pb.GetPetRequest,
) (*pet_pb.GetPetResponse, error) {
	return s.getPet(ctx, req)
}

func (s *service) PutPet(
	ctx context.Context,
	req *pet_pb.PutPetRequest,
) (*pet_pb.PutPetResponse, error) {
	return s.putPet(ctx, req)
}

func (s *service) DeletePet(
	ctx context.Context,
	req *pet_pb.DeletePetRequest,
) (*pet_pb.DeletePetResponse, error) {
	return s.deletePet(ctx, req)
}