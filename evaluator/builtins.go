package evaluator

import (
	"monkey/object"
)

var builtinMethods = map[string]*object.BuiltinMethod{
	"parseInt": &object.BuiltinMethod{
		Fn: func(args ...object.Object) object.Object {
			return NULL
		},
	},
}
