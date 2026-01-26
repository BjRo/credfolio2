"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";

interface WorkExperience {
  company: string;
  title: string;
  location?: string;
  startDate?: string;
  endDate?: string;
  isCurrent: boolean;
  description?: string;
}

interface Education {
  institution: string;
  degree?: string;
  field?: string;
  startDate?: string;
  endDate?: string;
  gpa?: string;
  achievements?: string;
}

interface ResumeExtractedData {
  name: string;
  email?: string;
  phone?: string;
  location?: string;
  summary?: string;
  experience: WorkExperience[];
  education: Education[];
  skills: string[];
  extractedAt: string;
  confidence: number;
}

interface Resume {
  id: string;
  status: string;
  extractedData?: ResumeExtractedData;
  createdAt: string;
  updatedAt: string;
}

type LoadingState = "loading" | "success" | "error" | "not-found";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const resumeId = params.id as string;
  const [state, setState] = useState<LoadingState>("loading");
  const [resume, setResume] = useState<Resume | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchResume = async () => {
      try {
        const response = await fetch(GRAPHQL_ENDPOINT, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: `
              query GetResume($id: ID!) {
                resume(id: $id) {
                  id
                  status
                  extractedData {
                    name
                    email
                    phone
                    location
                    summary
                    experience {
                      company
                      title
                      location
                      startDate
                      endDate
                      isCurrent
                      description
                    }
                    education {
                      institution
                      degree
                      field
                      startDate
                      endDate
                      gpa
                      achievements
                    }
                    skills
                    extractedAt
                    confidence
                  }
                  createdAt
                  updatedAt
                }
              }
            `,
            variables: { id: resumeId },
          }),
        });

        const result = await response.json();

        if (result.errors?.length) {
          setError(result.errors[0].message);
          setState("error");
          return;
        }

        const resumeData = result.data?.resume;
        if (!resumeData) {
          setState("not-found");
          return;
        }

        setResume(resumeData);
        setState("success");
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load profile");
        setState("error");
      }
    };

    fetchResume();
  }, [resumeId]);

  if (state === "loading") {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <div className="animate-pulse">
            <div className="bg-white shadow rounded-lg p-8">
              <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-8"></div>
              <div className="space-y-4">
                <div className="h-4 bg-gray-200 rounded"></div>
                <div className="h-4 bg-gray-200 rounded"></div>
                <div className="h-4 bg-gray-200 rounded w-3/4"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (state === "not-found") {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Profile Not Found</h1>
          <p className="text-gray-600 mb-6">
            The resume you&apos;re looking for doesn&apos;t exist or has been removed.
          </p>
          <button
            type="button"
            onClick={() => router.push("/upload-resume")}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Upload a Resume
          </button>
        </div>
      </div>
    );
  }

  if (state === "error") {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error Loading Profile</h1>
          <p className="text-gray-600 mb-6">{error}</p>
          <button
            type="button"
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  const data = resume?.extractedData;

  if (!data) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Profile Processing</h1>
          <p className="text-gray-600 mb-6">Your resume is still being processed. Please wait...</p>
          <div className="animate-spin h-8 w-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {/* Header / Contact Info */}
        <div className="bg-white shadow rounded-lg p-8 mb-6">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{data.name}</h1>
              <div className="mt-2 space-y-1 text-gray-600">
                {data.email && (
                  <p className="flex items-center gap-2">
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      role="img"
                      aria-label="Email"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                      />
                    </svg>
                    {data.email}
                  </p>
                )}
                {data.phone && (
                  <p className="flex items-center gap-2">
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      role="img"
                      aria-label="Phone"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                      />
                    </svg>
                    {data.phone}
                  </p>
                )}
                {data.location && (
                  <p className="flex items-center gap-2">
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      role="img"
                      aria-label="Location"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                    </svg>
                    {data.location}
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

        {/* Work Experience */}
        {data.experience.length > 0 && (
          <div className="bg-white shadow rounded-lg p-8 mb-6">
            <h2 className="text-xl font-bold text-gray-900 mb-6">Work Experience</h2>
            <div className="space-y-6">
              {data.experience.map((exp, index) => (
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
        )}

        {/* Education */}
        {data.education.length > 0 && (
          <div className="bg-white shadow rounded-lg p-8 mb-6">
            <h2 className="text-xl font-bold text-gray-900 mb-6">Education</h2>
            <div className="space-y-6">
              {data.education.map((edu, index) => (
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
        )}

        {/* Skills */}
        {data.skills.length > 0 && (
          <div className="bg-white shadow rounded-lg p-8">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Skills</h2>
            <div className="flex flex-wrap gap-2">
              {data.skills.map((skill) => (
                <span
                  key={skill}
                  className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium"
                >
                  {skill}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Back Link */}
        <div className="mt-8 text-center">
          <button
            type="button"
            onClick={() => router.push("/upload-resume")}
            className="text-blue-600 hover:text-blue-700 text-sm"
          >
            Upload another resume
          </button>
        </div>
      </div>
    </div>
  );
}
