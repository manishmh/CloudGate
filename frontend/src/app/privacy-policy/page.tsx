import { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Privacy Policy - CloudGate",
  description:
    "CloudGate Privacy Policy - How we collect, use, and protect your information",
};

async function getPrivacyPolicy() {
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

    return await response.text();
  } catch (error) {
    console.error("Error fetching privacy policy:", error);
    return null;
  }
}

export default async function PrivacyPolicyPage() {
  const privacyPolicy = await getPrivacyPolicy();

  if (!privacyPolicy) {
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
