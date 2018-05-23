package trigger

type Trigger interface {
	On(event string, task interface{}) error
	Fire(event string, params ...interface{}) (error)
	FireBackground(event string, params ...interface{}) (error)
	Clear(event string) error
	ClearEvents()
	HasEvent(event string) bool
	Events() map[string]string
	EventCount() int
}

var defaultTrigger = New()

// Default global trigger options.
func On(event string, task interface{}) error {
	return defaultTrigger.On(event, task)
}

func Fire(event string, params ...interface{}) (error) {
	return defaultTrigger.Fire(event, params...)
}

func FireBackground(event string, params ...interface{}) (error) {
	return defaultTrigger.FireBackground(event, params...)
}

func Clear(event string) error {
	return defaultTrigger.Clear(event)
}

func ClearEvents() {
	defaultTrigger.ClearEvents()
}

func HasEvent(event string) bool {
	return defaultTrigger.HasEvent(event)
}

func Events() map[string]string {
	return defaultTrigger.Events()
}

func EventCount() int {
	return defaultTrigger.EventCount()
}
