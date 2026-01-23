---
# credfolio2-si1z
title: Profile version snapshots
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:29:24Z
updated_at: 2026-01-23T16:29:50Z
parent: credfolio2-09w1
blocking:
    - credfolio2-ppen
    - credfolio2-1twf
---

Database and service layer for storing profile versions.

## Table: profile_versions

- id (UUID, PK)
- profile_id (FK)
- version_number (incrementing)
- snapshot_data (JSONB - full profile state)
- change_type (initial, extraction, manual_edit, merge, undo)
- change_source (resume, reference_letter_id, user)
- change_description (human-readable summary)
- created_at

## Service Methods

- CreateVersion(profileId, changeType, source) - Snapshot current state
- GetVersion(versionId) - Get specific version
- GetVersionHistory(profileId) - List all versions
- RestoreVersion(versionId) - Revert profile to this version

## Snapshot Strategy

Full snapshot approach (simpler):
- Store complete profile JSON at each version
- Easy to restore, compare
- More storage but profiles are small

## Checklist

- [ ] Create migration for profile_versions table
- [ ] Create Bun model
- [ ] Create repository
- [ ] Implement version service methods
- [ ] Add version creation hooks to mutations
- [ ] Write tests