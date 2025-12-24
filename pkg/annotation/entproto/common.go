package entproto

const (
	fieldAnnotationName = "entproto_field_annotation"
	apiAnnotationName   = "entproto_api_annotation"
)

func bitTwiddling[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint](a T) []T {
	if a == 0 {
		return nil
	}
	vals := make([]T, 0)
	for v := a; v != 0; v &= v - 1 {
		lowest := v & -v
		vals = append(vals, lowest)
	}
	return vals
}
