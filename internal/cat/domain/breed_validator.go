package domain

import "context"

type BreedValidator interface {
	ValidateBreed(ctx context.Context, breed string) (bool, error)
}
