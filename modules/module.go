package modules

/*
 *	Modules should all have the ability to run and fetch
 *	a new value.
 */
type Module interface {
	Run() string
}
