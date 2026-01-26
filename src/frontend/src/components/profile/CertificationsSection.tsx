import { Award, ExternalLink } from "lucide-react";

export interface Certification {
  name: string;
  issuer: string;
  date?: string | null;
  url?: string | null;
}

interface CertificationItemProps {
  certification: Certification;
}

function CertificationItem({ certification }: CertificationItemProps) {
  return (
    <div className="flex items-start gap-3 py-3 first:pt-0 last:pb-0 border-b border-gray-100 last:border-0">
      <div className="w-8 h-8 rounded-lg bg-amber-50 flex items-center justify-center flex-shrink-0">
        <Award className="w-4 h-4 text-amber-600" aria-hidden="true" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h3 className="font-medium text-gray-900">
              {certification.url ? (
                <a
                  href={certification.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:text-blue-600 inline-flex items-center gap-1"
                >
                  {certification.name}
                  <ExternalLink className="w-3 h-3" aria-hidden="true" />
                  <span className="sr-only">(opens in new tab)</span>
                </a>
              ) : (
                certification.name
              )}
            </h3>
            <p className="text-sm text-gray-600">{certification.issuer}</p>
          </div>
          {certification.date && (
            <span className="text-sm text-gray-500 flex-shrink-0">{certification.date}</span>
          )}
        </div>
      </div>
    </div>
  );
}

interface CertificationsSectionProps {
  certifications: Certification[];
}

export function CertificationsSection({ certifications }: CertificationsSectionProps) {
  if (certifications.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Certifications</h2>
      <div className="divide-y divide-gray-100">
        {certifications.map((cert, index) => (
          <CertificationItem key={`${cert.name}-${cert.issuer}-${index}`} certification={cert} />
        ))}
      </div>
    </div>
  );
}
