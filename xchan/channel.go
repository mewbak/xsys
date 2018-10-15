package xchan

import "github.com/pkg/errors"

// TODO: check whether safe when multi threads close channel simutaneously?
// Safe close chan with wrap function
func SafeCloseChanStruct(ch chan struct{}) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("BrokenPipe")
		}
	}()

	// assume ch != nil here.
	close(ch) // panic if ch is closed
	return nil
}