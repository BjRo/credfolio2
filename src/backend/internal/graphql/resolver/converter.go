package resolver

import (
	"encoding/json"

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

	// Convert ExtractedData from json.RawMessage to map[string]any
	var extractedData map[string]any
	if len(rl.ExtractedData) > 0 {
		// Unmarshal JSON into map; ignore errors (leave nil on failure)
		_ = json.Unmarshal(rl.ExtractedData, &extractedData) //nolint:errcheck // Best effort parsing
	}

	return &model.ReferenceLetter{
		ID:            rl.ID.String(),
		Title:         rl.Title,
		AuthorName:    rl.AuthorName,
		AuthorTitle:   rl.AuthorTitle,
		Organization:  rl.Organization,
		DateWritten:   rl.DateWritten,
		RawText:       rl.RawText,
		ExtractedData: extractedData,
		Status:        status,
		CreatedAt:     rl.CreatedAt,
		UpdatedAt:     rl.UpdatedAt,
		User:          user,
		File:          file,
	}
}
