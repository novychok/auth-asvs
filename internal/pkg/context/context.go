package context

import "context"

func New() context.Context {
	return context.Background()
}
