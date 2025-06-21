"use client";

import DashboardLayout from "@/components/DashboardLayout";
import { PROFILE_CONFIG, PROFILE_FIELDS, PROFILE_MESSAGES } from "@/constants";
import { useKeycloak } from "@react-keycloak/web";
import Image from "next/image";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

interface UserProfile {
  given_name: string;
  family_name: string;
  email: string;
  preferred_username: string;
  profile_picture?: string;
}

function ProfileContent() {
  const { keycloak, initialized } = useKeycloak();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [profile, setProfile] = useState<UserProfile>({
    given_name: "",
    family_name: "",
    email: "",
    preferred_username: "",
  });
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);
  const [profilePicture, setProfilePicture] = useState<string | null>(null);
  const [sendingVerification, setSendingVerification] = useState(false);
  const [emailVerified, setEmailVerified] = useState<boolean | null>(null);

  useEffect(() => {
    if (initialized && keycloak?.authenticated && keycloak.tokenParsed) {
      const userData = keycloak.tokenParsed;
      setProfile({
        given_name: userData.given_name || "",
        family_name: userData.family_name || "",
        email: userData.email || "",
        preferred_username: userData.preferred_username || "",
      });

      // Load profile picture from localStorage (in a real app, this would come from backend)
      const savedPicture = localStorage.getItem(
        `profile_picture_${userData.sub}`
      );
      if (savedPicture) {
        setProfilePicture(savedPicture);
      }

      // Check if email was verified locally
      const localEmailVerified = localStorage.getItem(
        `email_verified_${userData.sub}`
      );
      if (localEmailVerified === "true") {
        setEmailVerified(true);
      } else {
        setEmailVerified(userData.email_verified || false);
      }
    }

    // Check for verification status in URL
    const verification = searchParams.get("verification");
    if (verification) {
      if (verification === "success") {
        setMessage({
          type: "success",
          text: PROFILE_MESSAGES.EMAIL_VERIFIED_SUCCESS,
        });

        // Mark email as verified locally
        const userId = keycloak?.tokenParsed?.sub;
        if (userId) {
          localStorage.setItem(`email_verified_${userId}`, "true");
          setEmailVerified(true);
        }
      } else if (verification === "invalid") {
        setMessage({
          type: "error",
          text: PROFILE_MESSAGES.EMAIL_VERIFICATION_INVALID,
        });
      } else if (verification === "error") {
        setMessage({
          type: "error",
          text: PROFILE_MESSAGES.EMAIL_VERIFICATION_ERROR,
        });
      }

      // Clean up URL
      router.replace("/dashboard/profile");
    }
  }, [initialized, keycloak, searchParams, router]);

  const handleInputChange = (field: string, value: string) => {
    setProfile((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleSave = async () => {
    try {
      setLoading(true);
      setMessage(null);

      // Validate required fields
      const requiredFields = PROFILE_FIELDS.filter((field) => field.required);
      const missingFields = requiredFields.filter(
        (field) => !profile[field.id as keyof UserProfile]
      );

      if (missingFields.length > 0) {
        setMessage({
          type: "error",
          text: PROFILE_MESSAGES.VALIDATION_ERROR,
        });
        return;
      }

      // In a real application, you would send this to your backend
      // For now, we'll simulate a save operation
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Save to localStorage (in a real app, this would be saved to backend)
      const userId = keycloak?.tokenParsed?.sub;
      if (userId) {
        localStorage.setItem(`user_profile_${userId}`, JSON.stringify(profile));
      }

      setMessage({
        type: "success",
        text: PROFILE_MESSAGES.SAVE_SUCCESS,
      });
      setIsEditing(false);
    } catch (error) {
      console.error("Failed to save profile:", error);
      setMessage({
        type: "error",
        text: PROFILE_MESSAGES.SAVE_ERROR,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      setLoading(true);
      setMessage(null);

      // Validate file size
      if (file.size > PROFILE_CONFIG.MAX_FILE_SIZE) {
        setMessage({
          type: "error",
          text: PROFILE_MESSAGES.FILE_TOO_LARGE,
        });
        return;
      }

      // Validate file type
      if (
        !(PROFILE_CONFIG.ALLOWED_FILE_TYPES as readonly string[]).includes(
          file.type
        )
      ) {
        setMessage({
          type: "error",
          text: PROFILE_MESSAGES.INVALID_FILE_TYPE,
        });
        return;
      }

      // Convert to base64 for storage (in a real app, you'd upload to a server)
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64 = e.target?.result as string;

        // Simulate upload delay
        await new Promise((resolve) => setTimeout(resolve, 1000));

        setProfilePicture(base64);

        // Save to localStorage (in a real app, this would be saved to backend)
        const userId = keycloak?.tokenParsed?.sub;
        if (userId) {
          localStorage.setItem(`profile_picture_${userId}`, base64);
        }

        setMessage({
          type: "success",
          text: PROFILE_MESSAGES.UPLOAD_SUCCESS,
        });
      };

      reader.readAsDataURL(file);
    } catch (error) {
      console.error("Failed to upload file:", error);
      setMessage({
        type: "error",
        text: PROFILE_MESSAGES.UPLOAD_ERROR,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleSendVerification = async () => {
    try {
      setSendingVerification(true);
      setMessage(null);

      // Simulate sending verification email
      await new Promise((resolve) => setTimeout(resolve, 2000));

      setMessage({
        type: "success",
        text: PROFILE_MESSAGES.VERIFICATION_SENT,
      });
    } catch (error) {
      console.error("Failed to send verification:", error);
      setMessage({
        type: "error",
        text: PROFILE_MESSAGES.VERIFICATION_ERROR,
      });
    } finally {
      setSendingVerification(false);
    }
  };

  const getInitials = (firstName: string, lastName: string) => {
    return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();
  };

  const isEmailVerified = () => {
    return emailVerified === true;
  };

  const saveAction = (
    <div className="flex items-center space-x-2">
      {isEditing && (
        <>
          <button
            onClick={() => setIsEditing(false)}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 cursor-pointer"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={loading}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50 cursor-pointer"
          >
            {loading ? "Saving..." : "Save Changes"}
          </button>
        </>
      )}
      {!isEditing && (
        <button
          onClick={() => setIsEditing(true)}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 cursor-pointer"
        >
          Edit Profile
        </button>
      )}
    </div>
  );

  return (
    <DashboardLayout>
      {/* Message Display */}
      {message && (
        <div
          className={`mb-6 p-4 rounded-md ${
            message.type === "success"
              ? "bg-green-50 text-green-800 border border-green-200"
              : "bg-red-50 text-red-800 border border-red-200"
          }`}
        >
          <div className="flex justify-between items-center">
            <p className="text-sm">{message.text}</p>
            <button
              onClick={() => setMessage(null)}
              className="text-gray-400 hover:text-gray-600 cursor-pointer"
            >
              ×
            </button>
          </div>
        </div>
      )}

      <div className="max-w-4xl mx-auto">
        <div className="bg-white shadow rounded-lg">
          {/* Profile Header */}
          <div className="px-6 py-8 border-b border-gray-200">
            <div className="flex items-center space-x-6">
              <div className="relative">
                <div className="h-24 w-24 rounded-full overflow-hidden bg-gray-100 flex items-center justify-center">
                  {profilePicture ? (
                    <Image
                      src={profilePicture}
                      alt="Profile"
                      width={96}
                      height={96}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <span className="text-2xl font-medium text-gray-600">
                      {getInitials(profile.given_name, profile.family_name)}
                    </span>
                  )}
                </div>
                {isEditing && (
                  <label className="absolute bottom-0 right-0 bg-blue-600 rounded-full p-2 cursor-pointer hover:bg-blue-700">
                    <svg
                      className="h-4 w-4 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                    </svg>
                    <input
                      type="file"
                      className="hidden"
                      accept="image/*"
                      onChange={handleFileUpload}
                    />
                  </label>
                )}
              </div>
              <div>
                <h2 className="text-2xl font-bold text-gray-900">
                  {profile.given_name} {profile.family_name}
                </h2>
                <p className="text-gray-600">@{profile.preferred_username}</p>
                <div className="flex items-center mt-2">
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      isEmailVerified()
                        ? "bg-green-100 text-green-800"
                        : "bg-yellow-100 text-yellow-800"
                    }`}
                  >
                    {isEmailVerified() ? "Email Verified" : "Email Unverified"}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Profile Form */}
          <div className="px-6 py-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {PROFILE_FIELDS.map((field) => (
                <div key={field.id}>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {field.label}
                    {field.required && (
                      <span className="text-red-500 ml-1">*</span>
                    )}
                  </label>
                  <input
                    type={field.type}
                    value={profile[field.id as keyof UserProfile] || ""}
                    onChange={(e) =>
                      handleInputChange(field.id, e.target.value)
                    }
                    disabled={!isEditing || field.disabled}
                    placeholder={field.placeholder}
                    className={`w-full px-3 py-2 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 ${
                      !isEditing || field.disabled
                        ? "bg-gray-50 text-gray-500"
                        : "bg-white text-gray-900"
                    } border-gray-300`}
                  />
                  {field.helpText && (
                    <p className="mt-1 text-sm text-gray-500">
                      {field.helpText}
                    </p>
                  )}
                </div>
              ))}
            </div>

            {/* Email Verification Section */}
            {!isEmailVerified() && (
              <div className="mt-8 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-sm font-medium text-yellow-800">
                      Email Verification Required
                    </h3>
                    <p className="text-sm text-yellow-700 mt-1">
                      Please verify your email address to secure your account.
                    </p>
                  </div>
                  <button
                    onClick={handleSendVerification}
                    disabled={sendingVerification}
                    className="px-4 py-2 text-sm font-medium text-yellow-800 bg-yellow-100 border border-yellow-300 rounded-md hover:bg-yellow-200 disabled:opacity-50 cursor-pointer"
                  >
                    {sendingVerification ? "Sending..." : "Send Verification"}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}

export default function ProfilePage() {
  return (
    <Suspense
      fallback={
        <DashboardLayout>
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        </DashboardLayout>
      }
    >
      <ProfileContent />
    </Suspense>
  );
}
