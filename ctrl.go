package sscc

// Controller is an interface for providing access to all functionalities
// of `sscc` pkg.
type Controller interface {
	Procer
	Dbuser
	Searcher
}

// control is type for controller of all sscc operations.
type control struct {
	*proc
	*dbus
	*web
}

// NewControl returns default instance of cotrol type.
func NewControl() (Controller, error) {
	c := &control{newProc(), newDbus(), newWeb()}
	return c, c.init()
}
