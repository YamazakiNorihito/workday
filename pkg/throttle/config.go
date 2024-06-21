package throttle

type Config struct {
	BatchSize int
	Sleep     func()
}
