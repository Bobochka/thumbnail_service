package service

import (
	"time"

	"log"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/go-errors/errors"
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
	ErrOnStore               = errors.New("unable to store processed data")
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
	isLocked := m.Lock() == nil

	defer func() {
		if r := recover(); r != nil {
			defer m.Unlock()
			panic(r)
		}
	}()

	if !isLocked {
		value := s.pollStoredValue(key)

		if len(value) > 0 {
			return value, nil
		} else {
			if attempt < MaxLoops-1 {
				return s.syncedPerform(key, imgBytes, t, attempt+1)
			}
		}
	}

	data, err := s.instrumentedPerform(key, imgBytes, t)

	// Just unlocking the lock after the job is done, will result in stampede:
	// if 2 goroutines are performing same request,
	// one can finish the the job and release the mutex,
	// yet another that already checked store on lines 66:68 will take mutex
	// and perform the job again which is unwanted.
	//
	// In case there's no error (meaning data was saved to store successfully),
	// instead of unlocking, mutex is extended to the time that is most definitely more
	// than the time between checking store and mutex acquire.
	//
	// In case of error (either wasn't able to perform transformation or got error writing to store)
	// need to unlock the mutex, so that next process will try:
	// strictly speaking this is a violation of "perform only once" rule, but probably it's better
	// to perform same job more than once, than return error in case store is not accessible.
	//
	// if panic will happen during execution, just unlock should work fine.
	//
	// Although this solution is not exactly fault tolerant, if process will crash in the middle of operation,
	// next process (because it won't acquire mutex) will spend (3 * avg (transform + store) time) on store polling
	// and then the state will be ok. During the burst, `processed only once` rule might be violated as well.

	if isLocked {
		if err == nil {
			m.Extend()
		} else {
			m.Unlock()
		}
	}

	// swallow store error
	if err == ErrOnStore {
		err = nil
	}

	return data, err
}

func (s *Service) instrumentedPerform(key string, data []byte, t Transformation) ([]byte, error) {
	start := time.Now()

	data, err := s.perform(key, data, t)
	if err == nil {
		s.counter.Incr(time.Since(start).Nanoseconds())
	}

	return data, err
}

func (s *Service) perform(key string, data []byte, t Transformation) ([]byte, error) {
	res, err := t.Perform(data)
	if err != nil {
		return []byte{}, err
	}

	err = s.store.Set(key, res)
	if err != nil {
		log.Println("error writing data to store: ", err)
		return res, ErrOnStore
	}

	return res, nil
}

func (s *Service) pollStoredValue(key string) []byte {
	for i := 0; i < StorePollTries; i++ {
		// sleep half of avg execution time each round,
		// so there will be good chance that concurrent performer is done
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
