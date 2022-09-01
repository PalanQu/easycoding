package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
)

func (s *service) putPet(
	ctx context.Context,
	req *pet_pb.PutPetRequest,
) (*pet_pb.PutPetResponse, error) {
	pet, err := s.DB.Pet.Create().SetName(req.Name).SetType(int8(req.PetType)).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &pet_pb.PutPetResponse{
		Pet: &pet_pb.Pet{
			PetId:   int32(pet.ID),
			Name:    pet.Name,
			PetType: pet_pb.PetType(pet.Type),
		},
	}, nil
}
