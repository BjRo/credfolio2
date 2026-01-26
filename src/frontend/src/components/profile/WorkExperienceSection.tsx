import type { WorkExperience } from "./types";

interface WorkExperienceSectionProps {
  experience: WorkExperience[];
}

export function WorkExperienceSection({ experience }: WorkExperienceSectionProps) {
  if (experience.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-6">Work Experience</h2>
      <div className="space-y-6">
        {experience.map((exp, index) => (
          <div
            key={`${exp.company}-${exp.title}-${index}`}
            className={index > 0 ? "pt-6 border-t border-gray-200" : ""}
          >
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">{exp.title}</h3>
                <p className="text-gray-700">{exp.company}</p>
                {exp.location && <p className="text-sm text-gray-500">{exp.location}</p>}
              </div>
              <div className="text-sm text-gray-500 text-right">
                {exp.startDate && (
                  <p>
                    {exp.startDate} - {exp.isCurrent ? "Present" : exp.endDate || "N/A"}
                  </p>
                )}
              </div>
            </div>
            {exp.description && (
              <p className="mt-2 text-gray-600 whitespace-pre-line">{exp.description}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
