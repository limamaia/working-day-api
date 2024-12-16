package services

type HealthService struct{}

func (s *HealthService) GetStatus() string {
	return "App it's alive!"
}
