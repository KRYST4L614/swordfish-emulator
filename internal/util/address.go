package util

func Addr[T any](t T) *T { return &t }
