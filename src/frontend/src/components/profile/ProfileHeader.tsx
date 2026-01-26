import { Mail, MapPin, Phone } from "lucide-react";
import type { ProfileData } from "./types";

interface ProfileHeaderProps {
  data: ProfileData;
}

export function ProfileHeader({ data }: ProfileHeaderProps) {
  return (
    <div className="bg-white shadow rounded-lg p-8">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{data.name}</h1>
          <div className="mt-2 space-y-1 text-gray-600">
            {data.email && (
              <p className="flex items-center gap-2">
                <Mail className="w-4 h-4" aria-hidden="true" />
                <span>{data.email}</span>
              </p>
            )}
            {data.phone && (
              <p className="flex items-center gap-2">
                <Phone className="w-4 h-4" aria-hidden="true" />
                <span>{data.phone}</span>
              </p>
            )}
            {data.location && (
              <p className="flex items-center gap-2">
                <MapPin className="w-4 h-4" aria-hidden="true" />
                <span>{data.location}</span>
              </p>
            )}
          </div>
        </div>
        <div className="text-right text-sm text-gray-500">
          <p>Confidence: {Math.round(data.confidence * 100)}%</p>
        </div>
      </div>

      {data.summary && (
        <div className="mt-6">
          <h2 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
            Summary
          </h2>
          <p className="text-gray-700">{data.summary}</p>
        </div>
      )}
    </div>
  );
}
