package cli

import "reflect"

// Option is the info of a Command option field.
type Option struct {
	Name  string
	Alias string
	Help  string
}

// setName sets the Option name from the given Command struct field.
func (o *Option) setName(field reflect.StructField) {
	name := field.Tag.Get("optName")
	if name != "" {
		o.Name = cmdOptionNamePrefix + name
	}
}

// setAlias sets the Option alias from the given Command struct field.
func (o *Option) setAlias(field reflect.StructField) {
	alias := field.Tag.Get("optAlias")
	if alias != "" {
		o.Alias = cmdOptionAliasPrefix + alias
	}
}

// setHelp sets the Option help from the given Command struct field.
func (o *Option) setHelp(field reflect.StructField) {
	help := field.Tag.Get("optHelp")
	if help != "" {
		o.Help = help
	}
}
