package repo

import "context"

type Repo interface {
	GetHashedPasswordByUsername(ctx context.Context, username string) (string, error)
}
