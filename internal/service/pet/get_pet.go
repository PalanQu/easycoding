package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
	"easycoding/pkg/errors"
)

func (s *service) getPet(
	ctx context.Context,
	req *pet_pb.GetPetRequest,
) (*pet_pb.GetPetResponse, error) {
	pet, err := s.DB.Pet.Get(ctx, int(req.PetId))
	if err != nil {
		return nil, errors.WithMessage(err, "get pet failed")
	}
	return &pet_pb.GetPetResponse{
		Pet: &pet_pb.Pet{
			PetId:   int32(pet.ID),
			Name:    pet.Name,
			PetType: pet_pb.PetType(pet.Type),
		},
	}, nil
}
