// Package resolver contains the GraphQL resolver implementations.
package resolver

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain"
	"backend/internal/logger"
	"backend/internal/service"

	model "backend/internal/graphql/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the root resolver for the GraphQL schema.
// It holds dependencies needed by query and mutation resolvers.
type Resolver struct {
	userRepo              domain.UserRepository
	fileRepo              domain.FileRepository
	refLetterRepo         domain.ReferenceLetterRepository
	resumeRepo            domain.ResumeRepository
	profileRepo           domain.ProfileRepository
	profileExpRepo        domain.ProfileExperienceRepository
	profileEduRepo        domain.ProfileEducationRepository
	profileSkillRepo      domain.ProfileSkillRepository
	authorRepo            domain.AuthorRepository
	testimonialRepo       domain.TestimonialRepository
	skillValidationRepo   domain.SkillValidationRepository
	expValidationRepo     domain.ExperienceValidationRepository
	storage               domain.Storage
	jobEnqueuer           domain.JobEnqueuer
	documentExtractor     domain.DocumentExtractor
	materializationSvc    *service.MaterializationService
	log                   logger.Logger
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
	profileSkillRepo domain.ProfileSkillRepository,
	authorRepo domain.AuthorRepository,
	testimonialRepo domain.TestimonialRepository,
	skillValidationRepo domain.SkillValidationRepository,
	expValidationRepo domain.ExperienceValidationRepository,
	storage domain.Storage,
	jobEnqueuer domain.JobEnqueuer,
	documentExtractor domain.DocumentExtractor,
	materializationSvc *service.MaterializationService,
	log logger.Logger,
) *Resolver {
	return &Resolver{
		userRepo:              userRepo,
		fileRepo:              fileRepo,
		refLetterRepo:         refLetterRepo,
		resumeRepo:            resumeRepo,
		profileRepo:           profileRepo,
		profileExpRepo:        profileExpRepo,
		profileEduRepo:        profileEduRepo,
		profileSkillRepo:      profileSkillRepo,
		authorRepo:            authorRepo,
		testimonialRepo:       testimonialRepo,
		skillValidationRepo:   skillValidationRepo,
		expValidationRepo:     expValidationRepo,
		storage:               storage,
		jobEnqueuer:           jobEnqueuer,
		documentExtractor:     documentExtractor,
		materializationSvc:    materializationSvc,
		log:                   log,
	}
}

// loadProfileData fetches all nested data for a profile and returns the GraphQL model.
func (r *queryResolver) loadProfileData(ctx context.Context, profile *domain.Profile) (*model.Profile, error) {
	user, err := r.userRepo.GetByID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for profile: %w", err)
	}
	gqlUser := toGraphQLUser(user)

	experiences, err := r.profileExpRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiences for profile: %w", err)
	}
	gqlExperiences := toGraphQLProfileExperiences(experiences)

	educations, err := r.profileEduRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educations for profile: %w", err)
	}
	gqlEducations := toGraphQLProfileEducations(educations)

	skills, err := r.profileSkillRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills for profile: %w", err)
	}
	gqlSkills := toGraphQLProfileSkills(skills)

	var photoURL *string
	if profile.ProfilePhotoFileID != nil {
		photoFile, fileErr := r.fileRepo.GetByID(ctx, *profile.ProfilePhotoFileID)
		if fileErr == nil && photoFile != nil {
			url, urlErr := r.storage.GetPublicURL(ctx, photoFile.StorageKey, 24*time.Hour)
			if urlErr == nil {
				photoURL = &url
			}
		}
	}

	return toGraphQLProfile(profile, gqlUser, gqlExperiences, gqlEducations, gqlSkills, photoURL), nil
}
