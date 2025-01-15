package utils

import "context"

type ContextKey string

func RootContext() context.Context {
	return context.Background()
}
