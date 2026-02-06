package resolver

import (
	"encoding/json"
	"regexp"
	"strings"

	"backend/internal/domain"
	"backend/internal/graphql/model"
)

// jsonNull is the string representation of a JSON null value.
const jsonNull = "null"

// normalizeSkillCategory converts a domain SkillCategory to the GraphQL enum format (uppercase).
// Empty or unrecognized categories default to SOFT.
func normalizeSkillCategory(cat domain.SkillCategory) domain.SkillCategory {
	upper := domain.SkillCategory(strings.ToUpper(string(cat)))
	switch upper {
	case "TECHNICAL", "SOFT", "DOMAIN":
		return upper
	default:
		return "SOFT"
	}
}

// emailRegex is a basic regex for email validation.
// It's intentionally permissive to allow most valid emails.
var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// phoneRegex validates phone numbers with various formats.
// Allows digits, spaces, dashes, parentheses, dots, and leading +.
var phoneRegex = regexp.MustCompile(`^[\d\s\-().+]+$`)

// isValidEmail checks if the given string is a valid email format.
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// isValidPhone checks if the given string is a valid phone number format.
// This is a permissive check that allows common phone number formats.
func isValidPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if len(phone) < 5 || len(phone) > 30 {
		return false
	}
	return phoneRegex.MatchString(phone)
}

// normalizeSkillName produces a lowercase, trimmed version of a skill name for deduplication.
func normalizeSkillName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// toGraphQLUser converts a domain User to a GraphQL User model.
func toGraphQLUser(u *domain.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toGraphQLAuthor converts a domain Author to a GraphQL Author model.
func toGraphQLAuthor(a *domain.Author) *model.Author {
	if a == nil {
		return nil
	}
	return &model.Author{
		ID:          a.ID.String(),
		Name:        a.Name,
		Title:       a.Title,
		Company:     a.Company,
		LinkedInURL: a.LinkedInURL,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// toGraphQLAuthorWithImage converts a domain Author to a GraphQL Author model
// with an optional image URL for the author's profile picture.
func toGraphQLAuthorWithImage(a *domain.Author, imageURL *string) *model.Author {
	if a == nil {
		return nil
	}
	return &model.Author{
		ID:          a.ID.String(),
		Name:        a.Name,
		Title:       a.Title,
		Company:     a.Company,
		LinkedInURL: a.LinkedInURL,
		ImageURL:    imageURL,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// toGraphQLFile converts a domain File to a GraphQL File model.
// If user is provided, it will be set on the result.
func toGraphQLFile(f *domain.File, user *model.User) *model.File {
	if f == nil {
		return nil
	}
	return &model.File{
		ID:          f.ID.String(),
		Filename:    f.Filename,
		ContentType: f.ContentType,
		SizeBytes:   int(f.SizeBytes),
		StorageKey:  f.StorageKey,
		ContentHash: f.ContentHash,
		CreatedAt:   f.CreatedAt,
		User:        user,
	}
}

// toGraphQLReferenceLetter converts a domain ReferenceLetter to a GraphQL ReferenceLetter model.
// The user and file relations must be provided separately.
func toGraphQLReferenceLetter(rl *domain.ReferenceLetter, user *model.User, file *model.File) *model.ReferenceLetter {
	if rl == nil {
		return nil
	}

	// Convert domain status to GraphQL status
	var status model.ReferenceLetterStatus
	switch rl.Status {
	case domain.ReferenceLetterStatusPending:
		status = model.ReferenceLetterStatusPending
	case domain.ReferenceLetterStatusProcessing:
		status = model.ReferenceLetterStatusProcessing
	case domain.ReferenceLetterStatusCompleted:
		status = model.ReferenceLetterStatusCompleted
	case domain.ReferenceLetterStatusFailed:
		status = model.ReferenceLetterStatusFailed
	case domain.ReferenceLetterStatusApplied:
		status = model.ReferenceLetterStatusApplied
	default:
		status = model.ReferenceLetterStatusPending
	}

	return &model.ReferenceLetter{
		ID:            rl.ID.String(),
		Title:         rl.Title,
		AuthorName:    rl.AuthorName,
		AuthorTitle:   rl.AuthorTitle,
		Organization:  rl.Organization,
		DateWritten:   rl.DateWritten,
		RawText:       rl.RawText,
		ExtractedData: toGraphQLExtractedData(rl.ExtractedData),
		Status:        status,
		CreatedAt:     rl.CreatedAt,
		UpdatedAt:     rl.UpdatedAt,
		User:          user,
		File:          file,
	}
}

// toGraphQLExtractedData converts JSON raw message to typed ExtractedLetterData.
func toGraphQLExtractedData(raw json.RawMessage) *model.ExtractedLetterData {
	if len(raw) == 0 || string(raw) == jsonNull {
		return nil
	}

	var data domain.ExtractedLetterData
	if err := json.Unmarshal(raw, &data); err != nil {
		// Return nil for graceful degradation on parse errors
		return nil
	}

	return &model.ExtractedLetterData{
		Author:             toGraphQLExtractedAuthor(&data.Author),
		Testimonials:       toGraphQLExtractedTestimonials(data.Testimonials),
		SkillMentions:      toGraphQLExtractedSkillMentions(data.SkillMentions),
		ExperienceMentions: toGraphQLExtractedExperienceMentions(data.ExperienceMentions),
		DiscoveredSkills:   toGraphQLDiscoveredSkills(data.DiscoveredSkills),
		Metadata:           toGraphQLExtractionMetadata(&data.Metadata),
	}
}

// toGraphQLDiscoveredSkills converts a slice of domain DiscoveredSkill to GraphQL models.
func toGraphQLDiscoveredSkills(skills []domain.DiscoveredSkill) []*model.DiscoveredSkill {
	if len(skills) == 0 {
		return []*model.DiscoveredSkill{}
	}
	result := make([]*model.DiscoveredSkill, len(skills))
	for i, s := range skills {
		result[i] = &model.DiscoveredSkill{
			Skill:    s.Skill,
			Quote:    s.Quote,
			Context:  s.Context,
			Category: normalizeSkillCategory(s.Category),
		}
	}
	return result
}

// toGraphQLExtractedAuthor converts domain ExtractedAuthor to GraphQL model.
func toGraphQLExtractedAuthor(a *domain.ExtractedAuthor) *model.ExtractedAuthor {
	if a == nil {
		return nil
	}
	return &model.ExtractedAuthor{
		Name:         a.Name,
		Title:        a.Title,
		Company:      a.Company,
		Relationship: a.Relationship,
	}
}

// toGraphQLExtractedTestimonials converts a slice of domain ExtractedTestimonial to GraphQL models.
func toGraphQLExtractedTestimonials(testimonials []domain.ExtractedTestimonial) []*model.ExtractedTestimonial {
	if len(testimonials) == 0 {
		return []*model.ExtractedTestimonial{}
	}
	result := make([]*model.ExtractedTestimonial, len(testimonials))
	for i, t := range testimonials {
		result[i] = &model.ExtractedTestimonial{
			Quote:           t.Quote,
			SkillsMentioned: t.SkillsMentioned,
		}
	}
	return result
}

// toGraphQLExtractedSkillMentions converts a slice of domain ExtractedSkillMention to GraphQL models.
func toGraphQLExtractedSkillMentions(mentions []domain.ExtractedSkillMention) []*model.ExtractedSkillMention {
	if len(mentions) == 0 {
		return []*model.ExtractedSkillMention{}
	}
	result := make([]*model.ExtractedSkillMention, len(mentions))
	for i, m := range mentions {
		result[i] = &model.ExtractedSkillMention{
			Skill:   m.Skill,
			Quote:   m.Quote,
			Context: m.Context,
		}
	}
	return result
}

// toGraphQLExtractedExperienceMentions converts a slice of domain ExtractedExperienceMention to GraphQL models.
func toGraphQLExtractedExperienceMentions(mentions []domain.ExtractedExperienceMention) []*model.ExtractedExperienceMention {
	if len(mentions) == 0 {
		return []*model.ExtractedExperienceMention{}
	}
	result := make([]*model.ExtractedExperienceMention, len(mentions))
	for i, m := range mentions {
		result[i] = &model.ExtractedExperienceMention{
			Company: m.Company,
			Role:    m.Role,
			Quote:   m.Quote,
		}
	}
	return result
}

// toGraphQLExtractionMetadata converts domain ExtractionMetadata to GraphQL model.
func toGraphQLExtractionMetadata(m *domain.ExtractionMetadata) *model.ExtractionMetadata {
	if m == nil {
		return nil
	}
	return &model.ExtractionMetadata{
		ExtractedAt:      m.ExtractedAt,
		ModelVersion:     m.ModelVersion,
		ProcessingTimeMs: m.ProcessingTimeMs,
	}
}

// toGraphQLResume converts a domain Resume to a GraphQL Resume model.
// The user and file relations must be provided separately.
func toGraphQLResume(r *domain.Resume, user *model.User, file *model.File) *model.Resume {
	if r == nil {
		return nil
	}

	// Convert domain status to GraphQL status
	var status model.ResumeStatus
	switch r.Status {
	case domain.ResumeStatusPending:
		status = model.ResumeStatusPending
	case domain.ResumeStatusProcessing:
		status = model.ResumeStatusProcessing
	case domain.ResumeStatusCompleted:
		status = model.ResumeStatusCompleted
	case domain.ResumeStatusFailed:
		status = model.ResumeStatusFailed
	default:
		status = model.ResumeStatusPending
	}

	return &model.Resume{
		ID:            r.ID.String(),
		Status:        status,
		ExtractedData: toGraphQLResumeExtractedData(r.ExtractedData),
		ErrorMessage:  r.ErrorMessage,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		User:          user,
		File:          file,
	}
}

// toGraphQLResumeExtractedData converts JSON raw message to typed ResumeExtractedData.
func toGraphQLResumeExtractedData(raw json.RawMessage) *model.ResumeExtractedData {
	if len(raw) == 0 || string(raw) == jsonNull {
		return nil
	}

	var data domain.ResumeExtractedData
	if err := json.Unmarshal(raw, &data); err != nil {
		// Return nil for graceful degradation on parse errors
		return nil
	}

	// Convert experiences
	experiences := make([]*model.ExtractedWorkExperience, len(data.Experience))
	for i, exp := range data.Experience {
		experiences[i] = &model.ExtractedWorkExperience{
			Company:     exp.Company,
			Title:       exp.Title,
			Location:    exp.Location,
			StartDate:   exp.StartDate,
			EndDate:     exp.EndDate,
			IsCurrent:   exp.IsCurrent,
			Description: exp.Description,
		}
	}

	// Convert education
	educations := make([]*model.ExtractedEducation, len(data.Education))
	for i, edu := range data.Education {
		educations[i] = &model.ExtractedEducation{
			Institution:  edu.Institution,
			Degree:       edu.Degree,
			Field:        edu.Field,
			StartDate:    edu.StartDate,
			EndDate:      edu.EndDate,
			Gpa:          edu.GPA,
			Achievements: edu.Achievements,
		}
	}

	return &model.ResumeExtractedData{
		Name:        data.Name,
		Email:       data.Email,
		Phone:       data.Phone,
		Location:    data.Location,
		Summary:     data.Summary,
		Experiences: experiences,
		Educations:  educations,
		Skills:      data.Skills,
		ExtractedAt: data.ExtractedAt,
		Confidence:  data.Confidence,
	}
}

// toGraphQLProfile converts a domain Profile to a GraphQL Profile model.
// The user, experiences, educations, skills, and photoURL must be provided separately.
func toGraphQLProfile(p *domain.Profile, user *model.User, experiences []*model.ProfileExperience, educations []*model.ProfileEducation, skills []*model.ProfileSkill, photoURL *string) *model.Profile {
	if p == nil {
		return nil
	}
	return &model.Profile{
		ID:              p.ID.String(),
		User:            user,
		Name:            p.Name,
		Email:           p.Email,
		Phone:           p.Phone,
		Location:        p.Location,
		Summary:         p.Summary,
		ProfilePhotoURL: photoURL,
		Experiences:     experiences,
		Educations:      educations,
		Skills:          skills,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// domainProfileToGQL converts a domain Profile to a GraphQL Profile model
// without requiring relations (experiences, educations, skills).
// Use this when you only need the basic profile data without nested objects.
// The optional photoURL can be provided to populate the profile photo URL.
func domainProfileToGQL(p *domain.Profile, photoURL *string) *model.Profile {
	if p == nil {
		return nil
	}
	return &model.Profile{
		ID:              p.ID.String(),
		Name:            p.Name,
		Email:           p.Email,
		Phone:           p.Phone,
		Location:        p.Location,
		Summary:         p.Summary,
		ProfilePhotoURL: photoURL,
		Experiences:     []*model.ProfileExperience{},
		Educations:      []*model.ProfileEducation{},
		Skills:          []*model.ProfileSkill{},
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// toGraphQLProfileExperience converts a domain ProfileExperience to a GraphQL ProfileExperience model.
func toGraphQLProfileExperience(e *domain.ProfileExperience) *model.ProfileExperience {
	if e == nil {
		return nil
	}

	// Convert domain source to GraphQL source
	var source model.ExperienceSource
	switch e.Source {
	case domain.ExperienceSourceManual:
		source = model.ExperienceSourceManual
	case domain.ExperienceSourceResumeExtracted:
		source = model.ExperienceSourceResumeExtracted
	case domain.ExperienceSourceLetterDiscovered:
		source = model.ExperienceSourceLetterDiscovered
	default:
		source = model.ExperienceSourceManual
	}

	// Convert highlights from pq.StringArray to []string
	highlights := make([]string, len(e.Highlights))
	copy(highlights, e.Highlights)

	return &model.ProfileExperience{
		ID:           e.ID.String(),
		Company:      e.Company,
		Title:        e.Title,
		Location:     e.Location,
		StartDate:    e.StartDate,
		EndDate:      e.EndDate,
		IsCurrent:    e.IsCurrent,
		Description:  e.Description,
		Highlights:   highlights,
		DisplayOrder: e.DisplayOrder,
		Source:       source,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// stringPtr returns a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// toGraphQLProfileExperiences converts a slice of domain ProfileExperience to GraphQL models.
func toGraphQLProfileExperiences(experiences []*domain.ProfileExperience) []*model.ProfileExperience {
	if len(experiences) == 0 {
		return []*model.ProfileExperience{}
	}
	result := make([]*model.ProfileExperience, len(experiences))
	for i, e := range experiences {
		result[i] = toGraphQLProfileExperience(e)
	}
	return result
}

// toGraphQLProfileEducation converts a domain ProfileEducation to a GraphQL ProfileEducation model.
func toGraphQLProfileEducation(e *domain.ProfileEducation) *model.ProfileEducation {
	if e == nil {
		return nil
	}

	// Convert domain source to GraphQL source
	var source model.ExperienceSource
	switch e.Source {
	case domain.ExperienceSourceManual:
		source = model.ExperienceSourceManual
	case domain.ExperienceSourceResumeExtracted:
		source = model.ExperienceSourceResumeExtracted
	case domain.ExperienceSourceLetterDiscovered:
		source = model.ExperienceSourceLetterDiscovered
	default:
		source = model.ExperienceSourceManual
	}

	return &model.ProfileEducation{
		ID:           e.ID.String(),
		Institution:  e.Institution,
		Degree:       e.Degree,
		Field:        e.Field,
		StartDate:    e.StartDate,
		EndDate:      e.EndDate,
		IsCurrent:    e.IsCurrent,
		Description:  e.Description,
		Gpa:          e.GPA,
		DisplayOrder: e.DisplayOrder,
		Source:        source,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// toGraphQLProfileEducations converts a slice of domain ProfileEducation to GraphQL models.
func toGraphQLProfileEducations(educations []*domain.ProfileEducation) []*model.ProfileEducation {
	if len(educations) == 0 {
		return []*model.ProfileEducation{}
	}
	result := make([]*model.ProfileEducation, len(educations))
	for i, e := range educations {
		result[i] = toGraphQLProfileEducation(e)
	}
	return result
}

// toGraphQLProfileSkill converts a domain ProfileSkill to a GraphQL ProfileSkill model.
func toGraphQLProfileSkill(s *domain.ProfileSkill) *model.ProfileSkill {
	if s == nil {
		return nil
	}

	// Convert domain source to GraphQL source
	var source model.ExperienceSource
	switch s.Source {
	case domain.ExperienceSourceManual:
		source = model.ExperienceSourceManual
	case domain.ExperienceSourceResumeExtracted:
		source = model.ExperienceSourceResumeExtracted
	case domain.ExperienceSourceLetterDiscovered:
		source = model.ExperienceSourceLetterDiscovered
	default:
		source = model.ExperienceSourceManual
	}

	return &model.ProfileSkill{
		ID:             s.ID.String(),
		Name:           s.Name,
		NormalizedName: s.NormalizedName,
		Category:       domain.SkillCategory(strings.ToUpper(s.Category)),
		DisplayOrder:   s.DisplayOrder,
		Source:         source,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

// toGraphQLProfileSkills converts a slice of domain ProfileSkill to GraphQL models.
func toGraphQLProfileSkills(skills []*domain.ProfileSkill) []*model.ProfileSkill {
	if len(skills) == 0 {
		return []*model.ProfileSkill{}
	}
	result := make([]*model.ProfileSkill, len(skills))
	for i, s := range skills {
		result[i] = toGraphQLProfileSkill(s)
	}
	return result
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

// toGraphQLTestimonial converts a domain Testimonial to a GraphQL Testimonial model.
// The referenceLetter relation and validatedSkills can be provided separately if needed.
func toGraphQLTestimonial(t *domain.Testimonial, referenceLetter *model.ReferenceLetter, validatedSkills []*model.ProfileSkill) *model.Testimonial {
	if t == nil {
		return nil
	}

	// Convert domain relationship to GraphQL relationship
	var relationship model.TestimonialRelationship
	switch t.Relationship {
	case domain.TestimonialRelationshipManager:
		relationship = model.TestimonialRelationshipManager
	case domain.TestimonialRelationshipPeer:
		relationship = model.TestimonialRelationshipPeer
	case domain.TestimonialRelationshipDirectReport:
		relationship = model.TestimonialRelationshipDirectReport
	case domain.TestimonialRelationshipClient:
		relationship = model.TestimonialRelationshipClient
	default:
		relationship = model.TestimonialRelationshipOther
	}

	// Ensure validatedSkills is not nil (GraphQL requires non-null list)
	if validatedSkills == nil {
		validatedSkills = []*model.ProfileSkill{}
	}

	// Get author information - prefer Author relation, fall back to legacy fields
	var authorName string
	var authorTitle *string
	var authorCompany *string
	var author *model.Author

	if t.Author != nil {
		author = toGraphQLAuthor(t.Author)
		authorName = t.Author.Name
		authorTitle = t.Author.Title
		authorCompany = t.Author.Company
	} else {
		// Legacy: use denormalized fields
		if t.AuthorName != nil {
			authorName = *t.AuthorName
		}
		authorTitle = t.AuthorTitle
		authorCompany = t.AuthorCompany
	}

	return &model.Testimonial{
		ID:              t.ID.String(),
		Quote:           t.Quote,
		Author:          author,
		AuthorName:      authorName,
		AuthorTitle:     authorTitle,
		AuthorCompany:   authorCompany,
		Relationship:    relationship,
		ReferenceLetter: referenceLetter,
		CreatedAt:       t.CreatedAt,
		ValidatedSkills: validatedSkills,
	}
}

// toGraphQLTestimonials converts a slice of domain Testimonial to GraphQL models.
// validatedSkillsByRefLetter maps reference letter IDs to their validated skills.
// Only skills that are mentioned in the testimonial's SkillsMentioned field are included.
func toGraphQLTestimonials(testimonials []*domain.Testimonial, validatedSkillsByRefLetter map[string][]*model.ProfileSkill) []*model.Testimonial {
	if len(testimonials) == 0 {
		return []*model.Testimonial{}
	}
	result := make([]*model.Testimonial, len(testimonials))
	for i, t := range testimonials {
		var validatedSkills []*model.ProfileSkill
		if validatedSkillsByRefLetter != nil {
			allSkills := validatedSkillsByRefLetter[t.ReferenceLetterID.String()]
			// Filter to only include skills that are mentioned in this testimonial
			validatedSkills = filterSkillsByMentioned(allSkills, t.SkillsMentioned)
		}
		result[i] = toGraphQLTestimonial(t, nil, validatedSkills)
	}
	return result
}

// filterSkillsByMentioned filters a list of skills to only include those whose name
// (case-insensitive) appears in the skillsMentioned list.
func filterSkillsByMentioned(skills []*model.ProfileSkill, skillsMentioned []string) []*model.ProfileSkill {
	if len(skillsMentioned) == 0 {
		return []*model.ProfileSkill{}
	}

	// Build a set of mentioned skill names (lowercase for case-insensitive matching)
	mentionedSet := make(map[string]struct{}, len(skillsMentioned))
	for _, name := range skillsMentioned {
		mentionedSet[strings.ToLower(name)] = struct{}{}
	}

	// Filter skills to only include those mentioned
	var filtered []*model.ProfileSkill
	for _, skill := range skills {
		if _, ok := mentionedSet[strings.ToLower(skill.Name)]; ok {
			filtered = append(filtered, skill)
		}
	}

	if filtered == nil {
		return []*model.ProfileSkill{}
	}
	return filtered
}
