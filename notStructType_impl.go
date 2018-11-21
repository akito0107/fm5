// Code generated by "generr"; DO NOT EDIT.
package fm5

import "fmt"

func IsNotStructType(err error) (bool, string) {
	var typename string
	if e, ok := err.(notStructType); ok {
		typename = e.NotStructType()
		return true, typename
	}
	return false, typename
}

type NotStructType struct {
	Typename string
}

func (e *NotStructType) NotStructType() string {
	return e.Typename
}
func (e *NotStructType) Error() string {
	return fmt.Sprintf("notStructType Typename: %v", e.Typename)
}
