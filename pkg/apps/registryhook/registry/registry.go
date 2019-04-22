package registry

var registryUrl = "http://localhost:5000"

func Init(url string) {
	registryUrl = url
}

type DeleteStatus int

const (
	DeleteStatusDeleted DeleteStatus = iota
	DeleteStatusNotFound
	DeleteStatusUnknown
)

func (d DeleteStatus) String() string {
	switch d {
	case DeleteStatusDeleted:
		return "deleted"
	case DeleteStatusNotFound:
		return "not found"
	}
	return "unknown"
}
