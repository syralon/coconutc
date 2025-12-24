package entproto

type Option func(*options)

func (fn Option) applyService(sb *ServiceBuilder) {
	fn(&sb.options)
}

func (fn Option) applyEnt(eb *EntBuilder) {
	fn(&eb.options)
}

func (fn Option) applyGenerator(g *Generator) {
	fn(&g.options)
}

//func (fn Option) applyEvent(e *EventBuilder) {
//	fn(&e.options)
//}

func WithProtoPackage(pkg string) Option {
	return func(b *options) {
		b.protoPackage = pkg
	}
}
func WithGoPackage(pkg string) Option {
	return func(b *options) {
		b.goPackage = pkg
	}
}

func WithPath(path string) Option {
	return func(b *options) {
		b.path = path
	}
}

// UsePOSTInList use POST method instead of GET in List function.
func UsePOSTInList(v bool) Option {
	return func(b *options) {
		b.usePOSTInList = v
	}
}

type options struct {
	protoPackage  string
	goPackage     string
	path          string
	usePOSTInList bool
}
