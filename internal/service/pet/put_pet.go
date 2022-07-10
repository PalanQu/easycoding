package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
	"easycoding/pkg/orm"
)

func (s *service) putPet(
	ctx context.Context,
	req *pet_pb.PutPetRequest,
) (*pet_pb.PutPetResponse, error) {
	pet := &orm.Pet{
		Name: req.Name,
		Type: int32(req.PetType),
	}
	if err := pet.PutPet(s.DB); err != nil {
		return nil, err
	}
	return &pet_pb.PutPetResponse{
		Pet: &pet_pb.Pet{
			PetId:   pet.ID,
			Name:    pet.Name,
			PetType: pet_pb.PetType(pet.Type),
		},
	}, nil
}
