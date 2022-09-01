package pet

import (
	"context"
	pet_pb "easycoding/api/pet"
)

func (s *service) deletePet(
	ctx context.Context,
	req *pet_pb.DeletePetRequest,
) (*pet_pb.DeletePetResponse, error) {
	s.DB.Pet.DeleteOneID(int(req.PetId))
	return &pet_pb.DeletePetResponse{}, nil
}
