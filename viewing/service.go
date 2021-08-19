package viewing

import "fmt"

type Service interface {
	GetAllImages() ([]Img, error)
	//GetImagesByCriteria(filters ...string)
}

// Repository provides access to image repository.
type Repository interface {
	// AddBeer saves a given image to the repository.
	AddImage(Img) error
	// GetAllBeers returns all images saved in storage.
	GetAllImages() ([]Img, error)
	//Get Images by region and location
	GetImagesByCriteria(filters ...string) ([]Img, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAllImages() ([]Img, error) {
	imgList, err := s.r.GetAllImages()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return imgList, nil
}
