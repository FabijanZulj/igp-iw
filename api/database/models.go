// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import ()

type User struct {
	ID         int32
	Email      string
	Password   string
	Isverified bool
}

type Verifydatum struct {
	Userid int32
	Code   string
}