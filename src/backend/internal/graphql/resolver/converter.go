package resolver

import (
	"encoding/json"
	"strings"

	"backend/internal/domain"
	"backend/internal/graphql/model"
)

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
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var data domain.ExtractedLetterData
	if err := json.Unmarshal(raw, &data); err != nil {
		// Return nil for graceful degradation on parse errors
		return nil
	}

	return &model.ExtractedLetterData{
		Author:          toGraphQLExtractedAuthor(&data.Author),
		Skills:          toGraphQLExtractedSkills(data.Skills),
		Qualities:       toGraphQLExtractedQualities(data.Qualities),
		Accomplishments: toGraphQLExtractedAccomplishments(data.Accomplishments),
		Recommendation:  toGraphQLExtractedRecommendation(&data.Recommendation),
		Metadata:        toGraphQLExtractionMetadata(&data.Metadata),
	}
}

// toGraphQLExtractedAuthor converts domain ExtractedAuthor to GraphQL model.
func toGraphQLExtractedAuthor(a *domain.ExtractedAuthor) *model.ExtractedAuthor {
	if a == nil {
		return nil
	}
	return &model.ExtractedAuthor{
		Name:                a.Name,
		Title:               a.Title,
		Organization:        a.Organization,
		Relationship:        strings.ToUpper(string(a.Relationship)),
		RelationshipDetails: a.RelationshipDetails,
		Confidence:          a.Confidence,
	}
}

// toGraphQLExtractedSkills converts a slice of domain ExtractedSkill to GraphQL models.
func toGraphQLExtractedSkills(skills []domain.ExtractedSkill) []*model.ExtractedSkill {
	if len(skills) == 0 {
		return []*model.ExtractedSkill{}
	}
	result := make([]*model.ExtractedSkill, len(skills))
	for i, s := range skills {
		result[i] = &model.ExtractedSkill{
			Name:           s.Name,
			NormalizedName: s.NormalizedName,
			Category:       strings.ToUpper(string(s.Category)),
			Mentions:       s.Mentions,
			Context:        s.Context,
			Confidence:     s.Confidence,
		}
	}
	return result
}

// toGraphQLExtractedQualities converts a slice of domain ExtractedQuality to GraphQL models.
func toGraphQLExtractedQualities(qualities []domain.ExtractedQuality) []*model.ExtractedQuality {
	if len(qualities) == 0 {
		return []*model.ExtractedQuality{}
	}
	result := make([]*model.ExtractedQuality, len(qualities))
	for i, q := range qualities {
		result[i] = &model.ExtractedQuality{
			Trait:      q.Trait,
			Evidence:   q.Evidence,
			Confidence: q.Confidence,
		}
	}
	return result
}

// toGraphQLExtractedAccomplishments converts a slice of domain ExtractedAccomplishment to GraphQL models.
func toGraphQLExtractedAccomplishments(accomplishments []domain.ExtractedAccomplishment) []*model.ExtractedAccomplishment {
	if len(accomplishments) == 0 {
		return []*model.ExtractedAccomplishment{}
	}
	result := make([]*model.ExtractedAccomplishment, len(accomplishments))
	for i, a := range accomplishments {
		result[i] = &model.ExtractedAccomplishment{
			Description: a.Description,
			Impact:      a.Impact,
			Timeframe:   a.Timeframe,
			Confidence:  a.Confidence,
		}
	}
	return result
}

// toGraphQLExtractedRecommendation converts domain ExtractedRecommendation to GraphQL model.
func toGraphQLExtractedRecommendation(r *domain.ExtractedRecommendation) *model.ExtractedRecommendation {
	if r == nil {
		return nil
	}
	return &model.ExtractedRecommendation{
		Strength:   strings.ToUpper(string(r.Strength)),
		Sentiment:  r.Sentiment,
		KeyQuotes:  r.KeyQuotes,
		Summary:    r.Summary,
		Confidence: r.Confidence,
	}
}

// toGraphQLExtractionMetadata converts domain ExtractionMetadata to GraphQL model.
func toGraphQLExtractionMetadata(m *domain.ExtractionMetadata) *model.ExtractionMetadata {
	if m == nil {
		return nil
	}
	return &model.ExtractionMetadata{
		ExtractedAt:       m.ExtractedAt,
		ModelVersion:      m.ModelVersion,
		OverallConfidence: m.OverallConfidence,
		ProcessingTimeMs:  m.ProcessingTimeMs,
		WarningsCount:     m.WarningsCount,
		Warnings:          m.Warnings,
	}
}
