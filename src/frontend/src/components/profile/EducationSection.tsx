import type { Education } from "./types";

interface EducationSectionProps {
  education: Education[];
}

export function EducationSection({ education }: EducationSectionProps) {
  if (education.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-6">Education</h2>
      <div className="space-y-6">
        {education.map((edu, index) => (
          <div
            key={`${edu.institution}-${edu.degree || ""}-${index}`}
            className={index > 0 ? "pt-6 border-t border-gray-200" : ""}
          >
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">{edu.institution}</h3>
                {(edu.degree || edu.field) && (
                  <p className="text-gray-700">
                    {edu.degree}
                    {edu.degree && edu.field && " in "}
                    {edu.field}
                  </p>
                )}
                {edu.gpa && <p className="text-sm text-gray-500">GPA: {edu.gpa}</p>}
              </div>
              <div className="text-sm text-gray-500 text-right">
                {edu.startDate && (
                  <p>
                    {edu.startDate} - {edu.endDate || "Present"}
                  </p>
                )}
              </div>
            </div>
            {edu.achievements && <p className="mt-2 text-gray-600">{edu.achievements}</p>}
          </div>
        ))}
      </div>
    </div>
  );
}
