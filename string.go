package goboot

func String(s string) *string {
	return &s
}

func StringValue(s *string) string {
	return *s
}
