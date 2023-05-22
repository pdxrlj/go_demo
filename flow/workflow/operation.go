package workflow

type workloadHandler func([]byte, map[string][]string) ([]byte, error)

type Operation interface {
	GetId() string
	Encode() []byte
	GetProperties() map[string][]string
	Execute([]byte, map[string]interface{}) ([]byte, error)
}
