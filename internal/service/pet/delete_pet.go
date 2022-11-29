package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
	"easycoding/pkg/errors"
)

func (s *service) deletePet(
	ctx context.Context,
	req *pet_pb.DeletePetRequest,
) (*pet_pb.DeletePetResponse, error) {
	if err := s.DB.Pet.DeleteOneID(int(req.PetId)).Exec(ctx); err != nil {
		return nil, errors.ErrInvalid(err)
	}
	return &pet_pb.DeletePetResponse{}, nil
}
