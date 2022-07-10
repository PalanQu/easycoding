package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
	"easycoding/pkg/orm"
)

func (s *service) getPet(
	ctx context.Context,
	req *pet_pb.GetPetRequest,
) (*pet_pb.GetPetResponse, error) {
	pet := &orm.Pet{}
	if err := pet.GetPet(s.DB, req.PetId); err != nil {
		return nil, err
	}

	return &pet_pb.GetPetResponse{
		Pet: &pet_pb.Pet{
			PetId:   pet.ID,
			Name:    pet.Name,
			PetType: pet_pb.PetType(pet.Type),
		},
	}, nil
}
