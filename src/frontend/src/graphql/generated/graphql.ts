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

export type Query = {
  __typename?: 'Query';
  /** Get a file by ID. */
  file?: Maybe<File>;
  /** Get all files for a user. */
  files: Array<File>;
  /** Get a reference letter by ID. */
  referenceLetter?: Maybe<ReferenceLetter>;
  /** Get all reference letters for a user. */
  referenceLetters: Array<ReferenceLetter>;
  /** Get a user by ID. */
  user?: Maybe<User>;
};


export type QueryFileArgs = {
  id: Scalars['ID']['input'];
};


export type QueryFilesArgs = {
  userId: Scalars['ID']['input'];
};


export type QueryReferenceLetterArgs = {
  id: Scalars['ID']['input'];
};


export type QueryReferenceLettersArgs = {
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

/** Skill category classification. */
export enum SkillCategory {
  Domain = 'DOMAIN',
  Soft = 'SOFT',
  Technical = 'TECHNICAL'
}

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


export const TestConnectionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"TestConnection"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<TestConnectionQuery, TestConnectionQueryVariables>;
export const GetReferenceLetterDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetter"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetter"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"dateWritten"}},{"kind":"Field","name":{"kind":"Name","value":"rawText"}},{"kind":"Field","name":{"kind":"Name","value":"extractedData"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"author"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"relationship"}},{"kind":"Field","name":{"kind":"Name","value":"relationshipDetails"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"skills"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"normalizedName"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"mentions"}},{"kind":"Field","name":{"kind":"Name","value":"context"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"qualities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"trait"}},{"kind":"Field","name":{"kind":"Name","value":"evidence"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"accomplishments"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"impact"}},{"kind":"Field","name":{"kind":"Name","value":"timeframe"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"recommendation"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"strength"}},{"kind":"Field","name":{"kind":"Name","value":"sentiment"}},{"kind":"Field","name":{"kind":"Name","value":"keyQuotes"}},{"kind":"Field","name":{"kind":"Name","value":"summary"}},{"kind":"Field","name":{"kind":"Name","value":"confidence"}}]}},{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"extractedAt"}},{"kind":"Field","name":{"kind":"Name","value":"modelVersion"}},{"kind":"Field","name":{"kind":"Name","value":"overallConfidence"}},{"kind":"Field","name":{"kind":"Name","value":"processingTimeMs"}},{"kind":"Field","name":{"kind":"Name","value":"warningsCount"}},{"kind":"Field","name":{"kind":"Name","value":"warnings"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"file"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}}]}}]}}]}}]} as unknown as DocumentNode<GetReferenceLetterQuery, GetReferenceLetterQueryVariables>;
export const GetReferenceLettersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetReferenceLetters"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"referenceLetters"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"authorName"}},{"kind":"Field","name":{"kind":"Name","value":"authorTitle"}},{"kind":"Field","name":{"kind":"Name","value":"organization"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetReferenceLettersQuery, GetReferenceLettersQueryVariables>;
export const GetUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;
export const GetFilesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetFiles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"files"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"filename"}},{"kind":"Field","name":{"kind":"Name","value":"contentType"}},{"kind":"Field","name":{"kind":"Name","value":"sizeBytes"}},{"kind":"Field","name":{"kind":"Name","value":"storageKey"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<GetFilesQuery, GetFilesQueryVariables>;