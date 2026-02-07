"use client";

import {
  AlertCircle,
  Camera,
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  FileText,
  Linkedin,
  Loader2,
  MessageSquareQuote,
  MoreVertical,
  Pencil,
  Plus,
  Trash2,
} from "lucide-react";
import Image from "next/image";
import { useCallback, useMemo, useRef, useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { type GetTestimonialsQuery, TestimonialRelationship } from "@/graphql/generated/graphql";
import { GRAPHQL_UPLOAD_ENDPOINT } from "@/lib/urql/client";
import { buildViewerUrl } from "@/lib/viewer";
import { AuthorEditModal } from "./AuthorEditModal";
import { DeleteTestimonialDialog } from "./DeleteTestimonialDialog";

// Upload author image using XHR following GraphQL multipart request spec
async function uploadAuthorImageXhr(
  authorId: string,
  file: File
): Promise<{ success: boolean; error?: string }> {
  const operations = JSON.stringify({
    query: `
      mutation UploadAuthorImage($authorId: ID!, $file: Upload!) {
        uploadAuthorImage(authorId: $authorId, file: $file) {
          ... on UploadAuthorImageResult {
            __typename
            file {
              id
            }
            author {
              id
              imageUrl
            }
          }
          ... on FileValidationError {
            __typename
            message
            field
          }
        }
      }
    `,
    variables: {
      authorId,
      file: null,
    },
  });

  const map = JSON.stringify({
    "0": ["variables.file"],
  });

  const formData = new FormData();
  formData.append("operations", operations);
  formData.append("map", map);
  formData.append("0", file);

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();

    xhr.addEventListener("load", () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const result = JSON.parse(xhr.responseText);
          if (result.errors?.length) {
            reject(new Error(result.errors[0].message));
            return;
          }
          const data = result.data?.uploadAuthorImage;
          if (data?.__typename === "FileValidationError") {
            reject(new Error(data.message));
            return;
          }
          resolve({ success: true });
        } catch {
          reject(new Error("Failed to parse response"));
        }
      } else {
        reject(new Error(`Upload failed with status ${xhr.status}`));
      }
    });

    xhr.addEventListener("error", () => {
      reject(new Error("Network error during upload"));
    });

    xhr.open("POST", GRAPHQL_UPLOAD_ENDPOINT);
    xhr.send(formData);
  });
}

const RELATIONSHIP_LABELS: Record<TestimonialRelationship, string> = {
  [TestimonialRelationship.Manager]: "Manager",
  [TestimonialRelationship.Peer]: "Peer",
  [TestimonialRelationship.DirectReport]: "Direct Report",
  [TestimonialRelationship.Client]: "Client",
  [TestimonialRelationship.Other]: "Other",
};

// Number of quotes to show before collapsing
const COLLAPSE_THRESHOLD = 2;

type Testimonial = GetTestimonialsQuery["testimonials"][number];

// Author info extracted from testimonial (using author entity or legacy fields)
interface AuthorInfo {
  id: string;
  name: string;
  title: string | null | undefined;
  company: string | null | undefined;
  linkedInUrl: string | null | undefined;
  imageUrl: string | null | undefined;
  isLegacy: boolean;
}

// A group of testimonials from the same author
interface TestimonialGroup {
  author: AuthorInfo;
  relationship: TestimonialRelationship;
  testimonials: Testimonial[];
}

// Check if an author is unknown (needs details added)
function isUnknownAuthor(author: AuthorInfo): boolean {
  return !author.name || author.name.toLowerCase() === "unknown";
}

// Get author info from testimonial, preferring author entity over legacy fields
function getAuthorInfo(testimonial: Testimonial): AuthorInfo {
  if (testimonial.author) {
    return {
      id: testimonial.author.id,
      name: testimonial.author.name,
      title: testimonial.author.title,
      company: testimonial.author.company,
      linkedInUrl: testimonial.author.linkedInUrl,
      imageUrl: testimonial.author.imageUrl,
      isLegacy: false,
    };
  }
  // Fallback to legacy fields
  return {
    id: `legacy-${testimonial.authorName}`,
    name: testimonial.authorName,
    title: testimonial.authorTitle,
    company: testimonial.authorCompany,
    linkedInUrl: null,
    imageUrl: null,
    isLegacy: true,
  };
}

// Group testimonials by author
function groupTestimonialsByAuthor(testimonials: Testimonial[]): TestimonialGroup[] {
  const groups = new Map<string, TestimonialGroup>();

  for (const testimonial of testimonials) {
    const author = getAuthorInfo(testimonial);
    const existingGroup = groups.get(author.id);

    if (existingGroup) {
      existingGroup.testimonials.push(testimonial);
    } else {
      groups.set(author.id, {
        author,
        relationship: testimonial.relationship,
        testimonials: [testimonial],
      });
    }
  }

  return Array.from(groups.values());
}

interface QuoteItemProps {
  testimonial: Testimonial;
  sourceUrl?: string;
  onSkillClick?: (skillId: string) => void;
  onDeleteClick?: (testimonial: Testimonial) => void;
}

function QuoteItem({ testimonial, sourceUrl, onSkillClick, onDeleteClick }: QuoteItemProps) {
  const handleSkillClick = useCallback(
    (skillId: string) => {
      if (onSkillClick) {
        onSkillClick(skillId);
      } else {
        // Default behavior: scroll to the skill element
        const element = document.getElementById(`skill-${skillId}`);
        if (element) {
          element.scrollIntoView({ behavior: "smooth", block: "center" });
          // Briefly highlight the element
          element.classList.add("ring-2", "ring-primary", "ring-offset-2");
          setTimeout(() => {
            element.classList.remove("ring-2", "ring-primary", "ring-offset-2");
          }, 2000);
        }
      }
    },
    [onSkillClick]
  );

  const hasMenuActions = sourceUrl || onDeleteClick;

  return (
    <div className="group/quote pl-4 relative" data-testid="quote-item">
      {/* Kebab menu - appears on hover */}
      {hasMenuActions && (
        <div className="absolute -right-2 top-0 opacity-0 group-hover/quote:opacity-100 transition-opacity">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="p-1 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
                aria-label="More actions"
              >
                <MoreVertical className="h-4 w-4" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {sourceUrl && (
                <DropdownMenuItem asChild>
                  <a href={sourceUrl} target="_blank" rel="noopener noreferrer">
                    <FileText className="h-4 w-4" />
                    View source document
                  </a>
                </DropdownMenuItem>
              )}
              {sourceUrl && onDeleteClick && <DropdownMenuSeparator />}
              {onDeleteClick && (
                <DropdownMenuItem
                  onClick={() => onDeleteClick(testimonial)}
                  className="text-destructive focus:text-destructive"
                >
                  <Trash2 className="h-4 w-4" />
                  Delete
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      )}

      {/* Quote */}
      <blockquote>
        <p className="text-foreground italic leading-relaxed text-sm">
          <span
            className="text-2xl text-primary/20 font-serif leading-none align-text-top mr-0.5"
            data-testid="opening-quote-mark"
          >
            &ldquo;
          </span>
          {testimonial.quote}
          <span className="text-2xl text-primary/20 font-serif leading-none align-bottom">
            &rdquo;
          </span>
        </p>
      </blockquote>

      {/* Validated Skills */}
      {testimonial.validatedSkills && testimonial.validatedSkills.length > 0 && (
        <div className="mt-2 flex items-center gap-2 flex-wrap">
          <span className="text-xs text-muted-foreground flex items-center gap-1">
            <CheckCircle2 className="h-3 w-3" />
            Validates:
          </span>
          {testimonial.validatedSkills.map((skill, index) => (
            <span key={skill.id} className="inline-flex items-center">
              <button
                type="button"
                onClick={() => handleSkillClick(skill.id)}
                className="text-xs text-primary hover:text-primary/80 hover:underline transition-colors"
              >
                {skill.name}
              </button>
              {index < testimonial.validatedSkills.length - 1 && (
                <span className="text-muted-foreground mx-1">·</span>
              )}
            </span>
          ))}
        </div>
      )}
    </div>
  );
}

interface TestimonialGroupCardProps {
  group: TestimonialGroup;
  onSkillClick?: (skillId: string) => void;
  onDeleteClick?: (testimonial: Testimonial) => void;
  onEditAuthor?: (author: AuthorInfo) => void;
  onAvatarUploadSuccess?: () => void;
}

function TestimonialGroupCard({
  group,
  onSkillClick,
  onDeleteClick,
  onEditAuthor,
  onAvatarUploadSuccess,
}: TestimonialGroupCardProps) {
  const { author, relationship, testimonials } = group;
  const [isExpanded, setIsExpanded] = useState(testimonials.length <= COLLAPSE_THRESHOLD);
  const [isAvatarHovered, setIsAvatarHovered] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const avatarFileInputRef = useRef<HTMLInputElement>(null);

  const visibleTestimonials = isExpanded ? testimonials : testimonials.slice(0, COLLAPSE_THRESHOLD);
  const hiddenCount = testimonials.length - COLLAPSE_THRESHOLD;

  // Get source URL for each testimonial (for the per-quote kebab menu)
  const getSourceUrl = (testimonial: Testimonial) => {
    const letterId = testimonial.referenceLetter?.id;
    const hasFile = !!testimonial.referenceLetter?.file?.url;
    if (!letterId || !hasFile) return undefined;
    return buildViewerUrl(letterId, testimonial.quote);
  };

  const unknown = isUnknownAuthor(author);
  const canEditAuthor = onEditAuthor && !author.isLegacy;
  const canUploadAvatar = canEditAuthor && onAvatarUploadSuccess;

  // Get initials for avatar fallback
  const getInitials = (name: string): string => {
    if (!name || name.toLowerCase() === "unknown") return "?";
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const handleAvatarFileChange = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file || author.isLegacy) return;

      // Reset input so the same file can be selected again
      e.target.value = "";

      setIsUploading(true);
      try {
        await uploadAuthorImageXhr(author.id, file);
        onAvatarUploadSuccess?.();
      } catch (err) {
        console.error("Author image upload failed:", err);
      } finally {
        setIsUploading(false);
      }
    },
    [author.id, author.isLegacy, onAvatarUploadSuccess]
  );

  const handleAvatarClick = useCallback(() => {
    if (canUploadAvatar && !isUploading) {
      avatarFileInputRef.current?.click();
    }
  }, [canUploadAvatar, isUploading]);

  return (
    <div
      className={`group/card relative rounded-lg p-6 ${
        unknown
          ? "bg-muted/20 border-2 border-dashed border-warning/50"
          : "bg-muted/30 border border-border/50"
      }`}
    >
      {/* Unknown author banner */}
      {unknown && (
        <div className="mb-4 flex items-center gap-2 text-sm text-warning bg-warning/10 px-3 py-2 rounded-md">
          <AlertCircle className="h-4 w-4 flex-shrink-0" />
          <span>Author not detected</span>
          {canEditAuthor && (
            <button
              type="button"
              onClick={() => onEditAuthor(author)}
              className="ml-auto text-primary hover:text-primary/80 font-medium"
            >
              Add details
            </button>
          )}
        </div>
      )}

      {/* Author header */}
      <div className="flex items-start gap-3 mb-4">
        {/* Avatar */}
        <div className="relative flex-shrink-0">
          {canUploadAvatar ? (
            <button
              type="button"
              className={`w-10 h-10 rounded-full overflow-hidden flex items-center justify-center border-0 ${
                author.imageUrl ? "bg-muted" : unknown ? "bg-warning/20" : "bg-primary/10"
              } cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2`}
              aria-label={`Change photo for ${author.name}`}
              onClick={handleAvatarClick}
              onMouseEnter={() => setIsAvatarHovered(true)}
              onMouseLeave={() => setIsAvatarHovered(false)}
              onFocus={() => setIsAvatarHovered(true)}
              onBlur={() => setIsAvatarHovered(false)}
              disabled={isUploading}
            >
              {isUploading ? (
                <Loader2 className="w-5 h-5 text-primary animate-spin" aria-hidden="true" />
              ) : author.imageUrl ? (
                <Image
                  src={author.imageUrl}
                  alt={`Photo of ${author.name}`}
                  width={40}
                  height={40}
                  className="w-full h-full object-cover"
                  unoptimized
                />
              ) : (
                <span
                  className={`text-sm font-medium ${unknown ? "text-warning" : "text-primary"}`}
                >
                  {unknown ? "?" : getInitials(author.name)}
                </span>
              )}
            </button>
          ) : (
            <div
              className={`w-10 h-10 rounded-full overflow-hidden flex items-center justify-center ${
                author.imageUrl ? "bg-muted" : unknown ? "bg-warning/20" : "bg-primary/10"
              }`}
            >
              {author.imageUrl ? (
                <Image
                  src={author.imageUrl}
                  alt={`Photo of ${author.name}`}
                  width={40}
                  height={40}
                  className="w-full h-full object-cover"
                  unoptimized
                />
              ) : (
                <span
                  className={`text-sm font-medium ${unknown ? "text-warning" : "text-primary"}`}
                >
                  {unknown ? "?" : getInitials(author.name)}
                </span>
              )}
            </div>
          )}

          {/* Hover Overlay */}
          {canUploadAvatar && isAvatarHovered && !isUploading && (
            <div
              className="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center pointer-events-none"
              aria-hidden="true"
            >
              <Camera className="w-4 h-4 text-white" />
            </div>
          )}

          {/* Hidden File Input */}
          {canUploadAvatar && (
            <input
              ref={avatarFileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              className="hidden"
              onChange={handleAvatarFileChange}
              aria-label={`Upload author photo for ${author.name}`}
            />
          )}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p
              className={`font-semibold ${unknown ? "text-muted-foreground italic" : "text-foreground"}`}
            >
              {unknown ? "Unknown Author" : author.name}
            </p>
            {author.linkedInUrl && (
              <a
                href={author.linkedInUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-center w-4 h-4 bg-muted-foreground/80 hover:bg-foreground rounded-[3px] transition-colors"
                aria-label="LinkedIn profile"
              >
                <Linkedin className="h-3 w-3 text-background" strokeWidth={0} fill="currentColor" />
              </a>
            )}
          </div>
          {(author.title || author.company) && (
            <p className="text-sm text-muted-foreground">
              {author.title}
              {author.title && author.company && " at "}
              {author.company}
            </p>
          )}
          <span className="inline-flex items-center px-2 py-0.5 text-xs font-medium bg-secondary text-secondary-foreground rounded-full mt-1">
            {RELATIONSHIP_LABELS[relationship]}
          </span>
        </div>

        {/* Kebab menu for author */}
        {canEditAuthor && (
          <div className="opacity-0 group-hover/card:opacity-100 transition-opacity">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button
                  type="button"
                  className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
                  aria-label="Edit author"
                >
                  <MoreVertical className="h-4 w-4" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => onEditAuthor(author)}>
                  <Pencil className="h-4 w-4" />
                  Edit author
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        )}
      </div>

      {/* Quotes */}
      <div className="border-l-2 border-primary/20 space-y-4" data-testid="quotes-wrapper">
        {visibleTestimonials.map((testimonial) => (
          <QuoteItem
            key={testimonial.id}
            testimonial={testimonial}
            sourceUrl={getSourceUrl(testimonial)}
            onSkillClick={onSkillClick}
            onDeleteClick={onDeleteClick}
          />
        ))}
      </div>

      {/* Expand/collapse button */}
      {testimonials.length > COLLAPSE_THRESHOLD && (
        <button
          type="button"
          onClick={() => setIsExpanded(!isExpanded)}
          className="mt-4 flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          {isExpanded ? (
            <>
              <ChevronUp className="h-4 w-4" />
              Show less
            </>
          ) : (
            <>
              <ChevronDown className="h-4 w-4" />
              Show {hiddenCount} more
            </>
          )}
        </button>
      )}
    </div>
  );
}

interface TestimonialsSectionProps {
  testimonials: Testimonial[];
  isLoading?: boolean;
  onAddReference?: () => void;
  onSkillClick?: (skillId: string) => void;
  onTestimonialDeleted?: () => void;
  onAuthorUpdated?: () => void;
}

export function TestimonialsSection({
  testimonials,
  isLoading = false,
  onAddReference,
  onSkillClick,
  onTestimonialDeleted,
  onAuthorUpdated,
}: TestimonialsSectionProps) {
  // State for delete dialog
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [testimonialToDelete, setTestimonialToDelete] = useState<Testimonial | null>(null);

  // State for edit author modal
  const [editAuthorModalOpen, setEditAuthorModalOpen] = useState(false);
  const [authorToEdit, setAuthorToEdit] = useState<AuthorInfo | null>(null);

  // Group testimonials by author (must be called before any early returns)
  const groups = useMemo(() => groupTestimonialsByAuthor(testimonials), [testimonials]);

  const handleDeleteClick = useCallback((testimonial: Testimonial) => {
    setTestimonialToDelete(testimonial);
    setDeleteDialogOpen(true);
  }, []);

  const handleDeleteSuccess = useCallback(() => {
    setTestimonialToDelete(null);
    onTestimonialDeleted?.();
  }, [onTestimonialDeleted]);

  const handleEditAuthor = useCallback((author: AuthorInfo) => {
    setAuthorToEdit(author);
    setEditAuthorModalOpen(true);
  }, []);

  const handleEditAuthorSuccess = useCallback(() => {
    setAuthorToEdit(null);
    onAuthorUpdated?.();
  }, [onAuthorUpdated]);

  // Check if editing is enabled (we have callbacks for mutations)
  const isEditable = !!onTestimonialDeleted || !!onAuthorUpdated;

  // Don't render if no testimonials and no way to add one
  if (testimonials.length === 0 && !onAddReference) {
    return null;
  }

  return (
    <>
      <div id="testimonials" className="bg-card border rounded-lg p-6 sm:p-8">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <MessageSquareQuote className="h-5 w-5 text-primary" />
            <h2 className="text-xl font-bold text-foreground">What Others Say</h2>
          </div>
          {onAddReference && (
            <button
              type="button"
              onClick={onAddReference}
              className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
              aria-label="Add reference letter"
            >
              <Plus className="h-5 w-5" />
            </button>
          )}
        </div>

        {isLoading ? (
          <div className="space-y-4">
            {[1, 2].map((i) => (
              <div
                key={i}
                className="bg-muted/30 rounded-lg p-6 border border-border/50 animate-pulse"
              >
                <div className="h-20 bg-muted rounded" />
                <div className="mt-6 pt-4 border-t border-border/50">
                  <div className="flex items-start gap-3">
                    <div className="w-10 h-10 rounded-full bg-muted" />
                    <div className="flex-1 space-y-2">
                      <div className="h-4 bg-muted rounded w-32" />
                      <div className="h-3 bg-muted rounded w-48" />
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : testimonials.length > 0 ? (
          <div className="space-y-4">
            {groups.map((group) => (
              <TestimonialGroupCard
                key={group.author.id}
                group={group}
                onSkillClick={onSkillClick}
                onDeleteClick={isEditable ? handleDeleteClick : undefined}
                onEditAuthor={isEditable ? handleEditAuthor : undefined}
                onAvatarUploadSuccess={isEditable ? onAuthorUpdated : undefined}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <MessageSquareQuote className="h-12 w-12 text-muted-foreground/50 mx-auto mb-4" />
            <p className="text-muted-foreground mb-4">No testimonials yet.</p>
            <p className="text-sm text-muted-foreground mb-6">
              Add a reference letter to include testimonials from people who&apos;ve worked with
              you.
            </p>
            {onAddReference && (
              <button
                type="button"
                onClick={onAddReference}
                className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
              >
                <Plus className="h-4 w-4" />
                Add Reference Letter
              </button>
            )}
          </div>
        )}
      </div>

      {/* Delete confirmation dialog */}
      {testimonialToDelete && (
        <DeleteTestimonialDialog
          open={deleteDialogOpen}
          onOpenChange={setDeleteDialogOpen}
          testimonialId={testimonialToDelete.id}
          quote={testimonialToDelete.quote}
          authorName={getAuthorInfo(testimonialToDelete).name}
          onSuccess={handleDeleteSuccess}
        />
      )}

      {/* Edit author modal */}
      {authorToEdit && (
        <AuthorEditModal
          open={editAuthorModalOpen}
          onOpenChange={setEditAuthorModalOpen}
          author={{
            id: authorToEdit.id,
            name: authorToEdit.name,
            title: authorToEdit.title,
            company: authorToEdit.company,
            linkedInUrl: authorToEdit.linkedInUrl,
            imageUrl: authorToEdit.imageUrl,
          }}
          onSuccess={handleEditAuthorSuccess}
        />
      )}
    </>
  );
}
