package uid

/*
Generator provides an interface to generate unique IDs of length of 64 bits.
The uniqueness validity in terms of time, their frequency granularity and distribution depends on the implementation.
*/
type Generator interface {
	Get() (uint64, error)
}
