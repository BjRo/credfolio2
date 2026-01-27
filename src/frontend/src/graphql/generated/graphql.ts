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

/** Relationship type between letter author and candidate. */
export enum AuthorRelationship {
  Client = 'CLIENT',
  Colleague = 'COLLEAGUE',
  Manager = 'MANAGER',
  Mentor = 'MENTOR',
  Other = 'OTHER',
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

/** Result of a delete operation. */
export type DeleteResult = {
  __typename?: 'DeleteResult';
  /** ID of the deleted item. */
  deletedId: Scalars['ID']['output'];
  /** Whether the deletion was successful. */
  success: Scalars['Boolean']['output'];
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
  /** Manually entered by user. */
  Manual = 'MANUAL',
  /** Extracted from an uploaded resume. */
  ResumeExtracted = 'RESUME_EXTRACTED'
}

/** Error returned when experience validation fails. */
export type ExperienceValidationError = {
  __typename?: 'ExperienceValidationError';
  /** The field that failed validation. */
  field?: Maybe<Scalars['String']['output']>;
  /** Error message describing the validation failure. */
  message: Scalars['String']['output'];
};

/** An accomplishment cited in a reference letter. */
export type ExtractedAccomplishment = {
  __typename?: 'ExtractedAccomplishment';
  /** Confidence score for this extraction (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Description of the accomplishment. */
  description: Scalars['String']['output'];
  /** The impact or result of the accomplishment. */
  impact?: Maybe<Scalars['String']['output']>;
  /** When the accomplishment occurred (e.g., '2024', 'Q3 2023'). */
  timeframe?: Maybe<Scalars['String']['output']>;
};

/** Author details extracted from a reference letter. */
export type ExtractedAuthor = {
  __typename?: 'ExtractedAuthor';
  /** Confidence score for this extraction (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** The author's full name. */
  name: Scalars['String']['output'];
  /** The organization where the author works. */
  organization?: Maybe<Scalars['String']['output']>;
  /** The relationship type between author and candidate. */
  relationship: AuthorRelationship;
  /** Additional details about the relationship (e.g., 'Direct supervisor for 3 years'). */
  relationshipDetails?: Maybe<Scalars['String']['output']>;
  /** The author's job title or position. */
  title?: Maybe<Scalars['String']['output']>;
};

/** Structured data extracted from a reference letter. */
export type ExtractedLetterData = {
  __typename?: 'ExtractedLetterData';
  /** Accomplishments cited in the letter. */
  accomplishments: Array<ExtractedAccomplishment>;
  /** Details about the letter's author. */
  author: ExtractedAuthor;
  /** Metadata about the extraction process. */
  metadata: ExtractionMetadata;
  /** Qualities and traits mentioned in the letter. */
  qualities: Array<ExtractedQuality>;
  /** Overall recommendation assessment. */
  recommendation: ExtractedRecommendation;
  /** Skills mentioned in the letter. */
  skills: Array<ExtractedSkill>;
};

/** A quality or trait mentioned in a reference letter. */
export type ExtractedQuality = {
  __typename?: 'ExtractedQuality';
  /** Confidence score for this extraction (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Supporting evidence or quotes from the letter. */
  evidence?: Maybe<Array<Scalars['String']['output']>>;
  /** The quality or trait name. */
  trait: Scalars['String']['output'];
};

/** Overall recommendation assessment from a reference letter. */
export type ExtractedRecommendation = {
  __typename?: 'ExtractedRecommendation';
  /** Confidence score for this extraction (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Key quotes that demonstrate the recommendation strength. */
  keyQuotes?: Maybe<Array<Scalars['String']['output']>>;
  /** Sentiment score from -1.0 (negative) to 1.0 (positive). */
  sentiment: Scalars['Float']['output'];
  /** The overall strength of the recommendation. */
  strength: RecommendationStrength;
  /** Brief summary of the recommendation. */
  summary?: Maybe<Scalars['String']['output']>;
};

/** A skill mentioned in a reference letter. */
export type ExtractedSkill = {
  __typename?: 'ExtractedSkill';
  /** The category of skill. */
  category: SkillCategory;
  /** Confidence score for this extraction (0.0 to 1.0). */
  confidence: Scalars['Float']['output'];
  /** Contextual quotes where the skill was mentioned. */
  context?: Maybe<Array<Scalars['String']['output']>>;
  /** Number of times this skill was mentioned in the letter. */
  mentions: Scalars['Int']['output'];
  /** The skill name as it appears in the letter. */
  name: Scalars['String']['output'];
  /** Normalized skill name for aggregation (e.g., 'JavaScript' and 'JS' both become 'javascript'). */
  normalizedName: Scalars['String']['output'];
};

/** Metadata about the extraction process. */
export type ExtractionMetadata = {
  __typename?: 'ExtractionMetadata';
  /** When the extraction was performed. */
  extractedAt: Scalars['DateTime']['output'];
  /** The LLM model version used for extraction. */
  modelVersion: Scalars['String']['output'];
  /** Overall confidence score for the entire extraction (0.0 to 1.0). */
  overallConfidence: Scalars['Float']['output'];
  /** Time taken to process the extraction in milliseconds. */
  processingTimeMs?: Maybe<Scalars['Int']['output']>;
  /** Warning messages from the extraction process. */
  warnings?: Maybe<Array<Scalars['String']['output']>>;
  /** Number of warnings generated during extraction. */
  warningsCount: Scalars['Int']['output'];
};

/** An uploaded file stored in object storage. */
export type File = {
  __typename?: 'File';
  contentType: Scalars['String']['output'];
  createdAt: Scalars['DateTime']['output'];
  filename: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  sizeBytes: Scalars['Int']['output'];
  storageKey: Scalars['String']['output'];
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

export type Mutation = {
  __typename?: 'Mutation';
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
  /** Delete an education entry. */
  deleteEducation: DeleteResult;
  /** Delete a work experience. */
  deleteExperience: DeleteResult;
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
   * Upload a reference letter file for processing.
   * Accepts PDF, DOCX, or TXT files.
   * Creates a file record and queues the document for LLM extraction.
   */
  uploadFile: UploadFileResponse;
  /**
   * Upload a resume file for processing.
   * Accepts PDF, DOCX, or TXT files.
   * Creates a file record and queues the resume for LLM extraction.
   */
  uploadResume: UploadResumeResponse;
};


export type MutationCreateEducationArgs = {
  input: CreateEducationInput;
  userId: Scalars['ID']['input'];
};


export type MutationCreateExperienceArgs = {
  input: CreateExperienceInput;
  userId: Scalars['ID']['input'];
};


export type MutationDeleteEducationArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteExperienceArgs = {
  id: Scalars['ID']['input'];
};


export type MutationUpdateEducationArgs = {
  id: Scalars['ID']['input'];
  input: UpdateEducationInput;
};


export type MutationUpdateExperienceArgs = {
  id: Scalars['ID']['input'];
  input: UpdateExperienceInput;
};


export type MutationUploadFileArgs = {
  file: Scalars['Upload']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationUploadResumeArgs = {
  file: Scalars['Upload']['input'];
  userId: Scalars['ID']['input'];
};

/** A user's profile containing manually editable data. */
export type Profile = {
  __typename?: 'Profile';
  createdAt: Scalars['DateTime']['output'];
  /** Education entries. */
  educations: Array<ProfileEducation>;
  /** Work experience entries. */
  experiences: Array<ProfileExperience>;
  id: Scalars['ID']['output'];
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
};

export type Query = {
  __typename?: 'Query';
  /** Get a file by ID. */
  file?: Maybe<File>;
  /** Get all files for a user. */
  files: Array<File>;
  /**
   * Get a user's profile by user ID.
   * Returns null if no profile exists.
   */
  profile?: Maybe<Profile>;
  /** Get a single profile education entry by ID. */
  profileEducation?: Maybe<ProfileEducation>;
  /** Get a single profile experience by ID. */
  profileExperience?: Maybe<ProfileExperience>;
  /** Get a reference letter by ID. */
  referenceLetter?: Maybe<ReferenceLetter>;
  /** Get all reference letters for a user. */
  referenceLetters: Array<ReferenceLetter>;
  /** Get a resume by ID. */
  resume?: Maybe<Resume>;
  /** Get all resumes for a user. */
  resumes: Array<Resume>;
  /** Get a user by ID. */
  user?: Maybe<User>;
};


export type QueryFileArgs = {
  id: Scalars['ID']['input'];
};


export type QueryFilesArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryProfileArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryProfileEducationArgs = {
  id: Scalars['ID']['input'];
};


export type QueryProfileExperienceArgs = {
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


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};

/** Recommendation strength level. */
export enum RecommendationStrength {
  Moderate = 'MODERATE',
  Reserved = 'RESERVED',
  Strong = 'STRONG'
}

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
  /** Email address. */
  email?: Maybe<Scalars['String']['output']>;
  /** When the extraction was performed. */
  extractedAt: Scalars['DateTime']['output'];
  /** Location (city, state, country). */
  location?: Maybe<Scalars['String']['output']>;
  /** Full name of the candidate. */
  name: Scalars['String']['output'];
  /** Phone number. */
  phone?: Maybe<Scalars['String']['output']>;
  /** Skills list. */
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

/** Skill category classification. */
export enum SkillCategory {
  Domain = 'DOMAIN',
  Soft = 'SOFT',
  Technical = 'TECHNICAL'
}

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

/** Union type for upload result - either success or validation error. */
export type UploadFileResponse = FileValidationError | UploadFileResult;

/** Result of a file upload operation. */
export type UploadFileResult = {
  __typename?: 'UploadFileResult';
  /** The uploaded file metadata. */
  file: File;
  /** The reference letter created for processing. */
  referenceLetter: ReferenceLetter;
};

/** Union type for resume upload result - either success or validation error. */
export type UploadResumeResponse = FileValidationError | UploadResumeResult;

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

export type TestConnectionQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type TestConnectionQuery = { __typename?: 'Query', referenceLetters: Array<{ __typename?: 'ReferenceLetter', id: string, title?: string | null, status: ReferenceLetterStatus, authorName?: string | null, createdAt: string }> };

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

export type GetReferenceLetterQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetReferenceLetterQuery = { __typename?: 'Query', referenceLetter?: { __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, dateWritten?: string | null, rawText?: string | null, status: ReferenceLetterStatus, createdAt: string, updatedAt: string, extractedData?: { __typename?: 'ExtractedLetterData', author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, organization?: string | null, relationship: AuthorRelationship, relationshipDetails?: string | null, confidence: number }, skills: Array<{ __typename?: 'ExtractedSkill', name: string, normalizedName: string, category: SkillCategory, mentions: number, context?: Array<string> | null, confidence: number }>, qualities: Array<{ __typename?: 'ExtractedQuality', trait: string, evidence?: Array<string> | null, confidence: number }>, accomplishments: Array<{ __typename?: 'ExtractedAccomplishment', description: string, impact?: string | null, timeframe?: string | null, confidence: number }>, recommendation: { __typename?: 'ExtractedRecommendation', strength: RecommendationStrength, sentiment: number, keyQuotes?: Array<string> | null, summary?: string | null, confidence: number }, metadata: { __typename?: 'ExtractionMetadata', extractedAt: string, modelVersion: string, overallConfidence: number, processingTimeMs?: number | null, warningsCount: number, warnings?: Array<string> | null } } | null, user: { __typename?: 'User', id: string, email: string, name?: string | null }, file?: { __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number } | null } | null };

export type GetReferenceLettersQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetReferenceLettersQuery = { __typename?: 'Query', referenceLetters: Array<{ __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, status: ReferenceLetterStatus, createdAt: string }> };

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


export type GetResumeQuery = { __typename?: 'Query', resume?: { __typename?: 'Resume', id: string, status: ResumeStatus, errorMessage?: string | null, createdAt: string, updatedAt: string, extractedData?: { __typename?: 'ResumeExtractedData', name: string, email?: string | null, phone?: string | null, location?: string | null, summary?: string | null, skills: Array<string>, extractedAt: string, confidence: number } | null, user: { __typename?: 'User', id: string, email: string, name?: string | null }, file: { __typename?: 'File', id: string, filename: string, contentType: string } } | null };

export type GetUserResumesQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetUserResumesQuery = { __typename?: 'Query', resumes: Array<{ __typename?: 'Resume', id: string, status: ResumeStatus }> };

export type GetProfileQueryVariables = Exact<{
  userId: Scalars['ID']['input'];
}>;


export type GetProfileQuery = { __typename?: 'Query', profile?: { __typename?: 'Profile', id: string, createdAt: string, updatedAt: string, experiences: Array<{ __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }>, educations: Array<{ __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }> } | null };


export const TestConnectionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"TestConnection"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<TestConnectionQuery, TestConnectionQueryVariables>;
export const CreateExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateExperienceInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experience"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<CreateExperienceMutation, CreateExperienceMutationVariables>;
export const UpdateExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateExperienceInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"experience"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"ExperienceValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateExperienceMutation, UpdateExperienceMutationVariables>;
export const DeleteExperienceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteExperience"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteExperience"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteExperienceMutation, DeleteExperienceMutationVariables>;
export const CreateEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateEducationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"education"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<CreateEducationMutation, CreateEducationMutationVariables>;
export const UpdateEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateEducationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationResult"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"education"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"EducationValidationError"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"message"}},{"kind":"Field","name":{"kind":"Name","value":"field"}}]}}]}}]}}]} as unknown as DocumentNode<UpdateEducationMutation, UpdateEducationMutationVariables>;
export const DeleteEducationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteEducation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteEducation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}},{"kind":"Field","name":{"kind":"Name","value":"deletedId"}}]}}]}}]} as unknown as DocumentNode<DeleteEducationMutation, DeleteEducationMutationVariables>;
export const GetReferenceLetterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetter"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"dateWritten"}},{"kind":"Field","name":{"kind":"Name","value":"rawText"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}},{"kind":"Field","name":{"kind":"Name","value":"relationshipDetails"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"mentions"}},{"kind":"Field","name":{"kind":"Name","value":"context"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"qualities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"trait"}},{"kind":"Field","name":{"kind":"Name","value":"evidence"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"accomplishments"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"impact"}},{"kind":"Field","name":{"kind":"Name","value":"timeframe"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"recommendation"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"strength"}},{"kind":"Field","name":{"kind":"Name","value":"sentiment"}},{"kind":"Field","name":{"kind":"Name","value":"keyQuotes"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"modelVersion"}},{"kind":"Field","name":{"kind":"Name","value":"overallConfidence"}},{"kind":"Field","name":{"kind":"Name","value":"processingTimeMs"}},{"kind":"Field","name":{"kind":"Name","value":"warningsCount"}},{"kind":"Field","name":{"kind":"Name","value":"warnings"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterQuery, GetReferenceLetterQueryVariables>;
export const GetReferenceLettersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetters"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLettersQuery, GetReferenceLettersQueryVariables>;
export const GetUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;
export const GetFilesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetFiles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"files"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}},{"kind":"Field","name":{"kind":"Name","value":"storageKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetFilesQuery, GetFilesQueryVariables>;
export const GetResumeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetResume"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resume"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"skills"}},{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"errorMessage"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}}]}}]}}]}}]} as unknown as DocumentNode<GetResumeQuery, GetResumeQueryVariables>;
export const GetUserResumesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUserResumes"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resumes"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<GetUserResumesQuery, GetUserResumesQueryVariables>;
export const GetProfileDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetProfile"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"profile"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"experiences"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"educations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetProfileQuery, GetProfileQueryVariables>;