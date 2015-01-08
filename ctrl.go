package sscc

// Controller is an interface for providing access to pkg's functionalities.
type Controller interface {
	Procer
	Dbuser
	Searcher
}

// control is a default implementation of `Controller`.
type control struct {
	Procer
	Dbuser
	Searcher
}

// NewControl returns `Controller` implementation based on `Context`.
func NewControl(ctx *Context) Controller {
	return &control{
		ctx.procer(),
		ctx.dbuser(),
		ctx.searcher(),
	}
}

// Context is an initializing structure for `Controller` implementation.
type Context struct {
	Proc   Procer   // Proc is an implementation of `Procer`.
	Dbus   Dbuser   // Dbus is an implementation of `Dbuser`.
	Search Searcher // Search is an implementation of `Searcher`.
}

func (ctx *Context) procer() Procer {
	if ctx.Proc != nil {
		return ctx.Proc
	}
	return defaultProc
}

func (ctx *Context) dbuser() Dbuser {
	if ctx.Dbus != nil {
		return ctx.Dbus
	}
	return defaultDbus
}

func (ctx *Context) searcher() Searcher {
	if ctx.Search != nil {
		return ctx.Search
	}
	return defaultSearch
}
