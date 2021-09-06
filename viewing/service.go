package viewing

import "fmt"

type Service interface {
	GetAllImages() ([]Img, error)
	GetImagesByCriteria(reg, loc string) ([]Img, error)
}

// Repository provides access to image repository.
type Repository interface {
	// GetAllImages returns all images saved in storage.
	GetAllImages() ([]Img, error)
	//Get Images by region and location
	GetImagesByCriteria(region, location string) ([]Img, error)
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

func (s *service) GetImagesByCriteria(region, location string) ([]Img, error) {
	result, err := s.r.GetImagesByCriteria(region, location)
	if err != nil {
		return nil, err
	}
	return result, nil
}
