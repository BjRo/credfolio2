"use client";

import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { HighlightsEditor } from "./HighlightsEditor";
import {
  formatMonthYear,
  MonthYearPicker,
  type MonthYearValue,
  parseDateString,
} from "./MonthYearPicker";

export interface WorkExperienceFormData {
  company: string;
  title: string;
  location: string;
  startDate: string;
  endDate: string;
  isCurrent: boolean;
  description: string;
  highlights: string[];
}

interface WorkExperienceFormProps {
  initialData?: Partial<WorkExperienceFormData>;
  onSubmit: (data: WorkExperienceFormData) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
  mode: "create" | "edit";
}

interface FormErrors {
  company?: string;
  title?: string;
  startDate?: string;
  endDate?: string;
}

export function WorkExperienceForm({
  initialData,
  onSubmit,
  onCancel,
  isSubmitting = false,
  mode,
}: WorkExperienceFormProps) {
  const [company, setCompany] = useState(initialData?.company ?? "");
  const [title, setTitle] = useState(initialData?.title ?? "");
  const [location, setLocation] = useState(initialData?.location ?? "");
  const [startDate, setStartDate] = useState<MonthYearValue>(
    parseDateString(initialData?.startDate)
  );
  const [endDate, setEndDate] = useState<MonthYearValue>(parseDateString(initialData?.endDate));
  const [isCurrent, setIsCurrent] = useState(initialData?.isCurrent ?? false);
  const [description, setDescription] = useState(initialData?.description ?? "");
  const [highlights, setHighlights] = useState<string[]>(
    initialData?.highlights?.length ? initialData.highlights : [""]
  );
  const [errors, setErrors] = useState<FormErrors>({});

  const validate = (): boolean => {
    const newErrors: FormErrors = {};

    if (!company.trim()) {
      newErrors.company = "Company name is required";
    }

    if (!title.trim()) {
      newErrors.title = "Job title is required";
    }

    // Validate start date if year is provided
    if (startDate.year && !isValidDate(startDate)) {
      newErrors.startDate = "Invalid start date";
    }

    // Validate end date if not current and year is provided
    if (!isCurrent && endDate.year && !isValidDate(endDate)) {
      newErrors.endDate = "Invalid end date";
    }

    // Validate end date is after start date if both are provided
    if (startDate.year && endDate.year && !isCurrent) {
      const startNum = Number.parseInt(`${startDate.year}${startDate.month || "01"}`, 10);
      const endNum = Number.parseInt(`${endDate.year}${endDate.month || "12"}`, 10);
      if (endNum < startNum) {
        newErrors.endDate = "End date must be after start date";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const isValidDate = (date: MonthYearValue): boolean => {
    if (!date.year) return true;
    const year = Number.parseInt(date.year, 10);
    return year >= 1970 && year <= new Date().getFullYear() + 10;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    // Filter out empty highlights
    const filteredHighlights = highlights.filter((h) => h.trim() !== "");

    onSubmit({
      company: company.trim(),
      title: title.trim(),
      location: location.trim(),
      startDate: formatMonthYear(startDate),
      endDate: isCurrent ? "Present" : formatMonthYear(endDate),
      isCurrent,
      description: description.trim(),
      highlights: filteredHighlights,
    });
  };

  const handleCurrentChange = (checked: boolean) => {
    setIsCurrent(checked);
    if (checked) {
      setEndDate({ month: "", year: "" });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {/* Company */}
        <div className="space-y-2">
          <Label htmlFor="company">
            Company <span className="text-destructive">*</span>
          </Label>
          <Input
            id="company"
            value={company}
            onChange={(e) => setCompany(e.target.value)}
            placeholder="Company name"
            disabled={isSubmitting}
            aria-invalid={!!errors.company}
          />
          {errors.company && <p className="text-sm text-destructive">{errors.company}</p>}
        </div>

        {/* Title */}
        <div className="space-y-2">
          <Label htmlFor="title">
            Job Title <span className="text-destructive">*</span>
          </Label>
          <Input
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Job title"
            disabled={isSubmitting}
            aria-invalid={!!errors.title}
          />
          {errors.title && <p className="text-sm text-destructive">{errors.title}</p>}
        </div>
      </div>

      {/* Location */}
      <div className="space-y-2">
        <Label htmlFor="location">Location</Label>
        <Input
          id="location"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          placeholder="City, State, Country"
          disabled={isSubmitting}
        />
      </div>

      {/* Dates */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label>Start Date</Label>
          <MonthYearPicker value={startDate} onChange={setStartDate} disabled={isSubmitting} />
          {errors.startDate && <p className="text-sm text-destructive">{errors.startDate}</p>}
        </div>

        <div className="space-y-2">
          <Label>End Date</Label>
          <MonthYearPicker
            value={endDate}
            onChange={setEndDate}
            disabled={isSubmitting || isCurrent}
          />
          {errors.endDate && <p className="text-sm text-destructive">{errors.endDate}</p>}
          <div className="flex items-center gap-2 pt-1">
            <Checkbox
              id="isCurrent"
              checked={isCurrent}
              onCheckedChange={handleCurrentChange}
              disabled={isSubmitting}
            />
            <Label htmlFor="isCurrent" className="text-sm font-normal cursor-pointer">
              I currently work here
            </Label>
          </div>
        </div>
      </div>

      {/* Description */}
      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Describe your role and responsibilities..."
          rows={3}
          disabled={isSubmitting}
        />
      </div>

      {/* Highlights */}
      <div className="space-y-2">
        <Label>Key Achievements</Label>
        <HighlightsEditor
          highlights={highlights}
          onChange={setHighlights}
          disabled={isSubmitting}
          placeholder="Describe an achievement or key responsibility..."
        />
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting
            ? mode === "create"
              ? "Adding..."
              : "Saving..."
            : mode === "create"
              ? "Add Experience"
              : "Save Changes"}
        </Button>
      </div>
    </form>
  );
}
