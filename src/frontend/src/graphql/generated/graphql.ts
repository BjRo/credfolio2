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
  /** Skills discovered that aren't in the profile yet. */
  discoveredSkills: Array<Scalars['String']['output']>;
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
  /** Delete a skill. */
  deleteSkill: DeleteResult;
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
   * Update an existing skill.
   * Only updates fields that are provided.
   */
  updateSkill: SkillResponse;
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


export type MutationDeleteSkillArgs = {
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


export type MutationUpdateSkillArgs = {
  id: Scalars['ID']['input'];
  input: UpdateSkillInput;
};


export type MutationUploadFileArgs = {
  file: Scalars['Upload']['input'];
  userId: Scalars['ID']['input'];
};


export type MutationUploadResumeArgs = {
  file: Scalars['Upload']['input'];
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

/** A user's profile containing manually editable data. */
export type Profile = {
  __typename?: 'Profile';
  createdAt: Scalars['DateTime']['output'];
  /** Education entries. */
  educations: Array<ProfileEducation>;
  /** Work experience entries. */
  experiences: Array<ProfileExperience>;
  id: Scalars['ID']['output'];
  /** Skill entries. */
  skills: Array<ProfileSkill>;
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
  /** Get all testimonials for a profile. */
  testimonials: Array<Testimonial>;
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

/** Union type for skill create/update result. */
export type SkillResponse = SkillResult | SkillValidationError;

/** Result of a successful skill operation. */
export type SkillResult = {
  __typename?: 'SkillResult';
  /** The created or updated skill. */
  skill: ProfileSkill;
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
  /** Company/organization of the author. */
  authorCompany?: Maybe<Scalars['String']['output']>;
  /** Name of the person who provided the testimonial. */
  authorName: Scalars['String']['output'];
  /** Title/position of the author. */
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

/** Input for updating an existing skill. */
export type UpdateSkillInput = {
  /** Skill category. */
  category?: InputMaybe<SkillCategory>;
  /** Skill name. */
  name?: InputMaybe<Scalars['String']['input']>;
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

export type GetReferenceLetterQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetReferenceLetterQuery = { __typename?: 'Query', referenceLetter?: { __typename?: 'ReferenceLetter', id: string, title?: string | null, authorName?: string | null, authorTitle?: string | null, organization?: string | null, dateWritten?: string | null, rawText?: string | null, status: ReferenceLetterStatus, createdAt: string, updatedAt: string, extractedData?: { __typename?: 'ExtractedLetterData', discoveredSkills: Array<string>, author: { __typename?: 'ExtractedAuthor', name: string, title?: string | null, company?: string | null, relationship: AuthorRelationship }, testimonials: Array<{ __typename?: 'ExtractedTestimonial', quote: string, skillsMentioned?: Array<string> | null }>, skillMentions: Array<{ __typename?: 'ExtractedSkillMention', skill: string, quote: string, context?: string | null }>, experienceMentions: Array<{ __typename?: 'ExtractedExperienceMention', company: string, role: string, quote: string }>, metadata: { __typename?: 'ExtractionMetadata', extractedAt: string, modelVersion: string, processingTimeMs?: number | null } } | null, user: { __typename?: 'User', id: string, email: string, name?: string | null }, file?: { __typename?: 'File', id: string, filename: string, contentType: string, sizeBytes: number } | null } | null };

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


export type GetProfileQuery = { __typename?: 'Query', profile?: { __typename?: 'Profile', id: string, createdAt: string, updatedAt: string, experiences: Array<{ __typename?: 'ProfileExperience', id: string, company: string, title: string, location?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, highlights: Array<string>, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }>, educations: Array<{ __typename?: 'ProfileEducation', id: string, institution: string, degree: string, field?: string | null, startDate?: string | null, endDate?: string | null, isCurrent: boolean, description?: string | null, gpa?: string | null, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }>, skills: Array<{ __typename?: 'ProfileSkill', id: string, name: string, normalizedName: string, category: SkillCategory, displayOrder: number, source: ExperienceSource, createdAt: string, updatedAt: string }> } | null };

export type GetTestimonialsQueryVariables = Exact<{
  profileId: Scalars['ID']['input'];
}>;


export type GetTestimonialsQuery = { __typename?: 'Query', testimonials: Array<{ __typename?: 'Testimonial', id: string, quote: string, authorName: string, authorTitle?: string | null, authorCompany?: string | null, relationship: TestimonialRelationship, createdAt: string }> };


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
export const GetReferenceLetterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetter"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"dateWritten"}},{"kind":"Field","name":{"kind":"Name","value":"rawText"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}}]}},{"kind":"Field","name":{"kind":"Name","value":"testimonials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"skillsMentioned"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skillMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"skill"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"context"}}]}},{"kind":"Field","name":{"kind":"Name","value":"experienceMentions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"role"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}}]}},{"kind":"Field","name":{"kind":"Name","value":"discoveredSkills"}},{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"modelVersion"}},{"kind":"Field","name":{"kind":"Name","value":"processingTimeMs"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterQuery, GetReferenceLetterQueryVariables>;
export const GetReferenceLettersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetters"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLettersQuery, GetReferenceLettersQueryVariables>;
export const GetReferenceLetterStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetterStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterStatusQuery, GetReferenceLetterStatusQueryVariables>;
export const GetUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;
export const GetFilesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetFiles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"files"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}},{"kind":"Field","name":{"kind":"Name","value":"storageKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetFilesQuery, GetFilesQueryVariables>;
export const GetResumeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetResume"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resume"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"errorMessage"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}}]}}]}}]}}]} as unknown as DocumentNode<GetResumeQuery, GetResumeQueryVariables>;
export const GetUserResumesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUserResumes"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resumes"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<GetUserResumesQuery, GetUserResumesQueryVariables>;
export const GetProfileDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetProfile"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"profile"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"experiences"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"company"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"location"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"highlights"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"educations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"institution"}},{"kind":"Field","name":{"kind":"Name","value":"degree"}},{"kind":"Field","name":{"kind":"Name","value":"field"}},{"kind":"Field","name":{"kind":"Name","value":"startDate"}},{"kind":"Field","name":{"kind":"Name","value":"endDate"}},{"kind":"Field","name":{"kind":"Name","value":"isCurrent"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"gpa"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"displayOrder"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetProfileQuery, GetProfileQueryVariables>;
export const GetTestimonialsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetTestimonials"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"profileId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"testimonials"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"profileId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"profileId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"quote"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"authorCompany"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetTestimonialsQuery, GetTestimonialsQueryVariables>;