package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
	"easycoding/pkg/orm"
)

func (s *service) deletePet(
	ctx context.Context,
	req *pet_pb.DeletePetRequest,
) (*pet_pb.DeletePetResponse, error) {
	pet := &orm.Pet{}
	if err := pet.DeletePet(s.DB, req.PetId); err != nil {
		return nil, err
	}
	return &pet_pb.DeletePetResponse{}, nil
}
