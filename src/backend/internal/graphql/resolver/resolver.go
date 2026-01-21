// Package resolver contains the GraphQL resolver implementations.
package resolver

import "backend/internal/domain"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the root resolver for the GraphQL schema.
// It holds dependencies needed by query and mutation resolvers.
type Resolver struct {
	userRepo          domain.UserRepository
	fileRepo          domain.FileRepository
	refLetterRepo     domain.ReferenceLetterRepository
}

// NewResolver creates a new Resolver with the given repositories.
func NewResolver(
	userRepo domain.UserRepository,
	fileRepo domain.FileRepository,
	refLetterRepo domain.ReferenceLetterRepository,
) *Resolver {
	return &Resolver{
		userRepo:          userRepo,
		fileRepo:          fileRepo,
		refLetterRepo:     refLetterRepo,
	}
}
