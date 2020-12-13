package notify

type Notifier interface {
	Start()
	Error(err error)
	Success()
}
