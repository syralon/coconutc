package entproto

type PaginatorStyle int32

const (
	ClassicalPaginator PaginatorStyle = iota
	InfinitePaginator
)

func (s PaginatorStyle) String() string {
	switch s {
	case InfinitePaginator:
		return "InfinitePaginator"
	default:
		return "ClassicalPaginator"
	}
}
