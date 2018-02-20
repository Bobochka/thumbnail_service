package lib

type Mutex interface {
	Lock() error
	Unlock() bool
}
