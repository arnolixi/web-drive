package arc

type IClass interface {
	Build(reactor *Reactor)
	Name() string
}
