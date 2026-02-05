package contract

type Auth interface {
	ID() string
	Guest() bool
}
