package main

import (
	"github.com/google/uuid"
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Guid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func Ptr[T any](x T) *T {
	return &x
}

func Map[T any, U any](src []T, mapper func(T) U) []U {
	dest := make([]U, len(src))
	for k, v := range src {
		dest[k] = mapper(v)
	}
	return dest
}

func Min[T ~int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T ~int](a, b T) T {
	if a > b {
		return a
	}
	return b
}
