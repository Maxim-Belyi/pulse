package tracer

type AMQPCarrier map[string]any

func (c AMQPCarrier) Set(key string, value string) {
	c[key] = value
}

func (c AMQPCarrier) Get(key string) string {
	if strVal, ok := c[key].(string); ok {
		return strVal
	}
	return ""
}
func (c AMQPCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
