package task

//nolint:revive
type TaskStatus int8

const (
	TaskStatusIncomplete TaskStatus = iota // task is incomplete
	TaskStatusCompleted                    // task is completed
)

func (s TaskStatus) Valid() bool {
	return s == TaskStatusIncomplete || s == TaskStatusCompleted
}
