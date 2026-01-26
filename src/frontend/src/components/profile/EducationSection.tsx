import { GraduationCap } from "lucide-react";
import { formatDate } from "@/lib/utils";
import type { Education } from "./types";

interface EducationCardProps {
  education: Education;
  isFirst: boolean;
}

function EducationCard({ education, isFirst }: EducationCardProps) {
  const startDate = formatDate(education.startDate);
  const endDate = formatDate(education.endDate) || "Present";
  const dateRange = startDate ? `${startDate} - ${endDate}` : null;

  const degreeField = [education.degree, education.field].filter(Boolean).join(" in ");

  // Validate GPA - should be numeric, not a date or garbage
  const isValidGpa = education.gpa && /^[\d./]+$/.test(education.gpa.trim());

  return (
    <div className={!isFirst ? "pt-6 border-t border-gray-200" : ""}>
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-1 sm:gap-4">
        <div className="flex gap-3">
          <div className="hidden sm:flex w-10 h-10 rounded-lg bg-gray-100 items-center justify-center flex-shrink-0">
            <GraduationCap className="w-5 h-5 text-gray-500" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h3 className="text-lg font-semibold text-gray-900">{education.institution}</h3>
            {degreeField && <p className="text-gray-700">{degreeField}</p>}
            {isValidGpa && (
              <p className="text-sm text-gray-500">
                <span className="font-medium">GPA:</span> {education.gpa}
              </p>
            )}
          </div>
        </div>
        <div className="text-sm text-gray-500 sm:text-right flex-shrink-0">
          {dateRange && <p>{dateRange}</p>}
        </div>
      </div>
      {education.achievements && (
        <p className="mt-3 text-gray-600 sm:ml-13">{education.achievements}</p>
      )}
    </div>
  );
}

interface EducationSectionProps {
  education: Education[];
}

export function EducationSection({ education }: EducationSectionProps) {
  if (education.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-6">Education</h2>
      <div className="space-y-6">
        {education.map((edu, index) => (
          <EducationCard
            key={`${edu.institution}-${edu.degree || ""}-${index}`}
            education={edu}
            isFirst={index === 0}
          />
        ))}
      </div>
    </div>
  );
}
