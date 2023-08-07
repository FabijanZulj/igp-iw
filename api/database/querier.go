// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error)
	CreateVerifyData(ctx context.Context, db DBTX, arg CreateVerifyDataParams) (Verifydatum, error)
	DeleteUser(ctx context.Context, db DBTX, id int32) error
	GetUserByEmail(ctx context.Context, db DBTX, email string) (User, error)
	GetVerificationCode(ctx context.Context, db DBTX, email string) (Verifydatum, error)
	ListUsers(ctx context.Context, db DBTX) ([]User, error)
	UpdateUser(ctx context.Context, db DBTX, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)