package slg

func mergeProps(base, extra map[string]interface{}) map[string]interface{} {
	if extra == nil {
		return base
	}
	if base == nil {
		return extra
	}
	out := make(map[string]interface{}, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
