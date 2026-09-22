package model

// AuthResult describes the credential an auth utility endpoint accepted.
type AuthResult struct {
	Scheme AuthScheme
	// Credential is the username, token or API key that was accepted.
	Credential string
}
