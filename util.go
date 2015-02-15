package spotify

import "fmt"

// outerr converts ([]byte, error) pair into (string, error).
func outerr(out []byte, err error) (string, error) {
	o := "<no output>"
	if out != nil {
		o = string(out)
	}
	return o, err
}

// min returns the smaller of x or y.
func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf("[spotify]: "+format, args...)
}
