package modules

/*
 *	All modules should have the ability to take a template from a block and
 * 	then compile it into a resultant string. For example:
* 	cpu.Compile("<cpu>") --> "5%"
 */
type Module interface {
	Compile(template) string
}