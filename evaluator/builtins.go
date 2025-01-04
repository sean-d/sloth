package evaluator

import "github.com/sean-d/sloth/object"

/*
The most important part of this function is the call to Go’s len and the returning of a newly allocated object.Integer.
Besides that we have error checking that makes sure that we can’t call this function with the wrong number of arguments
or with an argument of an unsupported type.
*/
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
}
