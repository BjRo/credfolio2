// Package graphql contains the GraphQL API implementation.
package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"backend/internal/domain"
	"backend/internal/graphql/generated"
	"backend/internal/graphql/resolver"
	"backend/internal/logger"
)

// NewHandler creates a new GraphQL HTTP handler with the given repositories.
func NewHandler(
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
) http.Handler {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver.NewResolver(userRepo, fileRepo, refLetterRepo, resumeRepo, profileRepo, profileExpRepo, profileEduRepo, storage, jobEnqueuer, log),
		}),
	)
	return srv
}

// NewPlaygroundHandler creates a new GraphQL Playground HTTP handler.
// The endpoint parameter specifies the GraphQL endpoint URL.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
