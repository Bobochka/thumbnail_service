package service

import (
	"time"

	"github.com/Bobochka/thumbnail_service/lib"
)

type Store interface {
	Get(key string) []byte
	Set(key string, data []byte) error
}

type Transformation interface {
	Fingerprint(data []byte) string
	Perform(data []byte) ([]byte, error)
}

type Downloader interface {
	Download(url string) ([]byte, error)
}

type Locker interface {
	NewMutex(name string) lib.Mutex
}

type Config struct {
	Store      Store
	Downloader Downloader
	Locker     Locker
}

type Service struct {
	store      Store
	downloader Downloader
	locker     Locker
}

func New(config *Config) *Service {
	return &Service{
		store:      config.Store,
		downloader: config.Downloader,
		locker:     config.Locker,
	}
}

func (s *Service) Perform(url string, t Transformation) ([]byte, error) {
	imgBytes, err := s.downloader.Download(url)
	if err != nil {
		return nil, err
	}

	key := t.Fingerprint(imgBytes)

	if stored := s.store.Get(key); len(stored) > 0 {
		return stored, nil
	}

	return s.syncedPerform(key, imgBytes, t, 0)
}

const maxAttempts = 2

func (s *Service) syncedPerform(key string, imgBytes []byte, t Transformation, attempt int) ([]byte, error) {
	m := s.locker.NewMutex(key)

	err := m.Lock()

	if err == nil {
		defer m.Unlock()

		if stored := s.store.Get(key); len(stored) > 0 {
			return stored, nil
		}
	} else {
		value := s.pollStoredValue(key)

		if len(value) > 0 {
			return value, nil
		} else {
			if attempt < maxAttempts-1 {
				return s.syncedPerform(key, imgBytes, t, attempt+1)
			}
		}
	}

	return s.perform(key, imgBytes, t)
}

func (s *Service) perform(key string, data []byte, t Transformation) ([]byte, error) {
	res, err := t.Perform(data)
	if err != nil {
		return []byte{}, err
	}

	_ = s.store.Set(key, res)

	return res, err
}

func (s *Service) pollStoredValue(key string) []byte {
	// TODO: extract 3 to config
	for i := 0; i < 3; i++ {
		time.Sleep(200 * time.Millisecond) // TODO: config
		if stored := s.store.Get(key); len(stored) > 0 {
			return stored
		}
	}

	return nil
}
