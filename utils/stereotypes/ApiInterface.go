package stereotypes

type ApiRestInteface interface {
	GetRoot() string
	GetClientRoot() string
	GetEntityPath() string
	GetEntityId() string
}
