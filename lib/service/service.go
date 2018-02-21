package service

import (
	"time"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/paulbellamy/ratecounter"
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
	counter    *ratecounter.AvgRateCounter
}

func New(config *Config) *Service {
	return &Service{
		store:      config.Store,
		downloader: config.Downloader,
		locker:     config.Locker,
		counter:    ratecounter.NewAvgRateCounter(60 * time.Second),
	}
}

var (
	StorePollTries           = 3
	MaxLoops                 = 2
	DefaultPollSleepInterval = 200 * time.Millisecond
)

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
			if attempt < MaxLoops-1 {
				return s.syncedPerform(key, imgBytes, t, attempt+1)
			}
		}
	}

	return s.instrumentedPerform(key, imgBytes, t)
}

func (s *Service) instrumentedPerform(key string, data []byte, t Transformation) ([]byte, error) {
	start := time.Now()

	data, err := s.perform(key, data, t)
	if err != nil {
		s.counter.Incr(time.Since(start).Nanoseconds())
	}

	return data, err
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
	for i := 0; i < StorePollTries; i++ {
		// sleep half of avg execution time, so there will be good
		time.Sleep(s.pollSleepInterval())

		if stored := s.store.Get(key); len(stored) > 0 {
			return stored
		}
	}

	return nil
}

func (s *Service) pollSleepInterval() time.Duration {
	cnt := s.counter.Hits()
	if cnt == 0 {
		return DefaultPollSleepInterval
	}

	return time.Duration(s.counter.Rate() / 2)
}
