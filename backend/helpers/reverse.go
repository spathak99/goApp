
package helpers

import (
	"goApp/backend/types"
)


// Reverse a list
func Reverse(source []types.Post) []types.Post {
	destination := make([]types.Post, len(source))
	for i, j := 0, len(source)-1; i <= j; i, j = i+1, j-1 {
		destination[i], destination[j] = source[j], source[i]
	}
	return destination
}
