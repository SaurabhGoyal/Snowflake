package uid

type UID uint64

type UIDGenerator interface {
	Get() (uint64, error)
}
