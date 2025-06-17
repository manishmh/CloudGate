"use client";

import { useSearchParams } from "next/navigation";
import { Suspense, useEffect } from "react";

function OAuthCallbackContent() {
  const searchParams = useSearchParams();

  useEffect(() => {
    const code = searchParams.get("code");
    const error = searchParams.get("error");
    const state = searchParams.get("state");
    const provider = searchParams.get("provider");
    const email = searchParams.get("email");

    if (window.opener) {
      // We're in a popup window
      if (error) {
        window.opener.postMessage(
          {
            type: "oauth_error",
            error: error,
            provider: provider,
          },
          window.location.origin
        );
      } else if (code) {
        window.opener.postMessage(
          {
            type: "oauth_success",
            code: code,
            state: state,
            provider: provider,
            email: email,
          },
          window.location.origin
        );
      }
      window.close();
    } else {
      // We're in the main window (redirect flow)
      if (error) {
        window.location.href = `/dashboard/applications?error=${encodeURIComponent(
          error
        )}`;
      } else if (code && provider && email) {
        window.location.href = `/dashboard/applications?connected=${provider}&email=${email}`;
      } else {
        window.location.href = "/dashboard/applications";
      }
    }
  }, [searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="text-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p className="text-gray-600 text-lg">Processing OAuth callback...</p>
        <p className="text-gray-500 text-sm mt-2">
          Please wait while we complete the authentication.
        </p>
      </div>
    </div>
  );
}

function LoadingFallback() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="text-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p className="text-gray-600 text-lg">Loading OAuth callback...</p>
      </div>
    </div>
  );
}

export default function OAuthCallbackPage() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <OAuthCallbackContent />
    </Suspense>
  );
}
