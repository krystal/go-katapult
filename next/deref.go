package next

import "reflect"

// DerefOrEmpty returns the dereferenced value of the input pointer, or the zero
// value of the type if the input is nil.
func DerefOrEmpty[T any](in *T) T {
	if in == nil {
		tType := reflect.TypeOf(*new(T))

		//nolint:exhaustive
		switch tType.Kind() {
		case reflect.Slice, reflect.Array:
			return reflect.MakeSlice(tType, 0, 0).Interface().(T)
		case reflect.Map:
			return reflect.MakeMap(tType).Interface().(T)
		default:
			var empty T

			return empty
		}
	}

	return *in
}
