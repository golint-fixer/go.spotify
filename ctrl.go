package spotify

// Controller describes operations available through package.
type Controller interface {
	Execer
	Dbuser
	Searcher
}

// control is a default implementation of `Controller`.
type control struct {
	Execer
	Dbuser
	Searcher
}

// NewControl returns `Controller` implementation based on `Context`.
func NewControl(ctx *Context) Controller {
	return &control{
		ctx.execer(),
		ctx.dbuser(),
		ctx.searcher(),
	}
}

// Context is an initializing structure for `Controller` implementation.
type Context struct {
	Exec   Execer   // Exec is an implementation of `Procer`.
	Dbus   Dbuser   // Dbus is an implementation of `Dbuser`.
	Search Searcher // Search is an implementation of `Searcher`.
	Name   string   // Name is a name of process.
}

// execer returns `Execer` implementation.
func (ctx *Context) execer() Execer {
	if ctx.Exec != nil {
		return ctx.Exec
	}
	return NewExecer(ctx)
}

// dbuser returns `Dbuser` implementation.
func (ctx *Context) dbuser() Dbuser {
	if ctx.Dbus != nil {
		return ctx.Dbus
	}
	return NewDbuser()
}

// searcher returns `Searcher` implementation.
func (ctx *Context) searcher() Searcher {
	if ctx.Search != nil {
		return ctx.Search
	}
	return newSearcher()
}

// name returns name of Spotify process for `dbuser`.
func (ctx *Context) name() string {
	if ctx.Name != "" {
		return ctx.Name
	}
	return "spotify"
}
