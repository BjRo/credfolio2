package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/logger"
)

// ReferenceLetterProcessingArgs contains the arguments for a reference letter processing job.
type ReferenceLetterProcessingArgs struct { //nolint:govet // Field ordering prioritizes JSON consistency over memory alignment
	StorageKey        string    `json:"storage_key"`
	ReferenceLetterID uuid.UUID `json:"reference_letter_id"`
	FileID            uuid.UUID `json:"file_id"`
	ContentType       string    `json:"content_type"`
}

// Kind returns the job type identifier for River.
func (ReferenceLetterProcessingArgs) Kind() string {
	return "reference_letter_processing"
}

// ReferenceLetterProcessingWorker processes uploaded reference letters to extract credibility data.
type ReferenceLetterProcessingWorker struct {
	river.WorkerDefaults[ReferenceLetterProcessingArgs]
	letterRepo          domain.ReferenceLetterRepository
	fileRepo            domain.FileRepository
	profileRepo         domain.ProfileRepository
	profileSkillRepo    domain.ProfileSkillRepository
	authorRepo          domain.AuthorRepository
	testimonialRepo     domain.TestimonialRepository
	skillValidationRepo domain.SkillValidationRepository
	storage             domain.Storage
	extractor           domain.DocumentExtractor
	log                 logger.Logger
}

// NewReferenceLetterProcessingWorker creates a new reference letter processing worker.
func NewReferenceLetterProcessingWorker(
	letterRepo domain.ReferenceLetterRepository,
	fileRepo domain.FileRepository,
	profileRepo domain.ProfileRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	authorRepo domain.AuthorRepository,
	testimonialRepo domain.TestimonialRepository,
	skillValidationRepo domain.SkillValidationRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *ReferenceLetterProcessingWorker {
	return &ReferenceLetterProcessingWorker{
		letterRepo:          letterRepo,
		fileRepo:            fileRepo,
		profileRepo:         profileRepo,
		profileSkillRepo:    profileSkillRepo,
		authorRepo:          authorRepo,
		testimonialRepo:     testimonialRepo,
		skillValidationRepo: skillValidationRepo,
		storage:             storage,
		extractor:           extractor,
		log:                 log,
	}
}

// Work processes a reference letter and extracts credibility data using LLM.
func (w *ReferenceLetterProcessingWorker) Work(ctx context.Context, job *river.Job[ReferenceLetterProcessingArgs]) error {
	args := job.Args
	w.log.Info("Processing reference letter",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("file_id", args.FileID.String()),
		logger.String("storage_key", args.StorageKey),
	)

	// Update status to processing
	if err := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusProcessing, nil); err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Get content type from file record (in case it wasn't passed in args)
	contentType := args.ContentType
	if contentType == "" {
		file, err := w.fileRepo.GetByID(ctx, args.FileID)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get file record: %v", err)
			w.log.Error("Failed to get file record",
				logger.Feature("jobs"),
				logger.String("reference_letter_id", args.ReferenceLetterID.String()),
				logger.String("file_id", args.FileID.String()),
				logger.Err(err),
			)
			_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
			return fmt.Errorf("failed to get file record: %w", err)
		}
		if file == nil {
			errMsg := "file record not found"
			_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
			return fmt.Errorf("file record not found: %s", args.FileID)
		}
		contentType = file.ContentType
	}

	// Download file from storage
	reader, err := w.storage.Download(ctx, args.StorageKey)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download file: %v", err)
		w.log.Error("Storage download failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.String("storage_key", args.StorageKey),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close() //nolint:errcheck // Best effort cleanup

	// Read file data
	data, err := io.ReadAll(reader)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file data: %v", err)
		w.log.Error("Failed to read file data",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Get reference letter to find user ID
	letter, err := w.letterRepo.GetByID(ctx, args.ReferenceLetterID)
	if err != nil || letter == nil {
		errMsg := "reference letter not found"
		if err != nil {
			errMsg = fmt.Sprintf("failed to get reference letter: %v", err)
		}
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("reference letter not found: %s", args.ReferenceLetterID)
	}

	// Get the user's profile and existing skills for context
	profileSkills, existingSkillsMap := w.getProfileSkillsContext(ctx, letter.UserID)

	// Extract credibility data using LLM with profile skills context
	extractedData, err := w.extractLetterData(ctx, data, contentType, profileSkills)
	if err != nil {
		errMsg := fmt.Sprintf("failed to extract letter data: %v", err)
		w.log.Error("Letter extraction failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to extract letter data: %w", err)
	}

	// Save extracted data and mark as completed
	if saveErr := w.saveExtractedData(ctx, args.ReferenceLetterID, extractedData); saveErr != nil {
		errMsg := fmt.Sprintf("failed to save extracted data: %v", saveErr)
		w.log.Error("Failed to save extracted data",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(saveErr),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to save extracted data: %w", saveErr)
	}

	// Process extracted data: create testimonials, skill validations, and discovered skills
	if processErr := w.processExtractedData(ctx, letter.UserID, args.ReferenceLetterID, extractedData, existingSkillsMap); processErr != nil {
		// Log the error but don't fail the job - the extraction succeeded, just post-processing failed
		w.log.Warning("Failed to process extracted data into profile records",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(processErr),
		)
	}

	w.log.Info("Reference letter processing completed",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("author_name", extractedData.Author.Name),
		logger.Int("testimonials_count", len(extractedData.Testimonials)),
		logger.Int("skill_mentions_count", len(extractedData.SkillMentions)),
		logger.Int("experience_mentions_count", len(extractedData.ExperienceMentions)),
		logger.Int("discovered_skills_count", len(extractedData.DiscoveredSkills)),
	)

	return nil
}

// extractLetterData uses the LLM to extract structured credibility data from the reference letter.
// The profileSkills parameter provides context about existing profile skills, enabling the LLM to
// distinguish between mentions of existing skills (for validation) and newly discovered skills.
func (w *ReferenceLetterProcessingWorker) extractLetterData(ctx context.Context, data []byte, contentType string, profileSkills []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	// First, extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Then, use LLM to extract structured credibility data from the text
	extractedData, err := w.extractor.ExtractLetterData(ctx, text, profileSkills)
	if err != nil {
		return nil, fmt.Errorf("failed to extract letter data: %w", err)
	}

	// Set extraction timestamp
	extractedData.Metadata = domain.ExtractionMetadata{
		ExtractedAt:  time.Now(),
		ModelVersion: "claude-sonnet-4-20250514", // TODO: Get from extractor config
	}

	return extractedData, nil
}

// saveExtractedData saves the extracted data to the reference letter record and marks as completed.
func (w *ReferenceLetterProcessingWorker) saveExtractedData(ctx context.Context, letterID uuid.UUID, data *domain.ExtractedLetterData) error {
	letter, err := w.letterRepo.GetByID(ctx, letterID)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", letterID)
	}

	// Marshal extracted data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	// Update letter with extracted data
	letter.ExtractedData = jsonData
	letter.Status = domain.ReferenceLetterStatusCompleted
	letter.ErrorMessage = nil

	// Also populate the author fields for easier querying
	if data.Author.Name != "" {
		letter.AuthorName = &data.Author.Name
	}
	if data.Author.Title != nil {
		letter.AuthorTitle = data.Author.Title
	}
	if data.Author.Company != nil {
		letter.Organization = data.Author.Company
	}

	if err := w.letterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}

	return nil
}

// updateStatus updates the reference letter status.
func (w *ReferenceLetterProcessingWorker) updateStatus(ctx context.Context, id uuid.UUID, status domain.ReferenceLetterStatus, errMsg *string) error {
	letter, err := w.letterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", id)
	}

	letter.Status = status
	letter.ErrorMessage = errMsg

	if err := w.letterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}

	return nil
}

// updateStatusFailed is a helper to update status to failed with an error message.
func (w *ReferenceLetterProcessingWorker) updateStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return w.updateStatus(ctx, id, domain.ReferenceLetterStatusFailed, &errMsg)
}

// normalizeSkillName produces a lowercase, trimmed version of a skill name for matching.
func normalizeSkillName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// getProfileSkillsContext fetches the user's profile skills and returns:
// 1. A slice of ProfileSkillContext for the LLM
// 2. A map of normalized skill names to ProfileSkill records for matching
func (w *ReferenceLetterProcessingWorker) getProfileSkillsContext(ctx context.Context, userID uuid.UUID) ([]domain.ProfileSkillContext, map[string]*domain.ProfileSkill) {
	// Get user's profile
	profile, err := w.profileRepo.GetByUserID(ctx, userID)
	if err != nil || profile == nil {
		w.log.Warning("Could not get profile for skill context",
			logger.Feature("jobs"),
			logger.String("user_id", userID.String()),
			logger.Err(err),
		)
		return nil, nil
	}

	// Get profile skills
	skills, err := w.profileSkillRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		w.log.Warning("Could not get profile skills for context",
			logger.Feature("jobs"),
			logger.String("profile_id", profile.ID.String()),
			logger.Err(err),
		)
		return nil, nil
	}

	// Build context and lookup map
	skillContext := make([]domain.ProfileSkillContext, 0, len(skills))
	skillsMap := make(map[string]*domain.ProfileSkill, len(skills))
	for _, skill := range skills {
		skillContext = append(skillContext, domain.ProfileSkillContext{
			Name:           skill.Name,
			NormalizedName: skill.NormalizedName,
			Category:       skill.Category,
		})
		skillsMap[normalizeSkillName(skill.Name)] = skill
		// Also add by normalized name for better matching
		if skill.NormalizedName != "" {
			skillsMap[skill.NormalizedName] = skill
		}
	}

	return skillContext, skillsMap
}

// processExtractedData creates testimonials, skill validations, and discovered skills from extracted data.
//
//nolint:gocyclo // This orchestration function inherently handles multiple related operations
func (w *ReferenceLetterProcessingWorker) processExtractedData(
	ctx context.Context,
	userID uuid.UUID,
	letterID uuid.UUID,
	data *domain.ExtractedLetterData,
	existingSkills map[string]*domain.ProfileSkill,
) error {
	// Initialize map if nil to allow adding discovered skills
	if existingSkills == nil {
		existingSkills = make(map[string]*domain.ProfileSkill)
	}

	// Get or create profile for the user
	profile, err := w.profileRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get/create profile: %w", err)
	}

	// Find or create the author
	author, err := w.findOrCreateAuthor(ctx, profile.ID, &data.Author)
	if err != nil {
		return fmt.Errorf("failed to find/create author: %w", err)
	}

	// Create testimonials and collect skill mentions per testimonial
	testimonialSkillMentions := make(map[uuid.UUID][]string) // testimonial ID -> skill names mentioned
	for _, t := range data.Testimonials {
		testimonial := &domain.Testimonial{
			ProfileID:         profile.ID,
			ReferenceLetterID: letterID,
			AuthorID:          &author.ID,
			Quote:             t.Quote,
			AuthorName:        &author.Name,
			AuthorTitle:       author.Title,
			AuthorCompany:     author.Company,
			Relationship:      mapAuthorToTestimonialRelationship(data.Author.Relationship),
			SkillsMentioned:   t.SkillsMentioned,
		}
		if err := w.testimonialRepo.Create(ctx, testimonial); err != nil {
			w.log.Warning("Failed to create testimonial",
				logger.Feature("jobs"),
				logger.String("letter_id", letterID.String()),
				logger.Err(err),
			)
			continue
		}
		// Track which skills are mentioned in this testimonial
		testimonialSkillMentions[testimonial.ID] = t.SkillsMentioned
	}

	// Create skill validations for skills mentioned in testimonials
	for testimonialID, skillNames := range testimonialSkillMentions {
		for _, skillName := range skillNames {
			normalizedName := normalizeSkillName(skillName)
			if existingSkill, ok := existingSkills[normalizedName]; ok {
				// Create skill validation linking to this testimonial
				validation := &domain.SkillValidation{
					ProfileSkillID:    existingSkill.ID,
					ReferenceLetterID: letterID,
					TestimonialID:     &testimonialID,
					QuoteSnippet:      nil, // The quote is in the testimonial
				}
				if err := w.skillValidationRepo.Create(ctx, validation); err != nil {
					w.log.Warning("Failed to create skill validation",
						logger.Feature("jobs"),
						logger.String("skill_id", existingSkill.ID.String()),
						logger.String("testimonial_id", testimonialID.String()),
						logger.Err(err),
					)
				}
			}
		}
	}

	// Create skill validations for explicit skill mentions (with quotes)
	for _, sm := range data.SkillMentions {
		normalizedName := normalizeSkillName(sm.Skill)
		if existingSkill, ok := existingSkills[normalizedName]; ok {
			validation := &domain.SkillValidation{
				ProfileSkillID:    existingSkill.ID,
				ReferenceLetterID: letterID,
				TestimonialID:     nil, // Not from a specific testimonial
				QuoteSnippet:      &sm.Quote,
			}
			if err := w.skillValidationRepo.Create(ctx, validation); err != nil {
				w.log.Warning("Failed to create skill validation from mention",
					logger.Feature("jobs"),
					logger.String("skill", sm.Skill),
					logger.Err(err),
				)
			}
		}
	}

	// Create profile skills for discovered skills
	for _, ds := range data.DiscoveredSkills {
		normalizedName := normalizeSkillName(ds.Skill)
		// Skip if skill already exists
		if _, exists := existingSkills[normalizedName]; exists {
			continue
		}

		// Determine category from context
		category := "SOFT" // Default to soft skills for discovered skills
		if ds.Context != nil {
			contextLower := strings.ToLower(*ds.Context)
			if strings.Contains(contextLower, "technical") || strings.Contains(contextLower, "programming") {
				category = "TECHNICAL" //nolint:goconst // Well-known skill category value
			} else if strings.Contains(contextLower, "domain") || strings.Contains(contextLower, "industry") {
				category = "DOMAIN"
			}
		}

		// Create the new skill
		newSkill := &domain.ProfileSkill{
			ProfileID:               profile.ID,
			Name:                    ds.Skill,
			NormalizedName:          normalizedName,
			Category:                category,
			Source:                  domain.ExperienceSourceManual, // Will mark as "discovered" via SourceReferenceLetterID
			SourceReferenceLetterID: &letterID,
		}
		if err := w.profileSkillRepo.CreateIgnoreDuplicate(ctx, newSkill); err != nil {
			w.log.Warning("Failed to create discovered skill",
				logger.Feature("jobs"),
				logger.String("skill", ds.Skill),
				logger.Err(err),
			)
			continue
		}

		// Add to existing skills map for subsequent validations
		existingSkills[normalizedName] = newSkill
	}

	return nil
}

// findOrCreateAuthor finds an existing author or creates a new one.
func (w *ReferenceLetterProcessingWorker) findOrCreateAuthor(ctx context.Context, profileID uuid.UUID, extractedAuthor *domain.ExtractedAuthor) (*domain.Author, error) {
	// Try to find existing author
	author, err := w.authorRepo.FindByNameAndCompany(ctx, profileID, extractedAuthor.Name, extractedAuthor.Company)
	if err != nil {
		return nil, fmt.Errorf("failed to search for author: %w", err)
	}
	if author != nil {
		return author, nil
	}

	// Create new author
	author = &domain.Author{
		ProfileID: profileID,
		Name:      extractedAuthor.Name,
		Title:     extractedAuthor.Title,
		Company:   extractedAuthor.Company,
	}
	if err := w.authorRepo.Create(ctx, author); err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return author, nil
}

// mapAuthorToTestimonialRelationship maps an AuthorRelationship to a TestimonialRelationship.
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
