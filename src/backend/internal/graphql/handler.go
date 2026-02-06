// Package graphql contains the GraphQL API implementation.
package graphql

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"backend/internal/domain"
	"backend/internal/graphql/generated"
	"backend/internal/graphql/resolver"
	"backend/internal/logger"
	"backend/internal/service"
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
) http.Handler {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver.NewResolver(userRepo, fileRepo, refLetterRepo, resumeRepo, profileRepo, profileExpRepo, profileEduRepo, profileSkillRepo, authorRepo, testimonialRepo, skillValidationRepo, expValidationRepo, storage, jobEnqueuer, documentExtractor, materializationSvc, log),
		}),
	)

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		log.Error("GraphQL error",
			logger.Feature("graphql"),
			logger.String("message", gqlErr.Message),
			logger.String("path", gqlErr.Path.String()),
		)
		return gqlErr
	})

	return srv
}

// NewPlaygroundHandler creates a new GraphQL Playground HTTP handler.
// The endpoint parameter specifies the GraphQL endpoint URL.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
