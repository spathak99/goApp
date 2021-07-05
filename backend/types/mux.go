package types

import(
	"net/http"
)
// Server Struct
type Multiplexer struct {
	multiplexer http.ServeMux
}
