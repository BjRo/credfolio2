"use client";

import { Camera, X } from "lucide-react";
import Image from "next/image";
import { useCallback, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

export interface ProfileHeaderFormData {
  name: string;
  email: string;
  phone: string;
  location: string;
  summary: string;
  pendingImageFile?: File | null;
  imageRemoved?: boolean;
}

interface ProfileHeaderFormProps {
  initialData?: Partial<ProfileHeaderFormData>;
  photoUrl?: string | null;
  onSubmit: (data: ProfileHeaderFormData) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
}

interface FormErrors {
  name?: string;
  email?: string;
  phone?: string;
}

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const PHONE_REGEX = /^[\d\s\-().+]+$/;

export function ProfileHeaderForm({
  initialData,
  photoUrl,
  onSubmit,
  onCancel,
  isSubmitting = false,
}: ProfileHeaderFormProps) {
  const [name, setName] = useState(initialData?.name ?? "");
  const [email, setEmail] = useState(initialData?.email ?? "");
  const [phone, setPhone] = useState(initialData?.phone ?? "");
  const [location, setLocation] = useState(initialData?.location ?? "");
  const [summary, setSummary] = useState(initialData?.summary ?? "");
  const [errors, setErrors] = useState<FormErrors>({});

  // Image state
  const [previewImageUrl, setPreviewImageUrl] = useState<string | null>(photoUrl ?? null);
  const [pendingImageFile, setPendingImageFile] = useState<File | null>(null);
  const [imageRemoved, setImageRemoved] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Reset input so the same file can be selected again
    e.target.value = "";

    // Store the file for upload on submit
    setPendingImageFile(file);
    setImageRemoved(false);

    // Create a preview URL immediately
    const localPreviewUrl = URL.createObjectURL(file);
    setPreviewImageUrl(localPreviewUrl);
  }, []);

  const handleRemoveImage = useCallback(() => {
    setPreviewImageUrl(null);
    setPendingImageFile(null);
    setImageRemoved(true);
  }, []);

  const validate = (): boolean => {
    const newErrors: FormErrors = {};

    if (!name.trim()) {
      newErrors.name = "Name is required";
    }

    if (email.trim() && !EMAIL_REGEX.test(email.trim())) {
      newErrors.email = "Invalid email format";
    }

    if (phone.trim()) {
      const cleanPhone = phone.trim();
      if (cleanPhone.length < 5 || cleanPhone.length > 30 || !PHONE_REGEX.test(cleanPhone)) {
        newErrors.phone = "Invalid phone format";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    onSubmit({
      name: name.trim(),
      email: email.trim(),
      phone: phone.trim(),
      location: location.trim(),
      summary: summary.trim(),
      pendingImageFile: pendingImageFile ?? undefined,
      imageRemoved: imageRemoved || undefined,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Profile Photo */}
      <div className="flex flex-col items-center gap-3">
        <div className="relative w-20 h-20">
          <button
            type="button"
            className={`w-20 h-20 rounded-full overflow-hidden flex items-center justify-center border-2 border-dashed border-muted-foreground/30 ${
              previewImageUrl ? "bg-muted" : "bg-muted/50"
            } cursor-pointer hover:border-primary/50 transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2`}
            onClick={() => fileInputRef.current?.click()}
            aria-label="Upload profile photo"
            disabled={isSubmitting}
          >
            {previewImageUrl ? (
              <Image
                src={previewImageUrl}
                alt="Profile photo"
                width={80}
                height={80}
                className="w-full h-full object-cover"
                unoptimized
              />
            ) : (
              <div className="flex flex-col items-center gap-1 text-muted-foreground">
                <Camera className="w-6 h-6" />
                <span className="text-xs">Add photo</span>
              </div>
            )}
          </button>
          {previewImageUrl && !isSubmitting && (
            <button
              type="button"
              onClick={handleRemoveImage}
              className="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-destructive text-destructive-foreground flex items-center justify-center hover:bg-destructive/90"
              aria-label="Remove photo"
            >
              <X className="w-3 h-3" />
            </button>
          )}
        </div>
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          className="hidden"
          onChange={handleFileChange}
          aria-label="Upload profile photo"
          disabled={isSubmitting}
        />
        <p className="text-xs text-muted-foreground">Click to upload a profile photo (optional)</p>
      </div>

      {/* Name */}
      <div className="space-y-2">
        <Label htmlFor="header-name">
          Name <span className="text-destructive">*</span>
        </Label>
        <Input
          id="header-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Your full name"
          disabled={isSubmitting}
          aria-invalid={!!errors.name}
        />
        {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
      </div>

      {/* Email */}
      <div className="space-y-2">
        <Label htmlFor="header-email">Email</Label>
        <Input
          id="header-email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your.email@example.com"
          disabled={isSubmitting}
          aria-invalid={!!errors.email}
        />
        {errors.email && <p className="text-sm text-destructive">{errors.email}</p>}
      </div>

      {/* Phone */}
      <div className="space-y-2">
        <Label htmlFor="header-phone">Phone</Label>
        <Input
          id="header-phone"
          type="tel"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          placeholder="+1 (555) 123-4567"
          disabled={isSubmitting}
          aria-invalid={!!errors.phone}
        />
        {errors.phone && <p className="text-sm text-destructive">{errors.phone}</p>}
      </div>

      {/* Location */}
      <div className="space-y-2">
        <Label htmlFor="header-location">Location</Label>
        <Input
          id="header-location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          placeholder="City, State, Country"
          disabled={isSubmitting}
        />
      </div>

      {/* Summary */}
      <div className="space-y-2">
        <Label htmlFor="header-summary">Professional Summary</Label>
        <Textarea
          id="header-summary"
          value={summary}
          onChange={(e) => setSummary(e.target.value)}
          placeholder="Write a brief professional summary..."
          rows={5}
          disabled={isSubmitting}
        />
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Saving..." : "Save Changes"}
        </Button>
      </div>
    </form>
  );
}
