"use client";

import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  formatMonthYear,
  MonthYearPicker,
  type MonthYearValue,
  parseDateString,
} from "./MonthYearPicker";

export interface EducationFormData {
  institution: string;
  degree: string;
  field: string;
  startDate: string;
  endDate: string;
  isCurrent: boolean;
  description: string;
  gpa: string;
}

interface EducationFormProps {
  initialData?: Partial<EducationFormData>;
  onSubmit: (data: EducationFormData) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
  mode: "create" | "edit";
}

interface FormErrors {
  institution?: string;
  degree?: string;
  startDate?: string;
  endDate?: string;
}

export function EducationForm({
  initialData,
  onSubmit,
  onCancel,
  isSubmitting = false,
  mode,
}: EducationFormProps) {
  const [institution, setInstitution] = useState(initialData?.institution ?? "");
  const [degree, setDegree] = useState(initialData?.degree ?? "");
  const [field, setField] = useState(initialData?.field ?? "");
  const [startDate, setStartDate] = useState<MonthYearValue>(
    parseDateString(initialData?.startDate)
  );
  const [endDate, setEndDate] = useState<MonthYearValue>(parseDateString(initialData?.endDate));
  const [isCurrent, setIsCurrent] = useState(initialData?.isCurrent ?? false);
  const [description, setDescription] = useState(initialData?.description ?? "");
  const [gpa, setGpa] = useState(initialData?.gpa ?? "");
  const [errors, setErrors] = useState<FormErrors>({});

  const validate = (): boolean => {
    const newErrors: FormErrors = {};

    if (!institution.trim()) {
      newErrors.institution = "Institution name is required";
    }

    if (!degree.trim()) {
      newErrors.degree = "Degree is required";
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

    onSubmit({
      institution: institution.trim(),
      degree: degree.trim(),
      field: field.trim(),
      startDate: formatMonthYear(startDate),
      endDate: isCurrent ? "Present" : formatMonthYear(endDate),
      isCurrent,
      description: description.trim(),
      gpa: gpa.trim(),
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
        {/* Institution */}
        <div className="space-y-2">
          <Label htmlFor="institution">
            Institution <span className="text-destructive">*</span>
          </Label>
          <Input
            id="institution"
            value={institution}
            onChange={(e) => setInstitution(e.target.value)}
            placeholder="University or school name"
            disabled={isSubmitting}
            aria-invalid={!!errors.institution}
          />
          {errors.institution && <p className="text-sm text-destructive">{errors.institution}</p>}
        </div>

        {/* Degree */}
        <div className="space-y-2">
          <Label htmlFor="degree">
            Degree <span className="text-destructive">*</span>
          </Label>
          <Input
            id="degree"
            value={degree}
            onChange={(e) => setDegree(e.target.value)}
            placeholder="e.g., Bachelor of Science"
            disabled={isSubmitting}
            aria-invalid={!!errors.degree}
          />
          {errors.degree && <p className="text-sm text-destructive">{errors.degree}</p>}
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {/* Field of Study */}
        <div className="space-y-2">
          <Label htmlFor="field">Field of Study</Label>
          <Input
            id="field"
            value={field}
            onChange={(e) => setField(e.target.value)}
            placeholder="e.g., Computer Science"
            disabled={isSubmitting}
          />
        </div>

        {/* GPA */}
        <div className="space-y-2">
          <Label htmlFor="gpa">GPA</Label>
          <Input
            id="gpa"
            value={gpa}
            onChange={(e) => setGpa(e.target.value)}
            placeholder="e.g., 3.8"
            disabled={isSubmitting}
          />
        </div>
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
              I currently study here
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
          placeholder="Describe your studies, achievements, or activities..."
          rows={3}
          disabled={isSubmitting}
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
              ? "Add Education"
              : "Save Changes"}
        </Button>
      </div>
    </form>
  );
}
