interface SkillTagProps {
  skill: string;
}

function SkillTag({ skill }: SkillTagProps) {
  return (
    <span className="inline-flex items-center px-3 py-1.5 bg-blue-50 text-blue-700 rounded-full text-sm font-medium border border-blue-100 hover:bg-blue-100 transition-colors">
      {skill}
    </span>
  );
}

interface SkillsSectionProps {
  skills: string[];
}

export function SkillsSection({ skills }: SkillsSectionProps) {
  if (skills.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Skills</h2>
      <ul className="flex flex-wrap gap-2 list-none" aria-label="Skills list">
        {skills.map((skill) => (
          <li key={skill}>
            <SkillTag skill={skill} />
          </li>
        ))}
      </ul>
    </div>
  );
}
