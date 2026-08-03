package noop

type Noop struct{}

func New() (*Noop, error) {
	return &Noop{}, nil
}

func (n *Noop) HashFn(source []byte, salt []byte) ([]byte, error) {
	// TODO: consider whether we are effectively copying the source
	// so that it behaves more like other HashFn (i.e. prevents
	// modifications to the original.)
	return []byte(source), nil
}

func (n *Noop) HashFunc(source string, salt []byte) ([]byte, error) {
	// TODO: consider whether we are effectively copying the source
	// so that it behaves more like other HashFn (i.e. prevents
	// modifications to the original.)
	return []byte(source), nil
}
