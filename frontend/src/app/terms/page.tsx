"use client";

import Link from "next/link";
import { useEffect } from "react";

// Note: Metadata export doesn't work with 'use client', so we'll set it via Head or document title
// export const metadata: Metadata = {
//   title: "Terms of Service - CloudGate",
//   description:
//     "CloudGate Terms of Service - Terms and conditions for using our SSO platform",
// };

export default function TermsPage() {
  useEffect(() => {
    document.title = "Terms of Service - CloudGate";
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
              CloudGate Terms of Service
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
                  1. Acceptance of Terms
                </h2>
                <p>
                  By accessing or using CloudGate (&quot;Service&quot;), you
                  agree to be bound by these Terms of Service
                  (&quot;Terms&quot;). If you do not agree to these Terms, do
                  not use the Service.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  2. Description of Service
                </h2>
                <p>
                  CloudGate is a Single Sign-On (SSO) platform that provides:
                </p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Centralized authentication for multiple applications</li>
                  <li>User identity management</li>
                  <li>Security and access controls</li>
                  <li>Integration with third-party SaaS applications</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  3. User Accounts
                </h2>

                <h3 className="text-lg font-medium text-gray-800 mb-2">
                  3.1 Account Creation
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>You must provide accurate and complete information</li>
                  <li>You are responsible for maintaining account security</li>
                  <li>You must be at least 13 years old to use the Service</li>
                  <li>One person may not maintain multiple accounts</li>
                </ul>

                <h3 className="text-lg font-medium text-gray-800 mb-2 mt-4">
                  3.2 Account Security
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Keep your login credentials confidential</li>
                  <li>Notify us immediately of any unauthorized access</li>
                  <li>
                    You are responsible for all activities under your account
                  </li>
                  <li>
                    Use strong passwords and enable available security features
                  </li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  4. Acceptable Use
                </h2>

                <h3 className="text-lg font-medium text-gray-800 mb-2">
                  4.1 Permitted Uses
                </h3>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Access authorized applications through SSO</li>
                  <li>Manage your profile and preferences</li>
                  <li>Use security features as intended</li>
                </ul>

                <h3 className="text-lg font-medium text-gray-800 mb-2 mt-4">
                  4.2 Prohibited Uses
                </h3>
                <p>You may not:</p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Violate any laws or regulations</li>
                  <li>Attempt to gain unauthorized access to systems</li>
                  <li>Interfere with or disrupt the Service</li>
                  <li>Use the Service for illegal or harmful activities</li>
                  <li>Share your account credentials with others</li>
                  <li>Attempt to reverse engineer or copy the Service</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  5. Privacy and Data Protection
                </h2>
                <p>
                  Your privacy is important to us. Please review our Privacy
                  Policy, which also governs your use of the Service, to
                  understand our practices.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  6. Third-Party Integrations
                </h2>
                <p>
                  Our Service integrates with third-party applications. We are
                  not responsible for the content, privacy practices, or terms
                  of service of third-party applications. Your use of
                  third-party applications is subject to their respective terms
                  and policies.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  7. Service Availability
                </h2>
                <ul className="list-disc pl-6 space-y-1">
                  <li>We strive to maintain high service availability</li>
                  <li>
                    Scheduled maintenance may temporarily interrupt service
                  </li>
                  <li>We do not guarantee uninterrupted access</li>
                  <li>Emergency maintenance may occur without notice</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  8. Intellectual Property
                </h2>
                <p>
                  The Service and its original content, features, and
                  functionality are owned by CloudGate and are protected by
                  international copyright, trademark, patent, trade secret, and
                  other intellectual property laws.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  9. Termination
                </h2>

                <h3 className="text-lg font-medium text-gray-800 mb-2">
                  9.1 Termination by You
                </h3>
                <p>
                  You may terminate your account at any time by contacting us.
                </p>

                <h3 className="text-lg font-medium text-gray-800 mb-2 mt-4">
                  9.2 Termination by Us
                </h3>
                <p>We may terminate or suspend your account if you:</p>
                <ul className="list-disc pl-6 space-y-1">
                  <li>Violate these Terms</li>
                  <li>Engage in fraudulent or illegal activities</li>
                  <li>Pose a security risk to the Service</li>
                  <li>Have been inactive for an extended period</li>
                </ul>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  10. Disclaimers
                </h2>
                <p>
                  The Service is provided &quot;as is&quot; and &quot;as
                  available&quot; without warranties of any kind, either express
                  or implied, including but not limited to implied warranties of
                  merchantability, fitness for a particular purpose, and
                  non-infringement.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  11. Limitation of Liability
                </h2>
                <p>
                  In no event shall CloudGate be liable for any indirect,
                  incidental, special, consequential, or punitive damages,
                  including but not limited to loss of profits, data, use,
                  goodwill, or other intangible losses.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  12. Indemnification
                </h2>
                <p>
                  You agree to defend, indemnify, and hold harmless CloudGate
                  from and against any claims, damages, obligations, losses,
                  liabilities, costs, or debt arising from your use of the
                  Service or violation of these Terms.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  13. Governing Law
                </h2>
                <p>
                  These Terms shall be interpreted and governed by the laws of
                  the jurisdiction where CloudGate operates, without regard to
                  conflict of law provisions.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  14. Dispute Resolution
                </h2>
                <p>
                  Any disputes arising from these Terms or your use of the
                  Service will be resolved through binding arbitration in
                  accordance with the rules of the relevant arbitration
                  association.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  15. Changes to Terms
                </h2>
                <p>
                  We reserve the right to modify these Terms at any time. We
                  will notify users of any material changes by posting the new
                  Terms on our website and updating the &quot;Last Updated&quot;
                  date. Continued use of the Service after changes constitutes
                  acceptance of the new Terms.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  16. Severability
                </h2>
                <p>
                  If any provision of these Terms is held to be invalid or
                  unenforceable, the remaining provisions will remain in full
                  force and effect.
                </p>
              </section>

              <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-3">
                  17. Entire Agreement
                </h2>
                <p>
                  These Terms constitute the entire agreement between you and
                  CloudGate regarding the use of the Service and supersede all
                  prior agreements and understandings.
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
                  <strong>Service Provider:</strong> Manish Kumar Saw
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
                For questions about these Terms of Service, please contact us at
                the above email address. We will respond to your inquiry within
                5 business days.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
