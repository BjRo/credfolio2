interface SkillTagProps {
  skill: string;
}

function SkillTag({ skill }: SkillTagProps) {
  return (
    <span className="inline-flex items-center px-3 py-1.5 bg-primary/10 text-primary rounded-full text-sm font-medium border border-primary/20 hover:bg-primary/20 transition-colors">
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
    <div className="bg-card shadow rounded-lg p-6 sm:p-8">
      <h2 className="text-xl font-bold text-foreground mb-4">Skills</h2>
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
