package uploading

import "fmt"

type Service interface {
	AddImages(...Img) error
}

// Repository provides access to beer repository.
type Repository interface {
	// AddBeer saves a given beer to the repository.
	AddImage(Img) error
	// GetAllBeers returns all beers saved in storage.
	//GetAllImages() []listing.Beer
	//Get Images by region and location
	GetImagesByLoc(region, loc string)
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
