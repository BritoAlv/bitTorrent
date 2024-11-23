package common

import "errors"

func CastTo[T interface{}](obj interface{}) (T, error) {
	casted, isExpectedType := obj.(T)

	if !isExpectedType {
		return casted, errors.New("cannot cast the object to the given type")
	}

	return casted, nil
}
