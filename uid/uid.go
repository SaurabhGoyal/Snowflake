package uid

type UIDGenerator interface {
	Get() (uint64, error)
}
