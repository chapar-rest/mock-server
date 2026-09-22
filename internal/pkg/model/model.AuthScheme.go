package model

// AuthScheme is a credential scheme accepted by the auth utility endpoints.
type AuthScheme string

const (
	AuthSchemeBasic  AuthScheme = "basic"
	AuthSchemeBearer AuthScheme = "bearer"
	AuthSchemeApiKey AuthScheme = "api_key"
)
