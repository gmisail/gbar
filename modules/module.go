package modules

/*
 *	Modules should all have the ability to run and fetch
 *	a new value. The `Run()` function returns a map of variables
 *	which can be used from within templates.
 */
type Module interface {
	Run() map[string] interface{} 
}

