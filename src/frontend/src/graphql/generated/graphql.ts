/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: string; output: string; }
  JSON: { input: Record<string, unknown>; output: Record<string, unknown>; }
  Upload: { input: any; output: any; }
};

/** Counts of items applied from a reference letter. */
export type AppliedCount = {
  __typename?: 'AppliedCount';
  /** Number of experience validations applied. */
  experienceValidations: Scalars['Int']['output'];
  /** Number of new skills added. */
  newSkills: Scalars['Int']['output'];
  /** Number of skill validations applied. */
  skillValidations: Scalars['Int']['output'];
  /** Number of testimonials added. */
  testimonials: Scalars['Int']['output'];
};

/** Error returned when applying validations fails. */
export type ApplyValidationsError = {
  __typename?: 'ApplyValidationsError';
  /** The field that caused the error, if applicable. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the failure. */
  message: Scalars['String']['output'];
};

/** Input for applying selected validations from a reference letter. */
export type ApplyValidationsInput = {
  /** Experience validations to apply. */
  experienceValidations: Array<ExperienceValidationInput>;
  /** New skills discovered in the reference letter to add to the profile. */
  newSkills: Array<NewSkillInput>;
  /** The reference letter ID to apply validations from. */
  referenceLetterID: Scalars['ID']['input'];
  /** Skill validations to apply. */
  skillValidations: Array<SkillValidationInput>;
  /** Testimonials to add to the profile. */
  testimonials: Array<TestimonialInput>;
};

/** Union type for apply validations result. */
export type ApplyValidationsResponse = ApplyValidationsError | ApplyValidationsResult;

/** Result of applying reference letter validations. */
export type ApplyValidationsResult = {
  __typename?: 'ApplyValidationsResult';
  /** Counts of items that were applied. */
  appliedCount: AppliedCount;
  /** The updated profile. */
  profile: Profile;
  /** The updated reference letter. */
  referenceLetter: ReferenceLetter;
};

/**
 * A person who provided testimonials for the profile. Authors are deduplicated
 * by name and company within a profile.
 */
export type Author = {
  __typename?: 'Author';
  /** Company/organization of the author. */
  company?: Maybe<Scalars['String']['output']>;
  /** When the author record was created. */
  createdAt: Scalars['DateTime']['output'];
  /** Unique identifier for the author. */
  id: Scalars['ID']['output'];
  /** Presigned URL for the author's profile image. */
  imageUrl?: Maybe<Scalars['String']['output']>;
  /** LinkedIn profile URL of the author. */
  linkedInUrl?: Maybe<Scalars['String']['output']>;
  /** Name of the author. */
  name: Scalars['String']['output'];
  /** All testimonials from this author. */
  testimonials: Array<Testimonial>;
  /** Title/position of the author. */
  title?: Maybe<Scalars['String']['output']>;
  /** When the author record was last updated. */
  updatedAt: Scalars['DateTime']['output'];
};

/** Relationship type between letter author and candidate. */
export enum AuthorRelationship {
  Client = 'CLIENT',
  Colleague = 'COLLEAGUE',
  DirectReport = 'DIRECT_REPORT',
  Manager = 'MANAGER',
  Mentor = 'MENTOR',
  Other = 'OTHER',
  Peer = 'PEER',
  Professor = 'PROFESSOR'
}

/** Input for creating a new education entry. */
export type CreateEducationInput = {
  /** Degree or certification (required). */
  degree: Scalars['String']['input'];
  /** Description or achievements. */
  description?: InputMaybe<Scalars['String']['input']>;
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: InputMaybe<Scalars['String']['input']>;
  /** Field of study. */
  field?: InputMaybe<Scalars['String']['input']>;
  /** GPA if applicable. */
  gpa?: InputMaybe<Scalars['String']['input']>;
  /** Institution name (required). */
  institution: Scalars['String']['input'];
  /** Whether currently studying here. */
  isCurrent: Scalars['Boolean']['input'];
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: InputMaybe<Scalars['String']['input']>;
};

/** Input for creating a new work experience. */
export type CreateExperienceInput = {
  /** Company or organization name (required). */
  company: Scalars['String']['input'];
  /** Job description or responsibilities. */
  description?: InputMaybe<Scalars['String']['input']>;
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: InputMaybe<Scalars['String']['input']>;
  /** Key achievements or highlights (bullet points). */
  highlights?: InputMaybe<Array<Scalars['String']['input']>>;
  /** Whether this is the current job. */
  isCurrent: Scalars['Boolean']['input'];
  /** Location of the job. */
  location?: InputMaybe<Scalars['String']['input']>;
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: InputMaybe<Scalars['String']['input']>;
  /** Job title or position (required). */
  title: Scalars['String']['input'];
};

/** Input for creating a new skill. */
export type CreateSkillInput = {
  /** Skill category (required). */
  category: SkillCategory;
  /** Skill name (required). */
  name: Scalars['String']['input'];
};

/** Union type for profile photo deletion result. */
export type DeleteProfilePhotoResponse = DeleteProfilePhotoResult | ProfileHeaderValidationError;

/** Result of deleting a profile photo. */
export type DeleteProfilePhotoResult = {
  __typename?: 'DeleteProfilePhotoResult';
  /** Whether the deletion was successful. */
  success: Scalars['Boolean']['output'];
};

/** Result of a delete operation. */
export type DeleteResult = {
  __typename?: 'DeleteResult';
  /** ID of the deleted item. */
  deletedId: Scalars['ID']['output'];
  /** Whether the deletion was successful. */
  success: Scalars['Boolean']['output'];
};

/** Status of asynchronous document detection. */
export enum DetectionStatus {
  Completed = 'COMPLETED',
  Failed = 'FAILED',
  Pending = 'PENDING',
  Processing = 'PROCESSING'
}

/** A skill discovered in a reference letter that may not be on the profile. */
export type DiscoveredSkill = {
  __typename?: 'DiscoveredSkill';
  /** Category assigned by the LLM (TECHNICAL, SOFT, or DOMAIN). */
  category: SkillCategory;
  /** Context for the skill mention (e.g., 'technical skills', 'leadership'). */
  context?: Maybe<Scalars['String']['output']>;
  /** Quote from the letter mentioning this skill. */
  quote: Scalars['String']['output'];
  /** The skill name. */
  skill: Scalars['String']['output'];
};

/**
 * Result of lightweight document content detection.
 * Used to quickly classify a document before running full extraction.
 */
export type DocumentDetectionResult = {
  __typename?: 'DocumentDetectionResult';
  /** Confidence in the detection (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Hint about the document type. */
  documentTypeHint: DocumentTypeHint;
  /** ID of the stored file for subsequent processing. */
  fileId: Scalars['ID']['output'];
  /** Whether the document contains career information (resume/CV content). */
  hasCareerInfo: Scalars['Boolean']['output'];
  /** Whether the document contains testimonial/recommendation content. */
  hasTestimonial: Scalars['Boolean']['output'];
  /** Brief summary of what was found in the document. */
  summary: Scalars['String']['output'];
  /** Name of the testimonial author (person who wrote the recommendation), if detected. */
  testimonialAuthor?: Maybe<Scalars['String']['output']>;
};

/**
 * Status of a document detection job.
 * Poll this query to track detection progress after uploadForDetection.
 */
export type DocumentDetectionStatus = {
  __typename?: 'DocumentDetectionStatus';
  /** Detection results, available when status is COMPLETED. */
  detection?: Maybe<DocumentDetectionResult>;
  /** Error message, set when status is FAILED. */
  error?: Maybe<Scalars['String']['output']>;
  /** The file ID being detected. */
  fileId: Scalars['ID']['output'];
  /** Current detection status. */
  status: DetectionStatus;
};

/** Input for reporting feedback about document detection or extraction. */
export type DocumentFeedbackInput = {
  /** The type of feedback being reported. */
  feedbackType: DocumentFeedbackType;
  /** ID of the file the feedback relates to. */
  fileId: Scalars['ID']['input'];
  /** Free-text message describing the issue. */
  message: Scalars['String']['input'];
};

/** Result of reporting document feedback. */
export type DocumentFeedbackResult = {
  __typename?: 'DocumentFeedbackResult';
  /** Whether the feedback was recorded successfully. */
  success: Scalars['Boolean']['output'];
};

/** The type of feedback being reported about document processing. */
export enum DocumentFeedbackType {
  /** Feedback about incorrect document type detection. */
  DetectionCorrection = 'DETECTION_CORRECTION',
  /** Feedback about extraction quality issues. */
  ExtractionQuality = 'EXTRACTION_QUALITY'
}

/** Aggregated processing status across resume and reference letter extraction. */
export type DocumentProcessingStatus = {
  __typename?: 'DocumentProcessingStatus';
  /** Overall status: true when all requested extractions are complete (COMPLETED or FAILED). */
  allComplete: Scalars['Boolean']['output'];
  /** Reference letter processing status (null if letter extraction was not requested). */
  referenceLetter?: Maybe<ReferenceLetter>;
  /** Resume processing status (null if resume extraction was not requested). */
  resume?: Maybe<Resume>;
};

/** Hint about the document type from lightweight detection. */
export enum DocumentTypeHint {
  Hybrid = 'HYBRID',
  ReferenceLetter = 'REFERENCE_LETTER',
  Resume = 'RESUME',
  Unknown = 'UNKNOWN'
}

/** Result when a duplicate file is detected during upload. */
export type DuplicateFileDetected = {
  __typename?: 'DuplicateFileDetected';
  /** The existing file that matches the uploaded content. */
  existingFile: File;
  /** The existing reference letter created from this file (if uploading a reference letter). */
  existingReferenceLetter?: Maybe<ReferenceLetter>;
  /** The existing resume created from this file (if uploading a resume). */
  existingResume?: Maybe<Resume>;
  /** Message describing the duplicate detection. */
  message: Scalars['String']['output'];
};

/** Union type for education create/update result. */
export type EducationResponse = EducationResult | EducationValidationError;

/** Result of a successful education operation. */
export type EducationResult = {
  __typename?: 'EducationResult';
  /** The created or updated education entry. */
  education: ProfileEducation;
};

/** Error returned when education validation fails. */
export type EducationValidationError = {
  __typename?: 'EducationValidationError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** Union type for experience create/update result. */
export type ExperienceResponse = ExperienceResult | ExperienceValidationError;

/** Result of a successful experience operation. */
export type ExperienceResult = {
  __typename?: 'ExperienceResult';
  /** The created or updated experience. */
  experience: ProfileExperience;
};

/** Source of a profile experience entry. */
export enum ExperienceSource {
  /** Discovered in a reference letter. */
  LetterDiscovered = 'LETTER_DISCOVERED',
  /** Manually entered by user. */
  Manual = 'MANUAL',
  /** Extracted from an uploaded resume. */
  ResumeExtracted = 'RESUME_EXTRACTED'
}

/** A validation record linking an experience to a reference letter. */
export type ExperienceValidation = {
  __typename?: 'ExperienceValidation';
  /** When the validation was created. */
  createdAt: Scalars['DateTime']['output'];
  /** The experience being validated. */
  experience: ProfileExperience;
  /** Unique identifier for the validation. */
  id: Scalars['ID']['output'];
  /** Quote snippet from the reference letter supporting this experience. */
  quoteSnippet?: Maybe<Scalars['String']['output']>;
  /** The reference letter providing the validation. */
  referenceLetter: ReferenceLetter;
};

/** Error returned when experience validation fails. */
export type ExperienceValidationError = {
  __typename?: 'ExperienceValidationError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** Input for applying an experience validation from a reference letter. */
export type ExperienceValidationInput = {
  /** The profile experience ID to validate. */
  profileExperienceID: Scalars['ID']['input'];
  /** Quote snippet from the reference letter supporting this experience. */
  quoteSnippet: Scalars['String']['input'];
};

/** Author details extracted from a reference letter. */
export type ExtractedAuthor = {
  __typename?: 'ExtractedAuthor';
  /** The company where the author works. */
  company?: Maybe<Scalars['String']['output']>;
  /** The author's full name. */
  name: Scalars['String']['output'];
  /** The relationship type between author and candidate. */
  relationship: AuthorRelationship;
  /** The author's job title or position. */
  title?: Maybe<Scalars['String']['output']>;
};

/** An education entry extracted from a resume. */
export type ExtractedEducation = {
  __typename?: 'ExtractedEducation';
  /** Achievements or honors. */
  achievements?: Maybe<Scalars['String']['output']>;
  /** Degree earned. */
  degree?: Maybe<Scalars['String']['output']>;
  /** End date. */
  endDate?: Maybe<Scalars['String']['output']>;
  /** Field of study. */
  field?: Maybe<Scalars['String']['output']>;
  /** Grade point average. */
  gpa?: Maybe<Scalars['String']['output']>;
  /** Institution name. */
  institution: Scalars['String']['output'];
  /** Start date. */
  startDate?: Maybe<Scalars['String']['output']>;
};

/** A reference to a role/company mentioned in the letter. */
export type ExtractedExperienceMention = {
  __typename?: 'ExtractedExperienceMention';
  /** Company name mentioned. */
  company: Scalars['String']['output'];
  /** Quote from the letter about this experience. */
  quote: Scalars['String']['output'];
  /** Role/title mentioned. */
  role: Scalars['String']['output'];
};

/** Structured data extracted from a reference letter for credibility validation. */
export type ExtractedLetterData = {
  __typename?: 'ExtractedLetterData';
  /** Details about the letter's author. */
  author: ExtractedAuthor;
  /** Skills discovered in the letter that aren't in the profile yet (with attribution). */
  discoveredSkills: Array<DiscoveredSkill>;
  /** Experience/role mentions with quotes. */
  experienceMentions: Array<ExtractedExperienceMention>;
  /** Metadata about the extraction process. */
  metadata: ExtractionMetadata;
  /** Specific skill mentions with quotes. */
  skillMentions: Array<ExtractedSkillMention>;
  /** Full testimonial quotes for display. */
  testimonials: Array<ExtractedTestimonial>;
};

/** A specific skill mentioned in the reference letter with context. */
export type ExtractedSkillMention = {
  __typename?: 'ExtractedSkillMention';
  /** Context for the skill mention (e.g., 'technical skills', 'leadership'). */
  context?: Maybe<Scalars['String']['output']>;
  /** Quote from the letter mentioning this skill. */
  quote: Scalars['String']['output'];
  /** The skill name. */
  skill: Scalars['String']['output'];
};

/** A testimonial quote suitable for display on the profile. */
export type ExtractedTestimonial = {
  __typename?: 'ExtractedTestimonial';
  /** The full quote text. */
  quote: Scalars['String']['output'];
  /** Skills mentioned in this testimonial. */
  skillsMentioned?: Maybe<Array<Scalars['String']['output']>>;
};

/** A work experience entry extracted from a resume. */
export type ExtractedWorkExperience = {
  __typename?: 'ExtractedWorkExperience';
  /** Company or organization name. */
  company: Scalars['String']['output'];
  /** Role description or achievements. */
  description?: Maybe<Scalars['String']['output']>;
  /** End date (YYYY-MM format), null if current. */
  endDate?: Maybe<Scalars['String']['output']>;
  /** Whether this is the current position. */
  isCurrent: Scalars['Boolean']['output'];
  /** Location of the role. */
  location?: Maybe<Scalars['String']['output']>;
  /** Start date (YYYY-MM format). */
  startDate?: Maybe<Scalars['String']['output']>;
  /** Job title or role. */
  title: Scalars['String']['output'];
};

/** Metadata about the extraction process. */
export type ExtractionMetadata = {
  __typename?: 'ExtractionMetadata';
  /** When the extraction was performed. */
  extractedAt: Scalars['DateTime']['output'];
  /** The LLM model version used for extraction. */
  modelVersion: Scalars['String']['output'];
  /** Time taken to process the extraction in milliseconds. */
  processingTimeMs?: Maybe<Scalars['Int']['output']>;
};

/** An uploaded file stored in object storage. */
export type File = {
  __typename?: 'File';
  /** SHA-256 hash of the file content for duplicate detection. */
  contentHash?: Maybe<Scalars['String']['output']>;
  contentType: Scalars['String']['output'];
  createdAt: Scalars['DateTime']['output'];
  filename: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  sizeBytes: Scalars['Int']['output'];
  storageKey: Scalars['String']['output'];
  /** Presigned URL for downloading the file. Expires after a short time. */
  url: Scalars['String']['output'];
  user: User;
};

/** Error returned when file validation fails. */
export type FileValidationError = {
  __typename?: 'FileValidationError';
  /** The field that failed validation (e.g., 'contentType', 'size'). */
  field: Scalars['String']['output'];
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** Error returned when import document results validation fails. */
export type ImportDocumentResultsError = {
  __typename?: 'ImportDocumentResultsError';
  /** The field that caused the error. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the failure. */
  message: Scalars['String']['output'];
};

/** Input for importing extracted document results into profile tables. */
export type ImportDocumentResultsInput = {
  /** Reference letter ID whose validations should be applied (null to skip). */
  referenceLetterID?: InputMaybe<Scalars['ID']['input']>;
  /** Resume ID to materialize into profile (null to skip). */
  resumeId?: InputMaybe<Scalars['ID']['input']>;
  /** Discovered skills to import from reference letter with per-skill category. Null = import all. */
  selectedDiscoveredSkills?: InputMaybe<Array<SelectedDiscoveredSkillInput>>;
  /** Indices of education entries to import. Null = import all. */
  selectedEducationIndices?: InputMaybe<Array<Scalars['Int']['input']>>;
  /** Indices of experiences to import. Null = import all. */
  selectedExperienceIndices?: InputMaybe<Array<Scalars['Int']['input']>>;
  /** Skill names to import. Null = import all. */
  selectedSkills?: InputMaybe<Array<Scalars['String']['input']>>;
  /** Indices of testimonials to import. Null = import all. */
  selectedTestimonialIndices?: InputMaybe<Array<Scalars['Int']['input']>>;
};

/** Union type for import document results. */
export type ImportDocumentResultsResponse = ImportDocumentResultsError | ImportDocumentResultsResult;

/** Result of importing document results into profile tables. */
export type ImportDocumentResultsResult = {
  __typename?: 'ImportDocumentResultsResult';
  /** Counts of items that were imported. */
  importedCount: ImportedCount;
  /** The updated profile after materialization. */
  profile: Profile;
};

/** Counts of items imported into the profile. */
export type ImportedCount = {
  __typename?: 'ImportedCount';
  /** Number of education entries materialized. */
  educations: Scalars['Int']['output'];
  /** Number of work experiences materialized. */
  experiences: Scalars['Int']['output'];
  /** Number of skills materialized. */
  skills: Scalars['Int']['output'];
  /** Number of testimonials materialized from reference letters. */
  testimonials: Scalars['Int']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  /**
   * Apply selected validations from a reference letter to the profile.
   * Creates skill validations, experience validations, testimonials, and new skills.
   * Updates the reference letter status to indicate validations have been applied.
   */
  applyReferenceLetterValidations: ApplyValidationsResponse;
  /**
   * Create a new education entry for a user's profile.
   * Creates the profile if it doesn't exist.
   */
  createEducation: EducationResponse;
  /**
   * Create a new work experience for a user's profile.
   * Creates the profile if it doesn't exist.
   */
  createExperience: ExperienceResponse;
  /**
   * Create a new skill for a user's profile.
   * Creates the profile if it doesn't exist.
   */
  createSkill: SkillResponse;
  /** Delete an education entry. */
  deleteEducation: DeleteResult;
  /** Delete a work experience. */
  deleteExperience: DeleteResult;
  /** Delete the profile photo. */
  deleteProfilePhoto: DeleteProfilePhotoResponse;
  /** Delete a skill. */
  deleteSkill: DeleteResult;
  /** Delete a testimonial. */
  deleteTestimonial: DeleteResult;
  /**
   * Import extracted document results into profile tables.
   * Materializes resume data (experiences, education, skills) and/or applies
   * reference letter validations. The referenced entities must be in COMPLETED status.
   */
  importDocumentResults: ImportDocumentResultsResponse;
  /**
   * Start processing a previously uploaded document for extraction.
   * The file must already exist (created via uploadForDetection).
   * Creates resume and/or reference letter records as needed, enqueues extraction,
   * and returns IDs for polling via documentProcessingStatus.
   * At least one of extractCareerInfo or extractTestimonial must be true.
   */
  processDocument: ProcessDocumentResponse;
  /**
   * Report feedback about document detection or extraction quality.
   * Feedback is logged for analysis; no table is created for MVP.
   */
  reportDocumentFeedback: DocumentFeedbackResult;
  /**
   * Update an author's information.
   * Allows correcting author details and adding LinkedIn profile links.
   */
  updateAuthor: Author;
  /**
   * Update an existing education entry.
   * Only updates fields that are provided.
   */
  updateEducation: EducationResponse;
  /**
   * Update an existing work experience.
   * Only updates fields that are provided.
   */
  updateExperience: ExperienceResponse;
  /**
   * Update the profile header fields (name, email, phone, location, summary).
   * Creates the profile if it doesn't exist.
   * Only updates fields that are provided.
   */
  updateProfileHeader: ProfileHeaderResponse;
  /**
   * Update an existing skill.
   * Only updates fields that are provided.
   */
  updateSkill: SkillResponse;
  /**
   * Upload an author's profile image.
   * Accepts JPEG, PNG, GIF, or WebP images (max 5MB).
   * Updates the author record with the new image.
   * Returns the updated author with image URL.
   */
  uploadAuthorImage: UploadAuthorImageResponse;
  /**
   * Upload a reference letter file for processing.
   * Accepts PDF, DOCX, or TXT files.
   * Creates a file record and queues the document for LLM extraction.
   * If a file with the same content hash already exists, returns DuplicateFileDetected
   * unless forceReimport is true.
   */
  uploadFile: UploadFileResponse;
  /**
   * Upload a document and start asynchronous content detection.
   * Stores the file and enqueues a background detection job.
   * Poll documentDetectionStatus to get detection results.
   */
  uploadForDetection: UploadForDetectionResponse;
  /**
   * Upload a profile photo.
   * Accepts JPEG, PNG, GIF, or WebP images (max 5MB).
   * Creates the profile if it doesn't exist.
   * Returns the updated profile with the photo URL.
   */
  uploadProfilePhoto: UploadProfilePhotoResponse;
  /**
   * Upload a resume file for processing.
   * Accepts PDF, DOCX, or TXT files.
   * Creates a file record and queues the resume for LLM extraction.
   * If a file with the same content hash already exists, returns DuplicateFileDetected
   * unless forceReimport is true.
   */
  uploadResume: UploadResumeResponse;
};


export type MutationApplyReferenceLetterValidationsArgs = {
  input: ApplyValidationsInput;
  userId: Scalars['ID']['input'];
};


export type MutationCreateEducationArgs = {
  input: CreateEducationInput;
  userId: Scalars['ID']['input'];
};


export type MutationCreateExperienceArgs = {
  input: CreateExperienceInput;
  userId: Scalars['ID']['input'];
};


export type MutationCreateSkillArgs = {
  input: CreateSkillInput;
  userId: Scalars['ID']['input'];
};


export type MutationDeleteEducationArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteExperienceArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteProfilePhotoArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationDeleteSkillArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteTestimonialArgs = {
  id: Scalars['ID']['input'];
};


export type MutationImportDocumentResultsArgs = {
  input: ImportDocumentResultsInput;
  userId: Scalars['ID']['input'];
};


export type MutationProcessDocumentArgs = {
  input: ProcessDocumentInput;
  userId: Scalars['ID']['input'];
};


export type MutationReportDocumentFeedbackArgs = {
  input: DocumentFeedbackInput;
  userId: Scalars['ID']['input'];
};


export type MutationUpdateAuthorArgs = {
  id: Scalars['ID']['input'];
  input: UpdateAuthorInput;
};


export type MutationUpdateEducationArgs = {
  id: Scalars['ID']['input'];
  input: UpdateEducationInput;
};


export type MutationUpdateExperienceArgs = {
  id: Scalars['ID']['input'];
  input: UpdateExperienceInput;
};


export type MutationUpdateProfileHeaderArgs = {
  input: UpdateProfileHeaderInput;
  userId: Scalars['ID']['input'];
};


export type MutationUpdateSkillArgs = {
  id: Scalars['ID']['input'];
  input: UpdateSkillInput;
};


export type MutationUploadAuthorImageArgs = {
  authorId: Scalars['ID']['input'];
  file: Scalars['Upload']['input'];
};


export type MutationUploadFileArgs = {
  file: Scalars['Upload']['input'];
  forceReimport?: InputMaybe<Scalars['Boolean']['input']>;
  userId: Scalars['ID']['input'];
};


export type MutationUploadForDetectionArgs = {
  file: Scalars['Upload']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationUploadProfilePhotoArgs = {
  file: Scalars['Upload']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationUploadResumeArgs = {
  file: Scalars['Upload']['input'];
  forceReimport?: InputMaybe<Scalars['Boolean']['input']>;
  userId: Scalars['ID']['input'];
};

/** Input for adding a new skill discovered in the reference letter. */
export type NewSkillInput = {
  /** The skill category. */
  category: SkillCategory;
  /** The skill name. */
  name: Scalars['String']['input'];
  /** Quote context from the reference letter mentioning this skill. */
  quoteContext?: InputMaybe<Scalars['String']['input']>;
};

/** Error returned when process document validation fails. */
export type ProcessDocumentError = {
  __typename?: 'ProcessDocumentError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/**
 * Input for processing a previously uploaded document.
 * The fileId must reference a file previously stored via uploadForDetection.
 */
export type ProcessDocumentInput = {
  /** Whether to extract career/resume information from the document. */
  extractCareerInfo: Scalars['Boolean']['input'];
  /** Whether to extract testimonial/reference letter content from the document. */
  extractTestimonial: Scalars['Boolean']['input'];
  /** ID of the file already stored via uploadForDetection. */
  fileId: Scalars['ID']['input'];
};

/** Union type for process document result. */
export type ProcessDocumentResponse = ProcessDocumentError | ProcessDocumentResult;

/** IDs of the entities created for tracking extraction progress. */
export type ProcessDocumentResult = {
  __typename?: 'ProcessDocumentResult';
  /** Reference letter ID created for testimonial extraction (null if not requested). */
  referenceLetterID?: Maybe<Scalars['ID']['output']>;
  /** Resume ID created for career info extraction (null if not requested). */
  resumeId?: Maybe<Scalars['ID']['output']>;
};

/** A user's profile containing manually editable data. */
export type Profile = {
  __typename?: 'Profile';
  createdAt: Scalars['DateTime']['output'];
  /** Education entries. */
  educations: Array<ProfileEducation>;
  /** User-edited email (overrides resume extraction if set). */
  email?: Maybe<Scalars['String']['output']>;
  /** Work experience entries. */
  experiences: Array<ProfileExperience>;
  id: Scalars['ID']['output'];
  /** User-edited location (overrides resume extraction if set). */
  location?: Maybe<Scalars['String']['output']>;
  /** User-edited name (overrides resume extraction if set). */
  name?: Maybe<Scalars['String']['output']>;
  /** User-edited phone (overrides resume extraction if set). */
  phone?: Maybe<Scalars['String']['output']>;
  /** URL to the user's profile photo (presigned URL for direct access). */
  profilePhotoUrl?: Maybe<Scalars['String']['output']>;
  /** Skill entries. */
  skills: Array<ProfileSkill>;
  /** User-edited professional summary (overrides resume extraction if set). */
  summary?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['DateTime']['output'];
  /** The user who owns this profile. */
  user: User;
};

/** An education entry in a user's profile. */
export type ProfileEducation = {
  __typename?: 'ProfileEducation';
  createdAt: Scalars['DateTime']['output'];
  /** Degree or certification. */
  degree: Scalars['String']['output'];
  /** Description or achievements. */
  description?: Maybe<Scalars['String']['output']>;
  /** Display order for sorting. */
  displayOrder: Scalars['Int']['output'];
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: Maybe<Scalars['String']['output']>;
  /** Field of study. */
  field?: Maybe<Scalars['String']['output']>;
  /** GPA if applicable. */
  gpa?: Maybe<Scalars['String']['output']>;
  /** Unique identifier for the education entry. */
  id: Scalars['ID']['output'];
  /** Institution name. */
  institution: Scalars['String']['output'];
  /** Whether currently studying here. */
  isCurrent: Scalars['Boolean']['output'];
  /** Source of this education entry. */
  source: ExperienceSource;
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['DateTime']['output'];
};

/** A work experience entry in a user's profile. */
export type ProfileExperience = {
  __typename?: 'ProfileExperience';
  /** Company or organization name. */
  company: Scalars['String']['output'];
  createdAt: Scalars['DateTime']['output'];
  /** Job description or responsibilities. */
  description?: Maybe<Scalars['String']['output']>;
  /** Display order for sorting. */
  displayOrder: Scalars['Int']['output'];
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: Maybe<Scalars['String']['output']>;
  /** Key achievements or highlights (bullet points). */
  highlights: Array<Scalars['String']['output']>;
  /** Unique identifier for the experience. */
  id: Scalars['ID']['output'];
  /** Whether this is the current job. */
  isCurrent: Scalars['Boolean']['output'];
  /** Location of the job. */
  location?: Maybe<Scalars['String']['output']>;
  /** Source of this experience entry. */
  source: ExperienceSource;
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: Maybe<Scalars['String']['output']>;
  /** Job title or position. */
  title: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
  /** Number of reference letters validating this experience. */
  validationCount: Scalars['Int']['output'];
};

/** Union type for profile header update result. */
export type ProfileHeaderResponse = ProfileHeaderResult | ProfileHeaderValidationError;

/** Result of a successful profile header update. */
export type ProfileHeaderResult = {
  __typename?: 'ProfileHeaderResult';
  /** The updated profile. */
  profile: Profile;
};

/** Error returned when profile header validation fails. */
export type ProfileHeaderValidationError = {
  __typename?: 'ProfileHeaderValidationError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** A skill entry in a user's profile. */
export type ProfileSkill = {
  __typename?: 'ProfileSkill';
  /** Skill category. */
  category: SkillCategory;
  createdAt: Scalars['DateTime']['output'];
  /** Display order for sorting. */
  displayOrder: Scalars['Int']['output'];
  /** Unique identifier for the skill. */
  id: Scalars['ID']['output'];
  /** Skill name as displayed. */
  name: Scalars['String']['output'];
  /** Normalized skill name for deduplication. */
  normalizedName: Scalars['String']['output'];
  /** Source of this skill entry. */
  source: ExperienceSource;
  /** The reference letter this skill was discovered from (if source is reference letter). */
  sourceReferenceLetter?: Maybe<ReferenceLetter>;
  updatedAt: Scalars['DateTime']['output'];
  /** Number of reference letters validating this skill. */
  validationCount: Scalars['Int']['output'];
};

export type Query = {
  __typename?: 'Query';
  /** Get an author by ID. */
  author?: Maybe<Author>;
  /** Get all authors for a profile. */
  authors: Array<Author>;
  /**
   * Check if a file with the given content hash already exists for the user.
   * Returns the existing file and associated resume/reference letter if found.
   * Used for pre-upload duplicate detection.
   */
  checkDuplicateFile?: Maybe<DuplicateFileDetected>;
  /**
   * Get the detection status for an uploaded document.
   * Poll this after uploadForDetection to get detection results.
   */
  documentDetectionStatus?: Maybe<DocumentDetectionStatus>;
  /**
   * Get the processing status of a document.
   * Provide the resume ID and/or reference letter ID returned by processDocument.
   * Returns aggregated status across all requested extractions.
   */
  documentProcessingStatus?: Maybe<DocumentProcessingStatus>;
  /** Get all validations for a specific experience. */
  experienceValidations: Array<ExperienceValidation>;
  /** Get a file by ID. */
  file?: Maybe<File>;
  /** Get all files for a user. */
  files: Array<File>;
  /**
   * Get a profile by its ID.
   * Returns null if no profile exists.
   */
  profile?: Maybe<Profile>;
  /**
   * Get a user's profile by user ID.
   * Returns null if no profile exists.
   */
  profileByUserId?: Maybe<Profile>;
  /** Get a single profile education entry by ID. */
  profileEducation?: Maybe<ProfileEducation>;
  /** Get a single profile experience by ID. */
  profileExperience?: Maybe<ProfileExperience>;
  /** Get a single profile skill by ID. */
  profileSkill?: Maybe<ProfileSkill>;
  /** Get a reference letter by ID. */
  referenceLetter?: Maybe<ReferenceLetter>;
  /** Get all reference letters for a user. */
  referenceLetters: Array<ReferenceLetter>;
  /** Get a resume by ID. */
  resume?: Maybe<Resume>;
  /** Get all resumes for a user. */
  resumes: Array<Resume>;
  /** Get all validations for a specific skill. */
  skillValidations: Array<SkillValidation>;
  /** Get all testimonials for a profile. */
  testimonials: Array<Testimonial>;
  /** Get a user by ID. */
  user?: Maybe<User>;
};


export type QueryAuthorArgs = {
  id: Scalars['ID']['input'];
};


export type QueryAuthorsArgs = {
  profileId: Scalars['ID']['input'];
};


export type QueryCheckDuplicateFileArgs = {
  contentHash: Scalars['String']['input'];
  userId: Scalars['ID']['input'];
};


export type QueryDocumentDetectionStatusArgs = {
  fileId: Scalars['ID']['input'];
};


export type QueryDocumentProcessingStatusArgs = {
  referenceLetterID?: InputMaybe<Scalars['ID']['input']>;
  resumeId?: InputMaybe<Scalars['ID']['input']>;
};


export type QueryExperienceValidationsArgs = {
  experienceId: Scalars['ID']['input'];
};


export type QueryFileArgs = {
  id: Scalars['ID']['input'];
};


export type QueryFilesArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryProfileArgs = {
  id: Scalars['ID']['input'];
};


export type QueryProfileByUserIdArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryProfileEducationArgs = {
  id: Scalars['ID']['input'];
};


export type QueryProfileExperienceArgs = {
  id: Scalars['ID']['input'];
};


export type QueryProfileSkillArgs = {
  id: Scalars['ID']['input'];
};


export type QueryReferenceLetterArgs = {
  id: Scalars['ID']['input'];
};


export type QueryReferenceLettersArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryResumeArgs = {
  id: Scalars['ID']['input'];
};


export type QueryResumesArgs = {
  userId: Scalars['ID']['input'];
};


export type QuerySkillValidationsArgs = {
  skillId: Scalars['ID']['input'];
};


export type QueryTestimonialsArgs = {
  profileId: Scalars['ID']['input'];
};


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};

/** A reference letter with extracted data. */
export type ReferenceLetter = {
  __typename?: 'ReferenceLetter';
  authorName?: Maybe<Scalars['String']['output']>;
  authorTitle?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['DateTime']['output'];
  dateWritten?: Maybe<Scalars['DateTime']['output']>;
  /** Structured data extracted from the letter by LLM processing. */
  extractedData?: Maybe<ExtractedLetterData>;
  file?: Maybe<File>;
  id: Scalars['ID']['output'];
  organization?: Maybe<Scalars['String']['output']>;
  rawText?: Maybe<Scalars['String']['output']>;
  status: ReferenceLetterStatus;
  title?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['DateTime']['output'];
  user: User;
};

/** Processing status of a reference letter. */
export enum ReferenceLetterStatus {
  Applied = 'APPLIED',
  Completed = 'COMPLETED',
  Failed = 'FAILED',
  Pending = 'PENDING',
  Processing = 'PROCESSING'
}

/** An uploaded resume with extracted profile data. */
export type Resume = {
  __typename?: 'Resume';
  createdAt: Scalars['DateTime']['output'];
  /** Error message if processing failed. */
  errorMessage?: Maybe<Scalars['String']['output']>;
  /** Structured data extracted from the resume by LLM processing. */
  extractedData?: Maybe<ResumeExtractedData>;
  file: File;
  id: Scalars['ID']['output'];
  /** Processing status of the resume. */
  status: ResumeStatus;
  updatedAt: Scalars['DateTime']['output'];
  user: User;
};

/**
 * Structured data extracted from a resume.
 * Education and work experience are materialized into profile tables
 * and exposed via ProfileEducation/ProfileExperience types.
 */
export type ResumeExtractedData = {
  __typename?: 'ResumeExtractedData';
  /** Overall confidence score (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Education entries. */
  educations: Array<ExtractedEducation>;
  /** Email address. */
  email?: Maybe<Scalars['String']['output']>;
  /** Work experience entries. */
  experiences: Array<ExtractedWorkExperience>;
  /** When the extraction was performed. */
  extractedAt: Scalars['DateTime']['output'];
  /** Location (city, state, country). */
  location?: Maybe<Scalars['String']['output']>;
  /** Full name of the candidate. */
  name: Scalars['String']['output'];
  /** Phone number. */
  phone?: Maybe<Scalars['String']['output']>;
  /** Skill names. */
  skills: Array<Scalars['String']['output']>;
  /** Professional summary or objective. */
  summary?: Maybe<Scalars['String']['output']>;
};

/** Processing status of a resume. */
export enum ResumeStatus {
  Completed = 'COMPLETED',
  Failed = 'FAILED',
  Pending = 'PENDING',
  Processing = 'PROCESSING'
}

/** Input for a discovered skill selected for import, carrying both name and category. */
export type SelectedDiscoveredSkillInput = {
  /** The skill category (user may override the LLM-assigned category). */
  category: SkillCategory;
  /** The skill name. */
  name: Scalars['String']['input'];
};

/** Skill category classification. */
export enum SkillCategory {
  Domain = 'DOMAIN',
  Soft = 'SOFT',
  Technical = 'TECHNICAL'
}

/** Union type for skill create/update result. */
export type SkillResponse = SkillResult | SkillValidationError;

/** Result of a successful skill operation. */
export type SkillResult = {
  __typename?: 'SkillResult';
  /** The created or updated skill. */
  skill: ProfileSkill;
};

/** A validation record linking a skill to a reference letter. */
export type SkillValidation = {
  __typename?: 'SkillValidation';
  /** When the validation was created. */
  createdAt: Scalars['DateTime']['output'];
  /** Unique identifier for the validation. */
  id: Scalars['ID']['output'];
  /** Quote snippet from the reference letter supporting this skill. */
  quoteSnippet?: Maybe<Scalars['String']['output']>;
  /** The reference letter providing the validation. */
  referenceLetter: ReferenceLetter;
  /** The skill being validated. */
  skill: ProfileSkill;
  /** The specific testimonial that validates this skill (if applicable). */
  testimonial?: Maybe<Testimonial>;
};

/** Error returned when skill validation fails. */
export type SkillValidationError = {
  __typename?: 'SkillValidationError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** Input for applying a skill validation from a reference letter. */
export type SkillValidationInput = {
  /** The profile skill ID to validate. */
  profileSkillID: Scalars['ID']['input'];
  /** Quote snippet from the reference letter supporting this skill. */
  quoteSnippet: Scalars['String']['input'];
};

/** A testimonial quote from a reference letter displayed on the profile. */
export type Testimonial = {
  __typename?: 'Testimonial';
  /** The author who provided this testimonial. */
  author?: Maybe<Author>;
  /** Company/organization of the author. (Deprecated: use author.company) */
  authorCompany?: Maybe<Scalars['String']['output']>;
  /** Name of the person who provided the testimonial. (Deprecated: use author.name) */
  authorName: Scalars['String']['output'];
  /** Title/position of the author. (Deprecated: use author.title) */
  authorTitle?: Maybe<Scalars['String']['output']>;
  /** When the testimonial was created. */
  createdAt: Scalars['DateTime']['output'];
  /** Unique identifier for the testimonial. */
  id: Scalars['ID']['output'];
  /** The full quote text. */
  quote: Scalars['String']['output'];
  /** The reference letter this testimonial was extracted from. */
  referenceLetter?: Maybe<ReferenceLetter>;
  /** Relationship between the author and the profile owner. */
  relationship: TestimonialRelationship;
  /** Skills validated by this testimonial's reference letter. */
  validatedSkills: Array<ProfileSkill>;
};

/** Input for creating a testimonial from a reference letter. */
export type TestimonialInput = {
  /** The full quote text for the testimonial. */
  quote: Scalars['String']['input'];
  /** Skills mentioned in this testimonial. */
  skillsMentioned: Array<Scalars['String']['input']>;
};

/** Relationship type between testimonial author and the profile owner. */
export enum TestimonialRelationship {
  Client = 'CLIENT',
  DirectReport = 'DIRECT_REPORT',
  Manager = 'MANAGER',
  Other = 'OTHER',
  Peer = 'PEER'
}

/** Input for updating an author's information. */
export type UpdateAuthorInput = {
  /** Updated company/organization of the author. */
  company?: InputMaybe<Scalars['String']['input']>;
  /** ID of the uploaded image file for the author's profile picture. */
  imageId?: InputMaybe<Scalars['ID']['input']>;
  /** Updated LinkedIn profile URL of the author. */
  linkedInUrl?: InputMaybe<Scalars['String']['input']>;
  /** Updated name of the author. */
  name?: InputMaybe<Scalars['String']['input']>;
  /** Updated title/position of the author. */
  title?: InputMaybe<Scalars['String']['input']>;
};

/** Input for updating an existing education entry. */
export type UpdateEducationInput = {
  /** Degree or certification. */
  degree?: InputMaybe<Scalars['String']['input']>;
  /** Description or achievements. */
  description?: InputMaybe<Scalars['String']['input']>;
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: InputMaybe<Scalars['String']['input']>;
  /** Field of study. */
  field?: InputMaybe<Scalars['String']['input']>;
  /** GPA if applicable. */
  gpa?: InputMaybe<Scalars['String']['input']>;
  /** Institution name. */
  institution?: InputMaybe<Scalars['String']['input']>;
  /** Whether currently studying here. */
  isCurrent?: InputMaybe<Scalars['Boolean']['input']>;
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: InputMaybe<Scalars['String']['input']>;
};

/** Input for updating an existing work experience. */
export type UpdateExperienceInput = {
  /** Company or organization name. */
  company?: InputMaybe<Scalars['String']['input']>;
  /** Job description or responsibilities. */
  description?: InputMaybe<Scalars['String']['input']>;
  /** End date (e.g., 'Dec 2023', 'Present'). */
  endDate?: InputMaybe<Scalars['String']['input']>;
  /** Key achievements or highlights (bullet points). */
  highlights?: InputMaybe<Array<Scalars['String']['input']>>;
  /** Whether this is the current job. */
  isCurrent?: InputMaybe<Scalars['Boolean']['input']>;
  /** Location of the job. */
  location?: InputMaybe<Scalars['String']['input']>;
  /** Start date (e.g., 'Jan 2020', '2020'). */
  startDate?: InputMaybe<Scalars['String']['input']>;
  /** Job title or position. */
  title?: InputMaybe<Scalars['String']['input']>;
};

/** Input for updating profile header fields. */
export type UpdateProfileHeaderInput = {
  /** Email address. */
  email?: InputMaybe<Scalars['String']['input']>;
  /** Location (city, state, country). */
  location?: InputMaybe<Scalars['String']['input']>;
  /** Name to display. */
  name?: InputMaybe<Scalars['String']['input']>;
  /** Phone number. */
  phone?: InputMaybe<Scalars['String']['input']>;
  /** Professional summary. */
  summary?: InputMaybe<Scalars['String']['input']>;
};

/** Input for updating an existing skill. */
export type UpdateSkillInput = {
  /** Skill category. */
  category?: InputMaybe<SkillCategory>;
  /** Skill name. */
  name?: InputMaybe<Scalars['String']['input']>;
};

/** Union type for author image upload result. */
export type UploadAuthorImageResponse = FileValidationError | UploadAuthorImageResult;

/** Result of a successful author image upload. */
export type UploadAuthorImageResult = {
  __typename?: 'UploadAuthorImageResult';
  /** The updated author with image URL. */
  author: Author;
  /** The uploaded file metadata. */
  file: File;
};

/** Union type for upload result - either success, validation error, or duplicate detected. */
export type UploadFileResponse = DuplicateFileDetected | FileValidationError | UploadFileResult;

/** Result of a file upload operation. */
export type UploadFileResult = {
  __typename?: 'UploadFileResult';
  /** The uploaded file metadata. */
  file: File;
  /** The reference letter created for processing. */
  referenceLetter: ReferenceLetter;
};

/** Union type for upload for detection result. */
export type UploadForDetectionResponse = FileValidationError | UploadForDetectionResult;

/**
 * Result of a successful document upload for detection.
 * Returns the file ID for polling detection status.
 */
export type UploadForDetectionResult = {
  __typename?: 'UploadForDetectionResult';
  /** The ID of the uploaded file. */
  fileId: Scalars['ID']['output'];
};

/** Union type for profile photo upload result. */
export type UploadProfilePhotoResponse = FileValidationError | UploadProfilePhotoResult;

/** Result of a successful profile photo upload. */
export type UploadProfilePhotoResult = {
  __typename?: 'UploadProfilePhotoResult';
  /** The updated profile with photo URL. */
  profile: Profile;
};

/** Union type for resume upload result - either success, validation error, or duplicate detected. */
export type UploadResumeResponse = DuplicateFileDetected | FileValidationError | UploadResumeResult;

/** Result of a resume upload operation. */
export type UploadResumeResult = {
  __typename?: 'UploadResumeResult';
  /** The uploaded file metadata. */
  file: File;
  /** The resume created for processing. */
  resume: Resume;
};

/** A user account in the system. */
export type User = {
  __typename?: 'User';
  createdAt: Scalars['DateTime']['output'];
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['DateTime']['output'];
};

export type CreateExperienceMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: CreateExperienceInput;
}>;


export type CreateExperienceMutation = { __typename?: 'Mutation', createExperience:
    | { __typename?: 'ExperienceResult', experience: { __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'ExperienceValidationError', message: string, field?: string | null }
   };

export type UpdateExperienceMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateExperienceInput;
}>;


export type UpdateExperienceMutation = { __typename?: 'Mutation', updateExperience:
    | { __typename?: 'ExperienceResult', experience: { __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'ExperienceValidationError', message: string, field?: string | null }
   };

export type DeleteExperienceMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteExperienceMutation = { __typename?: 'Mutation', deleteExperience: { __typename?: 'DeleteResult', success: boolean, deletedId: string } };

export type CreateEducationMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: CreateEducationInput;
}>;


export type CreateEducationMutation = { __typename?: 'Mutation', createEducation:
    | { __typename?: 'EducationResult', education: { __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'EducationValidationError', message: string, field?: string | null }
   };

export type UpdateEducationMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateEducationInput;
}>;


export type UpdateEducationMutation = { __typename?: 'Mutation', updateEducation:
    | { __typename?: 'EducationResult', education: { __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'EducationValidationError', message: string, field?: string | null }
   };

export type DeleteEducationMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteEducationMutation = { __typename?: 'Mutation', deleteEducation: { __typename?: 'DeleteResult', success: boolean, deletedId: string } };

export type CreateSkillMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: CreateSkillInput;
}>;


export type CreateSkillMutation = { __typename?: 'Mutation', createSkill:
    | { __typename?: 'SkillResult', skill: { __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'SkillValidationError', message: string, field?: string | null }
   };

export type UpdateSkillMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateSkillInput;
}>;


export type UpdateSkillMutation = { __typename?: 'Mutation', updateSkill:
    | { __typename?: 'SkillResult', skill: { __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string } }
    | { __typename?: 'SkillValidationError', message: string, field?: string | null }
   };

export type DeleteSkillMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteSkillMutation = { __typename?: 'Mutation', deleteSkill: { __typename?: 'DeleteResult', success: boolean, deletedId: string } };

export type UploadReferenceLetterMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  file: Scalars['Upload']['input'];
}>;


export type UploadReferenceLetterMutation = { __typename?: 'Mutation', uploadFile:
    | { __typename?: 'DuplicateFileDetected' }
    | { __typename: 'FileValidationError', message: string, field: string }
    | { __typename: 'UploadFileResult', file: { __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number }, referenceLetter: { __typename?: 'ReferenceLetter', id: string, status: ReferenceLetterStatus } }
   };

export type ApplyReferenceLetterValidationsMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: ApplyValidationsInput;
}>;


export type ApplyReferenceLetterValidationsMutation = { __typename?: 'Mutation', applyReferenceLetterValidations:
    | { __typename: 'ApplyValidationsError', message: string, field?: string | null }
    | { __typename: 'ApplyValidationsResult', referenceLetter: { __typename?: 'ReferenceLetter', id: string, status: ReferenceLetterStatus }, profile: { __typename?: 'Profile', id: string, skills: Array<{ __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory }>, experiences: Array<{ __typename?: 'ProfileExperience', id: string, company: string, title: string }> }, appliedCount: { __typename?: 'AppliedCount', skillValidations: number, experienceValidations: number, testimonials: number, newSkills: number } }
   };

export type UpdateProfileHeaderMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: UpdateProfileHeaderInput;
}>;


export type UpdateProfileHeaderMutation = { __typename?: 'Mutation', updateProfileHeader:
    | { __typename: 'ProfileHeaderResult', profile: { __typename?: 'Profile', id: string, name?: string | null, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, profilePhotoUrl?: string | null, updatedAt: string } }
    | { __typename: 'ProfileHeaderValidationError', message: string, field?: string | null }
   };

export type UploadProfilePhotoMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  file: Scalars['Upload']['input'];
}>;


export type UploadProfilePhotoMutation = { __typename?: 'Mutation', uploadProfilePhoto:
    | { __typename: 'FileValidationError', message: string, field: string }
    | { __typename: 'UploadProfilePhotoResult', profile: { __typename?: 'Profile', id: string, profilePhotoUrl?: string | null, updatedAt: string } }
   };

export type DeleteProfilePhotoMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type DeleteProfilePhotoMutation = { __typename?: 'Mutation', deleteProfilePhoto:
    | { __typename: 'DeleteProfilePhotoResult', success: boolean }
    | { __typename: 'ProfileHeaderValidationError', message: string, field?: string | null }
   };

export type DeleteTestimonialMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteTestimonialMutation = { __typename?: 'Mutation', deleteTestimonial: { __typename?: 'DeleteResult', success: boolean, deletedId: string } };

export type UpdateAuthorMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: UpdateAuthorInput;
}>;


export type UpdateAuthorMutation = { __typename?: 'Mutation', updateAuthor: { __typename?: 'Author', id: string, name: string, title?: string | null, company?: string | null, linkedInUrl?: string | null, imageUrl?: string | null, updatedAt: string } };

export type UploadAuthorImageMutationVariables = Exact<{
  authorId: Scalars['ID']['input'];
  file: Scalars['Upload']['input'];
}>;


export type UploadAuthorImageMutation = { __typename?: 'Mutation', uploadAuthorImage:
    | { __typename: 'FileValidationError', message: string, field: string }
    | { __typename: 'UploadAuthorImageResult', file: { __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number }, author: { __typename?: 'Author', id: string, name: string, title?: string | null, company?: string | null, linkedInUrl?: string | null, imageUrl?: string | null, updatedAt: string } }
   };

export type UploadForDetectionMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  file: Scalars['Upload']['input'];
}>;


export type UploadForDetectionMutation = { __typename?: 'Mutation', uploadForDetection:
    | { __typename: 'FileValidationError', message: string, field: string }
    | { __typename: 'UploadForDetectionResult', fileId: string }
   };

export type ReportDocumentFeedbackMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: DocumentFeedbackInput;
}>;


export type ReportDocumentFeedbackMutation = { __typename?: 'Mutation', reportDocumentFeedback: { __typename?: 'DocumentFeedbackResult', success: boolean } };

export type ProcessDocumentMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: ProcessDocumentInput;
}>;


export type ProcessDocumentMutation = { __typename?: 'Mutation', processDocument:
    | { __typename: 'ProcessDocumentError', message: string, field?: string | null }
    | { __typename: 'ProcessDocumentResult', resumeId?: string | null, referenceLetterID?: string | null }
   };

export type ImportDocumentResultsMutationVariables = Exact<{
  userId: Scalars['ID']['input'];
  input: ImportDocumentResultsInput;
}>;


export type ImportDocumentResultsMutation = { __typename?: 'Mutation', importDocumentResults:
    | { __typename: 'ImportDocumentResultsError', message: string, field?: string | null }
    | { __typename: 'ImportDocumentResultsResult', profile: { __typename?: 'Profile', id: string }, importedCount: { __typename?: 'ImportedCount', experiences: number, educations: number, skills: number, testimonials: number } }
   };

export type GetReferenceLetterForViewerQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetReferenceLetterForViewerQuery = { __typename?: 'Query', referenceLetter?: { __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, file?: { __typename?: 'File', id: string, url: string, filename: string, contentType: string } | null } | null };

export type GetReferenceLetterQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetReferenceLetterQuery = { __typename?: 'Query', referenceLetter?: { __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, dateWritten?: string | null, rawText?: string | null, status: ReferenceLetterStatus, createdAt: string, updatedAt: string, extractedData?: { __typename?: 'ExtractedLetterData', author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, company?: string | null, relationship: AuthorRelationship }, testimonials: Array<{ __typename?: 'ExtractedTestimonial', quote: string, skillsMentioned?: Array<string> | null }>, skillMentions: Array<{ __typename?: 'ExtractedSkillMention', skill: string, quote: string, context?: string | null }>, experienceMentions: Array<{ __typename?: 'ExtractedExperienceMention', company: string, role: string, quote: string }>, discoveredSkills: Array<{ __typename?: 'DiscoveredSkill', skill: string, quote: string, context?: string | null, category: SkillCategory }>, metadata: { __typename?: 'ExtractionMetadata', extractedAt: string, modelVersion: string, processingTimeMs?: number | null } } | null, user: { __typename?: 'User', id: string, email: string, name?: string | null }, file?: { __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number } | null } | null };

export type GetReferenceLettersQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetReferenceLettersQuery = { __typename?: 'Query', referenceLetters: Array<{ __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, status: ReferenceLetterStatus, createdAt: string }> };

export type GetReferenceLetterStatusQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetReferenceLetterStatusQuery = { __typename?: 'Query', referenceLetter?: { __typename?: 'ReferenceLetter', id: string, status: ReferenceLetterStatus } | null };

export type GetUserQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetUserQuery = { __typename?: 'Query', user?: { __typename?: 'User', id: string, email: string, name?: string | null, createdAt: string, updatedAt: string } | null };

export type GetFilesQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetFilesQuery = { __typename?: 'Query', files: Array<{ __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number, storageKey: string, createdAt: string }> };

export type GetResumeQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetResumeQuery = { __typename?: 'Query', resume?: { __typename?: 'Resume', id: string, status: ResumeStatus, errorMessage?: string | null, createdAt: string, updatedAt: string, extractedData?: { __typename?: 'ResumeExtractedData', name: string, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, extractedAt: string, confidence: number } | null, user: { __typename?: 'User', id: string, email: string, name?: string | null }, file: { __typename?: 'File', id: string, filename: string, contentType: string } } | null };

export type GetUserResumesQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetUserResumesQuery = { __typename?: 'Query', resumes: Array<{ __typename?: 'Resume', id: string, status: ResumeStatus }> };

export type GetProfileQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetProfileQuery = { __typename?: 'Query', profileByUserId?: { __typename?: 'Profile', id: string, name?: string | null, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, profilePhotoUrl?: string | null, createdAt: string, updatedAt: string, experiences: Array<{ __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, validationCount: number, createdAt: string, updatedAt: string }>, educations: Array<{ __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }>, skills: Array<{ __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory, displayOrder: number, source: ExperienceSource, validationCount: number, createdAt: string, updatedAt: string }> } | null };

export type GetProfileByIdQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetProfileByIdQuery = { __typename?: 'Query', profile?: { __typename?: 'Profile', id: string, name?: string | null, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, profilePhotoUrl?: string | null, createdAt: string, updatedAt: string, user: { __typename?: 'User', id: string }, experiences: Array<{ __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, validationCount: number, createdAt: string, updatedAt: string }>, educations: Array<{ __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }>, skills: Array<{ __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory, displayOrder: number, source: ExperienceSource, validationCount: number, createdAt: string, updatedAt: string }> } | null };

export type GetTestimonialsQueryVariables = Exact<{
  profileId: Scalars['ID']['input'];
}>;


export type GetTestimonialsQuery = { __typename?: 'Query', testimonials: Array<{ __typename?: 'Testimonial', id: string, quote: string, authorName: string, authorTitle?: string | null, authorCompany?: string | null, relationship: TestimonialRelationship, createdAt: string, author?: { __typename?: 'Author', id: string, name: string, title?: string | null, company?: string | null, linkedInUrl?: string | null, imageUrl?: string | null } | null, validatedSkills: Array<{ __typename?: 'ProfileSkill', id: string, name: string }>, referenceLetter?: { __typename?: 'ReferenceLetter', id: string, file?: { __typename?: 'File', id: string, url: string } | null } | null }> };

export type GetSkillValidationsQueryVariables = Exact<{
  skillId: Scalars['ID']['input'];
}>;


export type GetSkillValidationsQuery = { __typename?: 'Query', skillValidations: Array<{ __typename?: 'SkillValidation', id: string, quoteSnippet?: string | null, createdAt: string, referenceLetter: { __typename?: 'ReferenceLetter', id: string, file?: { __typename?: 'File', id: string, url: string } | null, extractedData?: { __typename?: 'ExtractedLetterData', author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, company?: string | null, relationship: AuthorRelationship } } | null } }> };

export type GetExperienceValidationsQueryVariables = Exact<{
  experienceId: Scalars['ID']['input'];
}>;


export type GetExperienceValidationsQuery = { __typename?: 'Query', experienceValidations: Array<{ __typename?: 'ExperienceValidation', id: string, quoteSnippet?: string | null, createdAt: string, referenceLetter: { __typename?: 'ReferenceLetter', id: string, file?: { __typename?: 'File', id: string, url: string } | null, extractedData?: { __typename?: 'ExtractedLetterData', author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, company?: string | null, relationship: AuthorRelationship } } | null } }> };

export type GetDocumentDetectionStatusQueryVariables = Exact<{
  fileId: Scalars['ID']['input'];
}>;


export type GetDocumentDetectionStatusQuery = { __typename?: 'Query', documentDetectionStatus?: { __typename?: 'DocumentDetectionStatus', fileId: string, status: DetectionStatus, error?: string | null, detection?: { __typename?: 'DocumentDetectionResult', hasCareerInfo: boolean, hasTestimonial: boolean, testimonialAuthor?: string | null, confidence: number, summary: string, documentTypeHint: DocumentTypeHint, fileId: string } | null } | null };

export type GetDocumentProcessingStatusQueryVariables = Exact<{
  resumeId?: InputMaybe<Scalars['ID']['input']>;
  referenceLetterID?: InputMaybe<Scalars['ID']['input']>;
}>;


export type GetDocumentProcessingStatusQuery = { __typename?: 'Query', documentProcessingStatus?: { __typename?: 'DocumentProcessingStatus', allComplete: boolean, resume?: { __typename?: 'Resume', id: string, status: ResumeStatus, errorMessage?: string | null, extractedData?: { __typename?: 'ResumeExtractedData', name: string, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, extractedAt: string, confidence: number } | null } | null, referenceLetter?: { __typename?: 'ReferenceLetter', id: string, status: ReferenceLetterStatus, extractedData?: { __typename?: 'ExtractedLetterData', author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, company?: string | null, relationship: AuthorRelationship }, testimonials: Array<{ __typename?: 'ExtractedTestimonial', quote: string, skillsMentioned?: Array<string> | null }>, skillMentions: Array<{ __typename?: 'ExtractedSkillMention', skill: string, quote: string, context?: string | null }>, experienceMentions: Array<{ __typename?: 'ExtractedExperienceMention', company: string, role: string, quote: string }>, discoveredSkills: Array<{ __typename?: 'DiscoveredSkill', skill: string, quote: string, context?: string | null, category: SkillCategory }>, metadata: { __typename?: 'ExtractionMetadata', extractedAt: string, modelVersion: string, processingTimeMs?: number | null } } | null } | null } | null };


export const CreateExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateExperienceInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experience"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<CreateExperienceMutation, CreateExperienceMutationVariables>;
export const UpdateExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateExperienceInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experience"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateExperienceMutation, UpdateExperienceMutationVariables>;
export const DeleteExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteExperienceMutation, DeleteExperienceMutationVariables>;
export const CreateEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateEducationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"education"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<CreateEducationMutation, CreateEducationMutationVariables>;
export const UpdateEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateEducationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"education"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateEducationMutation, UpdateEducationMutationVariables>;
export const DeleteEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteEducationMutation, DeleteEducationMutationVariables>;
export const CreateSkillDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateSkill"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateSkillInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createSkill"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"SkillResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"SkillValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<CreateSkillMutation, CreateSkillMutationVariables>;
export const UpdateSkillDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateSkill"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateSkillInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateSkill"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"SkillResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"SkillValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateSkillMutation, UpdateSkillMutationVariables>;
export const DeleteSkillDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteSkill"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteSkill"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteSkillMutation, DeleteSkillMutationVariables>;
export const UploadReferenceLetterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UploadReferenceLetter"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"file"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Upload"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uploadFile"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"file"},"value":{"kind":"Variable","name":{"kind":"Name","value":"file"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"UploadFileResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"FileValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UploadReferenceLetterMutation, UploadReferenceLetterMutationVariables>;
export const ApplyReferenceLetterValidationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ApplyReferenceLetterValidations"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ApplyValidationsInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"applyReferenceLetterValidations"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ApplyValidationsResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}},{"kind":"Field","name":{"kind":"Name","value":"profile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"experiences"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"appliedCount"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skillValidations"}},{"kind":"Field","name":{"kind":"Name","value":"experienceValidations"}},{"kind":"Field","name":{"kind":"Name","value":"testimonials"}},{"kind":"Field","name":{"kind":"Name","value":"newSkills"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ApplyValidationsError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<ApplyReferenceLetterValidationsMutation, ApplyReferenceLetterValidationsMutationVariables>;
export const UpdateProfileHeaderDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateProfileHeader"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateProfileHeaderInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateProfileHeader"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ProfileHeaderResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"profile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"profilePhotoUrl"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ProfileHeaderValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateProfileHeaderMutation, UpdateProfileHeaderMutationVariables>;
export const UploadProfilePhotoDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UploadProfilePhoto"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"file"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Upload"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uploadProfilePhoto"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"file"},"value":{"kind":"Variable","name":{"kind":"Name","value":"file"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"UploadProfilePhotoResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"profile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"profilePhotoUrl"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"FileValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UploadProfilePhotoMutation, UploadProfilePhotoMutationVariables>;
export const DeleteProfilePhotoDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteProfilePhoto"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteProfilePhoto"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"DeleteProfilePhotoResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"success"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ProfileHeaderValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<DeleteProfilePhotoMutation, DeleteProfilePhotoMutationVariables>;
export const DeleteTestimonialDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteTestimonial"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteTestimonial"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteTestimonialMutation, DeleteTestimonialMutationVariables>;
export const UpdateAuthorDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateAuthor"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateAuthorInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateAuthor"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"linkedInUrl"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<UpdateAuthorMutation, UpdateAuthorMutationVariables>;
export const UploadAuthorImageDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UploadAuthorImage"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"authorId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"file"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Upload"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uploadAuthorImage"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"authorId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"authorId"}}},{"kind":"Argument","name":{"kind":"Name","value":"file"},"value":{"kind":"Variable","name":{"kind":"Name","value":"file"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"UploadAuthorImageResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}},{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"linkedInUrl"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"FileValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UploadAuthorImageMutation, UploadAuthorImageMutationVariables>;
export const UploadForDetectionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UploadForDetection"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"file"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Upload"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"uploadForDetection"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"file"},"value":{"kind":"Variable","name":{"kind":"Name","value":"file"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"UploadForDetectionResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"fileId"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"FileValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UploadForDetectionMutation, UploadForDetectionMutationVariables>;
export const ReportDocumentFeedbackDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ReportDocumentFeedback"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"DocumentFeedbackInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"reportDocumentFeedback"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}}]}}]}}]} as unknown as DocumentNode<ReportDocumentFeedbackMutation, ReportDocumentFeedbackMutationVariables>;
export const ProcessDocumentDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ProcessDocument"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ProcessDocumentInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"processDocument"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ProcessDocumentResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"resumeId"}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetterID"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ProcessDocumentError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<ProcessDocumentMutation, ProcessDocumentMutationVariables>;
export const ImportDocumentResultsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ImportDocumentResults"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ImportDocumentResultsInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"importDocumentResults"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ImportDocumentResultsResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"profile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"importedCount"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experiences"}},{"kind":"Field","name":{"kind":"Name","value":"educations"}},{"kind":"Field","name":{"kind":"Name","value":"skills"}},{"kind":"Field","name":{"kind":"Name","value":"testimonials"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ImportDocumentResultsError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<ImportDocumentResultsMutation, ImportDocumentResultsMutationVariables>;
export const GetReferenceLetterForViewerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetterForViewer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}}]}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterForViewerQuery, GetReferenceLetterForViewerQueryVariables>;
export const GetReferenceLetterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetter"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"dateWritten"}},{"kind":"Field","name":{"kind":"Name","value":"rawText"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}}]}},{"kind":"Field","name":{"kind":"Name","value":"testimonials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"skillsMentioned"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skillMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"context"}}]}},{"kind":"Field","name":{"kind":"Name","value":"experienceMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}}]}},{"kind":"Field","name":{"kind":"Name","value":"discoveredSkills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"context"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"modelVersion"}},{"kind":"Field","name":{"kind":"Name","value":"processingTimeMs"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterQuery, GetReferenceLetterQueryVariables>;
export const GetReferenceLettersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetters"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLettersQuery, GetReferenceLettersQueryVariables>;
export const GetReferenceLetterStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetterStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterStatusQuery, GetReferenceLetterStatusQueryVariables>;
export const GetUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;
export const GetFilesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetFiles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"files"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}},{"kind":"Field","name":{"kind":"Name","value":"storageKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetFilesQuery, GetFilesQueryVariables>;
export const GetResumeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetResume"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resume"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"errorMessage"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}}]}}]}}]}}]} as unknown as DocumentNode<GetResumeQuery, GetResumeQueryVariables>;
export const GetUserResumesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUserResumes"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resumes"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<GetUserResumesQuery, GetUserResumesQueryVariables>;
export const GetProfileDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetProfile"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"profileByUserId"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"profilePhotoUrl"}},{"kind":"Field","name":{"kind":"Name","value":"experiences"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"validationCount"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"educations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"validationCount"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetProfileQuery, GetProfileQueryVariables>;
export const GetProfileByIdDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetProfileById"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"profile"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"profilePhotoUrl"}},{"kind":"Field","name":{"kind":"Name","value":"experiences"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"validationCount"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"educations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"validationCount"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetProfileByIdQuery, GetProfileByIdQueryVariables>;
export const GetTestimonialsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetTestimonials"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"profileId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"testimonials"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"profileId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"profileId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"linkedInUrl"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}}]}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"authorCompany"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"validatedSkills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetTestimonialsQuery, GetTestimonialsQueryVariables>;
export const GetSkillValidationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetSkillValidations"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"skillId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skillValidations"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"skillId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"skillId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"quoteSnippet"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetSkillValidationsQuery, GetSkillValidationsQueryVariables>;
export const GetExperienceValidationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetExperienceValidations"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"experienceId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experienceValidations"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"experienceId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"experienceId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"quoteSnippet"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetExperienceValidationsQuery, GetExperienceValidationsQueryVariables>;
export const GetDocumentDetectionStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetDocumentDetectionStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"fileId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"documentDetectionStatus"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"fileId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"fileId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"fileId"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"detection"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hasCareerInfo"}},{"kind":"Field","name":{"kind":"Name","value":"hasTestimonial"}},{"kind":"Field","name":{"kind":"Name","value":"testimonialAuthor"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"documentTypeHint"}},{"kind":"Field","name":{"kind":"Name","value":"fileId"}}]}},{"kind":"Field","name":{"kind":"Name","value":"error"}}]}}]}}]} as unknown as DocumentNode<GetDocumentDetectionStatusQuery, GetDocumentDetectionStatusQueryVariables>;
export const GetDocumentProcessingStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetDocumentProcessingStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"resumeId"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"referenceLetterID"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"documentProcessingStatus"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"resumeId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"resumeId"}}},{"kind":"Argument","name":{"kind":"Name","value":"referenceLetterID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"referenceLetterID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"allComplete"}},{"kind":"Field","name":{"kind":"Name","value":"resume"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"errorMessage"}}]}},{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}}]}},{"kind":"Field","name":{"kind":"Name","value":"testimonials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"skillsMentioned"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skillMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"context"}}]}},{"kind":"Field","name":{"kind":"Name","value":"experienceMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}}]}},{"kind":"Field","name":{"kind":"Name","value":"discoveredSkills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"context"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"modelVersion"}},{"kind":"Field","name":{"kind":"Name","value":"processingTimeMs"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetDocumentProcessingStatusQuery, GetDocumentProcessingStatusQueryVariables>;