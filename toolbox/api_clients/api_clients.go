package api_clients

import "github.com/gilperopiola/obreros/core"

var _ core.InternalAPIs = (*APIClients)(nil)
var _ core.ExternalAPIs = (*APIClients)(nil)

type APIClients struct {
}

func NewAPIClients() *APIClients {
	return &APIClients{}
}
