package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"backend/internal/domain"
)

// MaterializationResult contains counts of materialized items from resume extraction.
type MaterializationResult struct {
	Experiences int
	Educations  int
	Skills      int
}

// ReferenceLetterMaterializationResult contains counts of materialized items from reference letter extraction.
type ReferenceLetterMaterializationResult struct {
	Testimonials int
}

// MaterializationService handles materializing extracted data into profile tables.
// It is used by both the ResumeProcessingWorker (auto-materialization) and the
// importDocumentResults mutation (user-confirmed import).
type MaterializationService struct {
	profileRepo      domain.ProfileRepository
	profileExpRepo   domain.ProfileExperienceRepository
	profileEduRepo   domain.ProfileEducationRepository
	profileSkillRepo domain.ProfileSkillRepository
	authorRepo       domain.AuthorRepository
	testimonialRepo  domain.TestimonialRepository
}

// NewMaterializationService creates a new MaterializationService.
func NewMaterializationService(
	profileRepo domain.ProfileRepository,
	profileExpRepo domain.ProfileExperienceRepository,
	profileEduRepo domain.ProfileEducationRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	authorRepo domain.AuthorRepository,
	testimonialRepo domain.TestimonialRepository,
) *MaterializationService {
	return &MaterializationService{
		profileRepo:      profileRepo,
		profileExpRepo:   profileExpRepo,
		profileEduRepo:   profileEduRepo,
		profileSkillRepo: profileSkillRepo,
		authorRepo:       authorRepo,
		testimonialRepo:  testimonialRepo,
	}
}

// MaterializeResumeData creates profile education, experience, and skill rows from extracted resume data.
// This makes profile tables the single source of truth for display.
// Each category (experiences, education, skills) is processed independently so that failures in one
// category don't prevent the others from being saved.
// Idempotent: deletes existing entries from the same resume before re-creating.
func (s *MaterializationService) MaterializeResumeData(
	ctx context.Context,
	resumeID uuid.UUID,
	userID uuid.UUID,
	data *domain.ResumeExtractedData,
) (*MaterializationResult, error) {
	// Get or create the user's profile
	profile, err := s.profileRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create profile: %w", err)
	}

	// Delete any existing entries from this resume (idempotent re-processing)
	if delErr := s.profileExpRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return nil, fmt.Errorf("failed to delete existing experiences for resume: %w", delErr)
	}
	if delErr := s.profileEduRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return nil, fmt.Errorf("failed to delete existing education for resume: %w", delErr)
	}
	if delErr := s.profileSkillRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return nil, fmt.Errorf("failed to delete existing skills for resume: %w", delErr)
	}

	result := &MaterializationResult{}
	var errors []error

	if expCount, expErr := s.materializeExperiences(ctx, resumeID, profile.ID, data.Experience); expErr != nil {
		errors = append(errors, fmt.Errorf("experiences: %w", expErr))
	} else {
		result.Experiences = expCount
	}

	if eduCount, eduErr := s.materializeEducation(ctx, resumeID, profile.ID, data.Education); eduErr != nil {
		errors = append(errors, fmt.Errorf("education: %w", eduErr))
	} else {
		result.Educations = eduCount
	}

	if skillCount, skillErr := s.materializeSkills(ctx, resumeID, profile.ID, data.Skills); skillErr != nil {
		errors = append(errors, fmt.Errorf("skills: %w", skillErr))
	} else {
		result.Skills = skillCount
	}

	if len(errors) > 0 {
		return result, fmt.Errorf("materialization errors: %v", errors)
	}

	return result, nil
}

func (s *MaterializationService) materializeExperiences(ctx context.Context, resumeID, profileID uuid.UUID, experiences []domain.WorkExperience) (int, error) {
	displayOrder, err := s.profileExpRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to get next experience display order: %w", err)
	}

	for i, exp := range experiences {
		originalJSON, marshalErr := json.Marshal(exp)
		if marshalErr != nil {
			return i, fmt.Errorf("failed to marshal experience original data: %w", marshalErr)
		}

		profileExp := &domain.ProfileExperience{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Company:        exp.Company,
			Title:          exp.Title,
			Location:       exp.Location,
			StartDate:      exp.StartDate,
			EndDate:        exp.EndDate,
			IsCurrent:      exp.IsCurrent,
			Description:    exp.Description,
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
			OriginalData:   originalJSON,
		}
		if createErr := s.profileExpRepo.Create(ctx, profileExp); createErr != nil {
			return i, fmt.Errorf("failed to create experience for %s at %s: %w", exp.Title, exp.Company, createErr)
		}
	}
	return len(experiences), nil
}

func (s *MaterializationService) materializeEducation(ctx context.Context, resumeID, profileID uuid.UUID, educations []domain.Education) (int, error) {
	displayOrder, err := s.profileEduRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to get next education display order: %w", err)
	}

	for i, edu := range educations {
		originalJSON, marshalErr := json.Marshal(edu)
		if marshalErr != nil {
			return i, fmt.Errorf("failed to marshal education original data: %w", marshalErr)
		}

		degree := "Degree"
		if edu.Degree != nil && *edu.Degree != "" {
			degree = *edu.Degree
		}

		profileEdu := &domain.ProfileEducation{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Institution:    edu.Institution,
			Degree:         degree,
			Field:          edu.Field,
			StartDate:      edu.StartDate,
			EndDate:        edu.EndDate,
			Description:    edu.Achievements, // Map achievements -> description
			GPA:            edu.GPA,
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
			OriginalData:   originalJSON,
		}
		if createErr := s.profileEduRepo.Create(ctx, profileEdu); createErr != nil {
			return i, fmt.Errorf("failed to create education for %s: %w", edu.Institution, createErr)
		}
	}
	return len(educations), nil
}

func (s *MaterializationService) materializeSkills(ctx context.Context, resumeID, profileID uuid.UUID, skills []string) (int, error) {
	displayOrder, err := s.profileSkillRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to get next skill display order: %w", err)
	}

	dedupedSkills := DeduplicateSkills(skills)

	for i, skillName := range dedupedSkills {
		profileSkill := &domain.ProfileSkill{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Name:           skillName,
			NormalizedName: strings.ToLower(skillName),
			Category:       "TECHNICAL",
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
		}
		if createErr := s.profileSkillRepo.CreateIgnoreDuplicate(ctx, profileSkill); createErr != nil {
			return i, fmt.Errorf("failed to create skill %q: %w", skillName, createErr)
		}
	}
	return len(dedupedSkills), nil
}

// MaterializeReferenceLetterData creates author and testimonial rows from extracted reference letter data.
// Idempotent: deletes existing testimonials for the same reference letter before re-creating.
func (s *MaterializationService) MaterializeReferenceLetterData(
	ctx context.Context,
	referenceLetterID uuid.UUID,
	userID uuid.UUID,
	data *domain.ExtractedLetterData,
) (*ReferenceLetterMaterializationResult, error) {
	result := &ReferenceLetterMaterializationResult{}

	if len(data.Testimonials) == 0 {
		return result, nil
	}

	// Get or create the user's profile
	profile, err := s.profileRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create profile: %w", err)
	}

	// Delete existing testimonials for this reference letter (idempotent re-processing)
	if delErr := s.testimonialRepo.DeleteByReferenceLetterID(ctx, referenceLetterID); delErr != nil {
		return nil, fmt.Errorf("failed to delete existing testimonials for reference letter: %w", delErr)
	}

	// Find or create the Author entity
	var author *domain.Author
	existingAuthor, findErr := s.authorRepo.FindByNameAndCompany(ctx, profile.ID, data.Author.Name, data.Author.Company)
	if findErr != nil {
		return nil, fmt.Errorf("failed to find existing author: %w", findErr)
	}

	if existingAuthor != nil {
		author = existingAuthor
	} else {
		author = &domain.Author{
			ProfileID: profile.ID,
			Name:      data.Author.Name,
			Title:     data.Author.Title,
			Company:   data.Author.Company,
		}
		if createErr := s.authorRepo.Create(ctx, author); createErr != nil {
			return nil, fmt.Errorf("failed to create author: %w", createErr)
		}
	}

	// Create testimonials
	relationship := mapAuthorToTestimonialRelationship(data.Author.Relationship)
	for _, extractedTestimonial := range data.Testimonials {
		testimonial := &domain.Testimonial{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			ReferenceLetterID: referenceLetterID,
			AuthorID:          &author.ID,
			Quote:             extractedTestimonial.Quote,
			Relationship:      relationship,
			SkillsMentioned:   extractedTestimonial.SkillsMentioned,
			AuthorName:        &data.Author.Name,
			AuthorTitle:       data.Author.Title,
			AuthorCompany:     data.Author.Company,
		}
		if createErr := s.testimonialRepo.Create(ctx, testimonial); createErr != nil {
			return nil, fmt.Errorf("failed to create testimonial: %w", createErr)
		}
		result.Testimonials++
	}

	return result, nil
}

// mapAuthorToTestimonialRelationship converts an extraction-level author relationship
// to a testimonial relationship type.
func mapAuthorToTestimonialRelationship(ar domain.AuthorRelationship) domain.TestimonialRelationship {
	switch ar {
	case domain.AuthorRelationshipManager:
		return domain.TestimonialRelationshipManager
	case domain.AuthorRelationshipPeer, domain.AuthorRelationshipColleague:
		return domain.TestimonialRelationshipPeer
	case domain.AuthorRelationshipDirectReport:
		return domain.TestimonialRelationshipDirectReport
	case domain.AuthorRelationshipClient:
		return domain.TestimonialRelationshipClient
	default:
		return domain.TestimonialRelationshipOther
	}
}

// DeduplicateSkills removes duplicate skills by normalized name, preserving the first occurrence.
// It also trims whitespace and filters out empty strings.
func DeduplicateSkills(skills []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(skills))

	for _, skill := range skills {
		trimmed := strings.TrimSpace(skill)
		if trimmed == "" {
			continue
		}
		normalized := strings.ToLower(trimmed)
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, trimmed)
		}
	}
	return result
}
