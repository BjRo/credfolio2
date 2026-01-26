interface SkillsSectionProps {
  skills: string[];
}

export function SkillsSection({ skills }: SkillsSectionProps) {
  if (skills.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Skills</h2>
      <div className="flex flex-wrap gap-2">
        {skills.map((skill) => (
          <span
            key={skill}
            className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium"
          >
            {skill}
          </span>
        ))}
      </div>
    </div>
  );
}
