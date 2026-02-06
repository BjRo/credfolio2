package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"backend/internal/domain"
)

// MaterializationResult contains counts of materialized items.
type MaterializationResult struct {
	Experiences  int
	Educations   int
	Skills       int
	Testimonials int
}

// CrossReferenceResult contains counts of auto-applied validations.
type CrossReferenceResult struct {
	SkillValidations      int
	ExperienceValidations int
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
	skillValRepo     domain.SkillValidationRepository
	expValRepo       domain.ExperienceValidationRepository
}

// NewMaterializationService creates a new MaterializationService.
func NewMaterializationService(
	profileRepo domain.ProfileRepository,
	profileExpRepo domain.ProfileExperienceRepository,
	profileEduRepo domain.ProfileEducationRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	authorRepo domain.AuthorRepository,
	testimonialRepo domain.TestimonialRepository,
	skillValRepo domain.SkillValidationRepository,
	expValRepo domain.ExperienceValidationRepository,
) *MaterializationService {
	return &MaterializationService{
		profileRepo:      profileRepo,
		profileExpRepo:   profileExpRepo,
		profileEduRepo:   profileEduRepo,
		profileSkillRepo: profileSkillRepo,
		authorRepo:       authorRepo,
		testimonialRepo:  testimonialRepo,
		skillValRepo:     skillValRepo,
		expValRepo:       expValRepo,
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

	// Populate empty profile header fields from extracted data (preserves user edits)
	if err := s.populateProfileHeader(ctx, profile, data); err != nil {
		return nil, err
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

// MaterializeReferenceLetterData creates testimonial rows from extracted reference letter data.
// It finds or creates an Author entity, then creates Testimonial records for each extracted testimonial.
// Idempotent: deletes existing testimonials from the same reference letter before re-creating.
func (s *MaterializationService) MaterializeReferenceLetterData(
	ctx context.Context,
	referenceLetterID uuid.UUID,
	userID uuid.UUID,
	data *domain.ExtractedLetterData,
) (*MaterializationResult, error) {
	// Get or create the user's profile
	profile, err := s.profileRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create profile: %w", err)
	}

	// Delete any existing testimonials from this reference letter (idempotent re-processing)
	// TODO: Wrap delete + create cycle in a DB transaction for atomicity. A crash between
	// delete and the last create could leave the user with lost testimonials. Same pattern
	// exists in MaterializeResumeData — address both together.
	if delErr := s.testimonialRepo.DeleteByReferenceLetterID(ctx, referenceLetterID); delErr != nil {
		return nil, fmt.Errorf("failed to delete existing testimonials for reference letter: %w", delErr)
	}

	result := &MaterializationResult{}

	if len(data.Testimonials) == 0 {
		return result, nil
	}

	// Find or create the Author entity
	author, err := s.findOrCreateAuthor(ctx, profile.ID, &data.Author)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create author: %w", err)
	}

	// Create testimonial records
	for _, extracted := range data.Testimonials {
		testimonial := &domain.Testimonial{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			ReferenceLetterID: referenceLetterID,
			Quote:             extracted.Quote,
			Relationship:      mapAuthorRelationship(data.Author.Relationship),
			SkillsMentioned:   extracted.SkillsMentioned,
			AuthorName:        &data.Author.Name,
			AuthorTitle:       data.Author.Title,
			AuthorCompany:     data.Author.Company,
		}

		if author != nil {
			testimonial.AuthorID = &author.ID
		}

		if createErr := s.testimonialRepo.Create(ctx, testimonial); createErr != nil {
			return result, fmt.Errorf("failed to create testimonial: %w", createErr)
		}
		result.Testimonials++
	}

	return result, nil
}

// findOrCreateAuthor finds an existing author by name and company, or creates a new one.
// TODO: Add a unique DB constraint on (profile_id, name, company) to guard against
// duplicate authors from concurrent imports (TOCTOU race in the find-then-create pattern).
// TODO: Consider updating existing author's Title when reusing, so newer reference letters
// with updated titles (e.g., promotions) refresh the canonical Author entity.
func (s *MaterializationService) findOrCreateAuthor(ctx context.Context, profileID uuid.UUID, extracted *domain.ExtractedAuthor) (*domain.Author, error) {
	// Try to find existing author with same name and company
	existing, err := s.authorRepo.FindByNameAndCompany(ctx, profileID, extracted.Name, extracted.Company)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing author: %w", err)
	}
	if existing != nil {
		return existing, nil
	}

	// Create new author
	author := &domain.Author{
		ID:        uuid.New(),
		ProfileID: profileID,
		Name:      extracted.Name,
		Title:     extracted.Title,
		Company:   extracted.Company,
	}
	if createErr := s.authorRepo.Create(ctx, author); createErr != nil {
		return nil, fmt.Errorf("failed to create author: %w", createErr)
	}
	return author, nil
}

// populateProfileHeader fills empty profile header fields from extracted resume data.
// Fields that already have values (e.g., from user edits) are preserved.
func (s *MaterializationService) populateProfileHeader(ctx context.Context, profile *domain.Profile, data *domain.ResumeExtractedData) error {
	updated := false
	if profile.Name == nil && data.Name != "" {
		profile.Name = &data.Name
		updated = true
	}
	if profile.Email == nil && data.Email != nil {
		profile.Email = data.Email
		updated = true
	}
	if profile.Phone == nil && data.Phone != nil {
		profile.Phone = data.Phone
		updated = true
	}
	if profile.Location == nil && data.Location != nil {
		profile.Location = data.Location
		updated = true
	}
	if profile.Summary == nil && data.Summary != nil {
		profile.Summary = data.Summary
		updated = true
	}
	if updated {
		if err := s.profileRepo.Update(ctx, profile); err != nil {
			return fmt.Errorf("failed to update profile header: %w", err)
		}
	}
	return nil
}

// CrossReferenceValidations matches reference letter skill/experience mentions
// against the profile's skills and experiences, creating validation records for matches.
// This is called during unified import when both a resume and reference letter are imported together.
func (s *MaterializationService) CrossReferenceValidations(
	ctx context.Context,
	profileID uuid.UUID,
	referenceLetterID uuid.UUID,
	letterData *domain.ExtractedLetterData,
) (*CrossReferenceResult, error) {
	result := &CrossReferenceResult{}

	// Get all profile skills for matching
	skills, err := s.profileSkillRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile skills: %w", err)
	}

	// Build normalized name -> skill lookup
	skillByNorm := make(map[string]*domain.ProfileSkill, len(skills))
	for _, sk := range skills {
		skillByNorm[sk.NormalizedName] = sk
	}

	// Match skill mentions to profile skills
	for _, mention := range letterData.SkillMentions {
		normalized := strings.ToLower(strings.TrimSpace(mention.Skill))
		if sk, ok := skillByNorm[normalized]; ok {
			quote := mention.Quote
			validation := &domain.SkillValidation{
				ID:                uuid.New(),
				ProfileSkillID:    sk.ID,
				ReferenceLetterID: referenceLetterID,
				QuoteSnippet:      &quote,
			}
			if createErr := s.skillValRepo.Create(ctx, validation); createErr != nil {
				// Ignore duplicate constraint errors
				if !strings.Contains(createErr.Error(), "duplicate") && !strings.Contains(createErr.Error(), "unique constraint") {
					return result, fmt.Errorf("failed to create skill validation: %w", createErr)
				}
				continue
			}
			result.SkillValidations++
		}
	}

	// Get all profile experiences for matching
	experiences, err := s.profileExpRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile experiences: %w", err)
	}

	// Build normalized company -> experience lookup
	expByCompany := make(map[string]*domain.ProfileExperience, len(experiences))
	for _, exp := range experiences {
		norm := strings.ToLower(strings.TrimSpace(exp.Company))
		expByCompany[norm] = exp
	}

	// Match experience mentions to profile experiences
	for _, mention := range letterData.ExperienceMentions {
		normalized := strings.ToLower(strings.TrimSpace(mention.Company))
		if exp, ok := expByCompany[normalized]; ok {
			quote := mention.Quote
			validation := &domain.ExperienceValidation{
				ID:                  uuid.New(),
				ProfileExperienceID: exp.ID,
				ReferenceLetterID:   referenceLetterID,
				QuoteSnippet:        &quote,
			}
			if createErr := s.expValRepo.Create(ctx, validation); createErr != nil {
				if !strings.Contains(createErr.Error(), "duplicate") && !strings.Contains(createErr.Error(), "unique constraint") {
					return result, fmt.Errorf("failed to create experience validation: %w", createErr)
				}
				continue
			}
			result.ExperienceValidations++
		}
	}

	return result, nil
}

// mapAuthorRelationship maps an AuthorRelationship to a TestimonialRelationship.
func mapAuthorRelationship(ar domain.AuthorRelationship) domain.TestimonialRelationship {
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
