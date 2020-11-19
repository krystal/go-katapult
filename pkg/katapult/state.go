package katapult

type State string

const (
	DraftState    State = "draft"
	FailedState   State = "failed"
	PendingState  State = "pending"
	CompleteState State = "complete"
	BuildingState State = "building"
)
