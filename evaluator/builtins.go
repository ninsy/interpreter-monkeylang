package evaluator

import (
	"errors"
	"monkey/object"
)

var builtinMethods = map[string]*object.BuiltinMethod{
	"parseInt": &object.BuiltinMethod{
		Fn: func(args ...object.Object) object.Object {
			return NULL
		},
	},
	"len": &object.BuiltinMethod{
		Fn: func(args ...object.Object) object.Object {
			return length.Fn(args...)
		},
	},
	"slice": &object.BuiltinMethod{
		Fn: func(args ...object.Object) object.Object {
			return slice.Fn(args...)
		},
	},
	"head": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return head.Fn(args...)			
		},
	},
	"tail": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return tail.Fn(args...)			
		},
	},
	"push": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return push.Fn(args...)			
		},
	},
	"map": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return mapFn.Fn(args...)			
		},
	},
	"reduce": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return reduce.Fn(args...)			
		},
	},
	"filter": &object.BuiltinMethod {
		Fn: func(args ...object.Object) object.Object {
			return filter.Fn(args...)			
		},
	},
}

var slice = &object.BuiltinMethod{
	Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 || len(args) > 3 {
			return newError(`slice(arr, begin, *end) parameters are:
				arr: array on which slice is performed,
				begin: zero-based index at which to begin extraction,
				end: (optional) Zero-based index before which to end extraction. slice extracts up to but not including end
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("slice() only supports arrays, got=%s", args[0].Type())
		}

		output := &object.Array{}
		var beginIdx, endIdx int64
		var okBegin, okEnd error

		if len(args) == 2 {
			beginIdx, okBegin = extractSliceIndex(args[1])
			if okBegin != nil {
				return newError("parameters 'begin' and 'end' of slice() must be integers!")					
			}
		} else {
			beginIdx, okBegin = extractSliceIndex(args[1])
			endIdx, okEnd = extractSliceIndex(args[2])
			if okBegin != nil || okEnd != nil {
				return newError("parameters 'begin' and 'end' of slice() must be integers!")
			}
		}

		for idx := beginIdx; idx < endIdx && idx < int64(len(arr.Elements)); idx++ {
			output.Elements = append(output.Elements, arr.Elements[idx])
		}

		return output
	},
}

var mapFn = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError(`map(arr, fn) parameters are:
				arr: array on which map is performed,
				fn: zero-based index at which to begin extraction,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("map() only supports arrays, got=%s", args[0].Type())
		}
		fn, ok := args[1].(*object.Function)
		if !ok {
			return newError("fn must be function, got=%s", args[1].Type())
		}

		output := make([]object.Object, len(arr.Elements), len(arr.Elements))

		for i, mapArg := range arr.Elements {
			output[i] = applyFunction(fn, []object.Object{mapArg})
		}

		return &object.Array{Elements: output}
	},
}

var reduce = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 || len(args) > 3 {
			return newError(`reduce(arr, fn, initial) parameters are:
				arr: array on which slice is performed,
				fn: zero-based index at which to call reducer,
				initial: Value to use as the first argument to the first call of the callback. If no initial value is supplied, the first element in the array will be used.
				 				 Calling reduce() on an empty array without an initial value is an error,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("reduce() only supports arrays, got=%s", args[0].Type())
		}
		fn, ok := args[1].(*object.Function)
		if !ok {
			return newError("fn must be functon, got=%s", args[1].Type())
		}
		
		var initialVal object.Object

		if len(args) == 2 {
			if (len(arr.Elements) == 0) {
				return newError("If reduce() hasn't received initial value, provided array musn't be empty", args[0].Type())
			}
			initialVal = arr.Elements[0]
		} else {
			initialVal = args[2]
		}

		switch initialVal.Type() {
		case object.INTEGER_OBJ:
			integer := initialVal.(*object.Integer)
			if len(args) == 2 {
				for i := 1; i < len(arr.Elements) ; i++ {
					integer = applyFunction(fn, []object.Object{integer, arr.Elements[i]}).(*object.Integer)
				}
			} else {
				for _, mapArg := range arr.Elements {
					integer = applyFunction(fn, []object.Object{integer, mapArg}).(*object.Integer)
				}
			}
			return integer
		default:
			return newError("Unsupported type %s of initial value", initialVal.Type())
		}
	},
}

var filter = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError(`filter(arr, fn) parameters are:
				arr: array on which filtering is performed,
				fn: zero-based index at which to begin extraction,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("map() only supports arrays, got=%s", args[0].Type())
		}
		fn, ok := args[1].(*object.Function)
		if !ok {
			return newError("fn must be function, got=%s", args[1].Type())
		}

		output := []object.Object{}

		for _, mapArg := range arr.Elements {
			result := applyFunction(fn, []object.Object{mapArg}).(*object.Boolean)
			if result.Value {
				output = append(output, mapArg)
			}
		}

		return &object.Array{Elements: output}
	},
}

var length = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("len() accepts single parameter, got=%d", len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value)) }
		default:
			return newError("len() accepts only	strings / arrays, got=%s", args[0].Type())
		}
	},
}

var head = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError(`head(arr) parameters are:
				arr: array on which head is performed,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("head() only supports arrays, got=%s", args[0].Type())
		}
		return slice.Fn(arr, &object.Integer{Value:0}, &object.Integer{Value:1})
	},
}

var tail = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError(`tail(arr) parameters are:
				arr: array on which tail is performed,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("tail() only supports arrays, got=%s", args[0].Type())
		}
		return slice.Fn(args[0], &object.Integer{Value:1}, &object.Integer{Value:int64(len(arr.Elements))})
	},
}

var push = &object.BuiltinMethod {
	Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError(`push(arr, ...elements) parameters are:
				arr: array on which slice is performed,
				elements:The elements to add to the end of the array,
			`)
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("push() only supports arrays, got=%s", args[0].Type())
		}
		totalLen := len(arr.Elements) + len(args) - 1
		newElements := make([]object.Object, totalLen, totalLen)
		copy(newElements, arr.Elements)
	
		idy := 1
		for idx := len(arr.Elements); idx < totalLen; idx++ {
			newElements[idx] = args[idy]
			idy++	
		}

		return &object.Array{Elements: newElements}
	},
}

func extractSliceIndex(idx object.Object) (int64, error) {
	switch idx.Type() {
	case object.INTEGER_OBJ:
		integer := idx.(*object.Integer)
		return int64(integer.Value), nil
	default:
		return -1, errors.New("slice() 'start' / 'end' must be integers or resolve to integers")
	}
}