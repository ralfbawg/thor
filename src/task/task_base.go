package task

type task interface {
	init() (bool, string)
}
