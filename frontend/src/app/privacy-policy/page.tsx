"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

// Note: Metadata export doesn't work with 'use client', so we'll set it via Head or document title
// export const metadata: Metadata = {
//   title: "Privacy Policy - CloudGate",
//   description:
//     "CloudGate Privacy Policy - How we collect, use, and protect your information",
// };

export default function PrivacyPolicyPage() {
  const [privacyPolicy, setPrivacyPolicy] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    // Set document title
    document.title = "Privacy Policy - CloudGate";

    async function fetchPrivacyPolicy() {
      try {
        const response = await fetch(
          `${
            process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"
          }/privacy-policy`,
          {
            cache: "force-cache", // Cache the policy content
          }
        );

        if (!response.ok) {
          throw new Error("Failed to fetch privacy policy");
        }

        const content = await response.text();
        setPrivacyPolicy(content);
      } catch (error) {
        console.error("Error fetching privacy policy:", error);
        setError(true);
      } finally {
        setLoading(false);
      }
    }

    fetchPrivacyPolicy();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="bg-white shadow-sm rounded-lg p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-6">
              Privacy Policy
            </h1>
            <div className="animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-5/6 mb-4"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !privacyPolicy) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="bg-white shadow-sm rounded-lg p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-6">
              Privacy Policy
            </h1>
            <div className="text-red-600">
              <p>Unable to load privacy policy. Please try again later.</p>
              <p className="mt-2">
                For privacy-related questions, contact us at:{" "}
                <a href="mailto:manishmh982@gmail.com" className="underline">
                  manishmh982@gmail.com
                </a>
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow-sm rounded-lg p-8">
          {/* Navigation */}
          <div className="mb-6">
            <Link
              href="/"
              className="text-blue-600 hover:text-blue-800 text-sm font-medium"
            >
              ← Back to Home
            </Link>
          </div>

          {/* Content */}
          <div className="prose prose-lg max-w-none">
            <div
              className="text-gray-800 leading-relaxed"
              dangerouslySetInnerHTML={{ __html: privacyPolicy }}
            />
          </div>

          {/* Footer */}
          <div className="mt-12 pt-8 border-t border-gray-200">
            <div className="text-sm text-gray-600">
              <p>
                <strong>Contact Information:</strong>
              </p>
              <p>
                Data Controller: Manish Kumar Saw
                <br />
                Email:{" "}
                <a
                  href="mailto:manishmh982@gmail.com"
                  className="text-blue-600 hover:text-blue-800"
                >
                  manishmh982@gmail.com
                </a>
              </p>
              <p className="mt-4">
                For privacy-related questions or data requests, please contact
                us at the above email address.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
