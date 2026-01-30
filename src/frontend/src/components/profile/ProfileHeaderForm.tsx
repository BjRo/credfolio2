"use client";

import { useState } from "react";

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
}

interface ProfileHeaderFormProps {
  initialData?: Partial<ProfileHeaderFormData>;
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
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
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
