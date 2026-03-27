package core

type TrackerService interface {
	Track(shortPath string) error
}

type TrackerPort interface {
	Track(shortPath string) error
}

type trackerService struct {
	port TrackerPort
}

func NewTrackerService(port TrackerPort) TrackerService {
	return &trackerService{port: port}
}

func (s *trackerService) Track(shortPath string) error {
	return s.port.Track(shortPath)
}
