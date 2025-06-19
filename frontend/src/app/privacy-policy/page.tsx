"use client";

import Link from "next/link";
import { useEffect } from "react";

// Note: Metadata export doesn't work with 'use client', so we'll set it via Head or document title
// export const metadata: Metadata = {
//   title: "Privacy Policy - CloudGate",
//   description:
//     "CloudGate Privacy Policy - How we collect, use, and protect your information",
// };

export default function PrivacyPolicyPage() {
  useEffect(() => {
    document.title = "Privacy Policy - CloudGate";
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow-sm rounded-lg p-8">
          {/* Navigation */}
          <div className="mb-6">
            <Link
              href="/"
              className="text-blue-600 hover:text-blue-800 text-sm font-medium cursor-pointer"
            >
              ← Back to Home
            </Link>
          </div>

          {/* Content */}
          <div className="prose prose-lg max-w-none">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              CloudGate Privacy Policy
            </h1>
            <div className="text-sm text-gray-600 mb-8">
              <p>
                <strong>Effective Date:</strong> January 2025
              </p>
              <p>
                <strong>Last Updated:</strong> January 2025
              </p>
            </div>

            <div className="text-gray-800 leading-relaxed space-y-6">
              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  1. Introduction
                </h2>
                <p>
                  CloudGate (&quot;we,&quot; &quot;our,&quot; or &quot;us&quot;)
                  is committed to protecting your privacy. This Privacy Policy
                  explains how we collect, use, disclose, and safeguard your
                  information when you use our Single Sign-On (SSO) service.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  2. Information We Collect
                </h2>

                <h3 className="text-lg font-medium text-gray-800 mb-2">
                  2.1 Personal Information
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>
                    <strong>Account Information:</strong> Name, email address,
                    username
                  </li>
                  <li>
                    <strong>Profile Information:</strong> Profile picture,
                    preferences, settings
                  </li>
                  <li>
                    <strong>Authentication Data:</strong> Login credentials,
                    session tokens
                  </li>
                  <li>
                    <strong>Contact Information:</strong> Email for verification
                    and communication
                  </li>
                </ul>

                <h3 className="text-lg font-medium text-gray-800 mb-2 mt-4">
                  2.2 Technical Information
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>
                    <strong>Log Data:</strong> IP addresses, browser type,
                    device information
                  </li>
                  <li>
                    <strong>Usage Data:</strong> Login times, accessed
                    applications, session duration
                  </li>
                  <li>
                    <strong>Security Data:</strong> Failed login attempts,
                    security events
                  </li>
                </ul>

                <h3 className="text-lg font-medium text-gray-800 mb-2 mt-4">
                  2.3 Third-Party Integration Data
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>
                    <strong>OAuth Tokens:</strong> Access tokens for connected
                    SaaS applications
                  </li>
                  <li>
                    <strong>Application Data:</strong> Connection status, usage
                    patterns
                  </li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  3. How We Use Your Information
                </h2>
                <p>We use your information to:</p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Provide and maintain our SSO service</li>
                  <li>Authenticate and authorize access to applications</li>
                  <li>Improve security and prevent fraud</li>
                  <li>Send important notifications and updates</li>
                  <li>Provide customer support</li>
                  <li>Comply with legal obligations</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  4. Information Sharing and Disclosure
                </h2>
                <p>
                  <strong>
                    We do not sell, trade, or rent your personal information.
                  </strong>{" "}
                  We may share information:
                </p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>
                    <strong>With Your Consent:</strong> When you explicitly
                    authorize sharing
                  </li>
                  <li>
                    <strong>Service Providers:</strong> With trusted partners
                    who assist in operations
                  </li>
                  <li>
                    <strong>Legal Requirements:</strong> When required by law or
                    legal process
                  </li>
                  <li>
                    <strong>Security:</strong> To protect rights, property, or
                    safety
                  </li>
                  <li>
                    <strong>Business Transfer:</strong> In connection with
                    mergers or acquisitions
                  </li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  5. Data Security
                </h2>
                <p>We implement industry-standard security measures:</p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Encryption in transit and at rest</li>
                  <li>Multi-factor authentication (MFA)</li>
                  <li>Regular security audits and monitoring</li>
                  <li>Access controls and employee training</li>
                  <li>Incident response procedures</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  6. Data Retention
                </h2>
                <p>
                  We retain your information for as long as necessary to provide
                  our services and comply with legal obligations. You may
                  request deletion of your account and associated data at any
                  time.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  7. Your Rights
                </h2>
                <p>You have the right to:</p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Access your personal information</li>
                  <li>Correct inaccurate information</li>
                  <li>Delete your account and data</li>
                  <li>Export your data</li>
                  <li>Opt-out of communications</li>
                  <li>File complaints with supervisory authorities</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  8. Cookies and Tracking
                </h2>
                <p>
                  We use cookies and similar technologies to enhance your
                  experience, maintain sessions, and analyze usage patterns. You
                  can control cookie preferences through your browser settings.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  9. Third-Party Services
                </h2>
                <p>
                  Our service integrates with third-party applications. This
                  Privacy Policy does not cover third-party services. Please
                  review their privacy policies separately.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  10. International Data Transfers
                </h2>
                <p>
                  Your information may be transferred to and processed in
                  countries other than your own. We ensure appropriate
                  safeguards are in place for international transfers.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  11. Children&apos;s Privacy
                </h2>
                <p>
                  Our service is not intended for children under 13. We do not
                  knowingly collect personal information from children under 13.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  12. Changes to This Policy
                </h2>
                <p>
                  We may update this Privacy Policy from time to time. We will
                  notify you of any material changes by posting the new policy
                  on our website and updating the &quot;Last Updated&quot; date.
                </p>
              </section>
            </div>
          </div>

          {/* Footer */}
          <div className="mt-12 pt-8 border-t border-gray-200">
            <div className="text-sm text-gray-600">
              <p className="font-semibold text-gray-800 mb-3">
                Contact Information
              </p>
              <div className="space-y-2">
              <p>
                  <strong>Data Controller:</strong> Manish Kumar Saw
              </p>
              <p>
                  <strong>Email:</strong>
                <a
                  href="mailto:manishmh982@gmail.com"
                    className="text-blue-600 hover:text-blue-800 ml-1"
                >
                  manishmh982@gmail.com
                </a>
              </p>
                <p>
                  <strong>Service:</strong> CloudGate SSO Platform
                </p>
              </div>
              <p className="mt-4 text-gray-500">
                For privacy-related questions, data requests, or concerns,
                please contact us at the above email address. We will respond to
                your inquiry within 30 days.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
