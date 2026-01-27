// Package resolver contains the GraphQL resolver implementations.
package resolver

import (
	"backend/internal/domain"
	"backend/internal/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the root resolver for the GraphQL schema.
// It holds dependencies needed by query and mutation resolvers.
type Resolver struct {
	userRepo       domain.UserRepository
	fileRepo       domain.FileRepository
	refLetterRepo  domain.ReferenceLetterRepository
	resumeRepo     domain.ResumeRepository
	profileRepo    domain.ProfileRepository
	profileExpRepo domain.ProfileExperienceRepository
	profileEduRepo domain.ProfileEducationRepository
	storage        domain.Storage
	jobEnqueuer    domain.JobEnqueuer
	log            logger.Logger
}

// NewResolver creates a new Resolver with the given repositories.
func NewResolver(
	userRepo domain.UserRepository,
	fileRepo domain.FileRepository,
	refLetterRepo domain.ReferenceLetterRepository,
	resumeRepo domain.ResumeRepository,
	profileRepo domain.ProfileRepository,
	profileExpRepo domain.ProfileExperienceRepository,
	profileEduRepo domain.ProfileEducationRepository,
	storage domain.Storage,
	jobEnqueuer domain.JobEnqueuer,
	log logger.Logger,
) *Resolver {
	return &Resolver{
		userRepo:       userRepo,
		fileRepo:       fileRepo,
		refLetterRepo:  refLetterRepo,
		resumeRepo:     resumeRepo,
		profileRepo:    profileRepo,
		profileExpRepo: profileExpRepo,
		profileEduRepo: profileEduRepo,
		storage:        storage,
		jobEnqueuer:    jobEnqueuer,
		log:            log,
	}
}
