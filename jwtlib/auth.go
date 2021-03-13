package jwtlib

type Auth interface {
	Sign() (string, error)
	Verify(token string) (map[string]interface{}, error)
}
