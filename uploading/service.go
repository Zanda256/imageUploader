package uploading

import (
	"fmt"
)

type Service interface {
	AddImages(...Img) error
}

// Repository provides access to image repository.
type Repository interface {
	// AddImage saves a given image to the repository.
	AddImage(Img) error
}

type service struct {
	r Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddImages(imgList ...Img) error {
	for _, v := range imgList {
		err := s.r.AddImage(v)
		if err != nil {
			fmt.Printf("failed to add image : %+v", err)
			return err
		}
	}
	return nil
}
