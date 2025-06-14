"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

// Note: Metadata export doesn't work with 'use client', so we'll set it via Head or document title
// export const metadata: Metadata = {
//   title: "Terms of Service - CloudGate",
//   description:
//     "CloudGate Terms of Service - Terms and conditions for using our SSO platform",
// };

export default function TermsPage() {
  const [terms, setTerms] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    // Set document title
    document.title = "Terms of Service - CloudGate";

    async function fetchTerms() {
      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"}/terms`,
          {
            cache: "force-cache", // Cache the terms content
          }
        );

        if (!response.ok) {
          throw new Error("Failed to fetch terms");
        }

        const content = await response.text();
        setTerms(content);
      } catch (error) {
        console.error("Error fetching terms:", error);
        setError(true);
      } finally {
        setLoading(false);
      }
    }

    fetchTerms();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="bg-white shadow-sm rounded-lg p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-6">
              Terms of Service
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

  if (error || !terms) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="bg-white shadow-sm rounded-lg p-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-6">
              Terms of Service
            </h1>
            <div className="text-red-600">
              <p>Unable to load terms of service. Please try again later.</p>
              <p className="mt-2">
                For questions about our terms, contact us at:{" "}
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
              dangerouslySetInnerHTML={{ __html: terms }}
            />
          </div>

          {/* Footer */}
          <div className="mt-12 pt-8 border-t border-gray-200">
            <div className="text-sm text-gray-600">
              <p>
                <strong>Contact Information:</strong>
              </p>
              <p>
                Service Provider: Manish Kumar Saw
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
                For questions about these Terms of Service, please contact us at
                the above email address.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
