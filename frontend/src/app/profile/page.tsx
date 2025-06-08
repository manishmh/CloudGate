"use client";

import { PROFILE_CONFIG, PROFILE_FIELDS, PROFILE_MESSAGES } from "@/constants";
import { useKeycloak } from "@react-keycloak/web";
import Image from "next/image";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

interface UserProfile {
  given_name: string;
  family_name: string;
  email: string;
  preferred_username: string;
  profile_picture?: string;
}

export default function ProfilePage() {
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
      router.replace("/profile");
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
      console.error("Failed to upload profile picture:", error);
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

      const response = await fetch("/api/send-verification", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: keycloak?.tokenParsed?.email,
          userId: keycloak?.tokenParsed?.sub,
          name: `${profile.given_name} ${profile.family_name}`,
        }),
      });

      if (response.ok) {
        setMessage({
          type: "success",
          text: PROFILE_MESSAGES.EMAIL_VERIFICATION_SENT,
        });
      } else {
        throw new Error("Failed to send verification email");
      }
    } catch (error) {
      console.error("Failed to send verification email:", error);
      setMessage({
        type: "error",
        text: PROFILE_MESSAGES.EMAIL_VERIFICATION_ERROR,
      });
    } finally {
      setSendingVerification(false);
    }
  };

  const getInitials = (firstName: string, lastName: string) => {
    return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();
  };

  const isEmailVerified = () => {
    return emailVerified !== null
      ? emailVerified
      : keycloak?.tokenParsed?.email_verified;
  };

  if (!initialized) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!keycloak?.authenticated) {
    router.push("/login");
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center">
              <button
                onClick={() => router.back()}
                className="mr-4 text-gray-600 hover:text-gray-900 cursor-pointer"
              >
                ← Back
              </button>
              <h1 className="text-3xl font-bold text-gray-900">User Profile</h1>
            </div>
            <div className="flex items-center space-x-4">
              <button
                onClick={() => router.push("/dashboard")}
                className="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer"
              >
                Dashboard
              </button>
              <button
                onClick={() => keycloak?.logout()}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto py-6 sm:px-6 lg:px-8">
        {/* Message */}
        {message && (
          <div
            className={`mb-6 p-4 rounded-md ${
              message.type === "success"
                ? "bg-green-50 border border-green-200 text-green-800"
                : "bg-red-50 border border-red-200 text-red-800"
            }`}
          >
            <div className="flex justify-between items-center">
              <span>{message.text}</span>
              <button
                onClick={() => setMessage(null)}
                className="text-gray-400 hover:text-gray-600 cursor-pointer"
              >
                ×
              </button>
            </div>
          </div>
        )}

        <div className="bg-white rounded-lg shadow-sm p-6">
          {/* Profile Picture Section */}
          <div className="flex flex-col items-center mb-8">
            <div className="relative">
              {profilePicture ? (
                <Image
                  src={profilePicture}
                  alt="Profile"
                  width={PROFILE_CONFIG.PROFILE_AVATAR_SIZE}
                  height={PROFILE_CONFIG.PROFILE_AVATAR_SIZE}
                  className="rounded-full object-cover border-4 border-gray-200"
                  unoptimized={true}
                />
              ) : (
                <div
                  className="bg-blue-600 text-white rounded-full flex items-center justify-center text-3xl font-bold border-4 border-gray-200"
                  style={{
                    width: `${PROFILE_CONFIG.PROFILE_AVATAR_SIZE}px`,
                    height: `${PROFILE_CONFIG.PROFILE_AVATAR_SIZE}px`,
                  }}
                >
                  {getInitials(profile.given_name, profile.family_name)}
                </div>
              )}

              {/* Upload Button */}
              <label className="absolute bottom-0 right-0 bg-blue-600 hover:bg-blue-700 text-white rounded-full p-2 cursor-pointer shadow-lg">
                <input
                  type="file"
                  accept={PROFILE_CONFIG.ALLOWED_FILE_TYPES.join(",")}
                  onChange={handleFileUpload}
                  className="hidden"
                  disabled={loading}
                />
                📷
              </label>
            </div>

            <h2 className="mt-4 text-xl font-semibold text-gray-900">
              {profile.given_name} {profile.family_name}
            </h2>
            <p className="text-gray-600">@{profile.preferred_username}</p>
          </div>

          {/* Profile Form */}
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="text-lg font-medium text-gray-900">
                Profile Information
              </h3>
              <button
                onClick={() => setIsEditing(!isEditing)}
                className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer"
              >
                {isEditing ? "Cancel" : "Edit"}
              </button>
            </div>

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
                    placeholder={field.placeholder}
                    disabled={!isEditing || field.readonly}
                    className={`w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 ${
                      !isEditing || field.readonly
                        ? "bg-gray-50 text-gray-500"
                        : "bg-white text-gray-900"
                    }`}
                  />
                  {field.readonly && (
                    <p className="mt-1 text-xs text-gray-500">
                      This field cannot be modified
                    </p>
                  )}
                </div>
              ))}
            </div>

            {/* Additional User Info */}
            <div className="border-t pt-6">
              <h4 className="text-md font-medium text-gray-900 mb-4">
                Account Information
              </h4>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    User ID
                  </label>
                  <p className="text-sm text-gray-900 font-mono bg-gray-50 p-2 rounded">
                    {keycloak?.tokenParsed?.sub || "N/A"}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Email Verified
                  </label>
                  <div className="flex items-center space-x-2">
                    <span
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        isEmailVerified()
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {isEmailVerified() ? "✓ Verified" : "✗ Not Verified"}
                    </span>
                    {!isEmailVerified() && (
                      <button
                        onClick={handleSendVerification}
                        disabled={sendingVerification}
                        className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-3 py-1 rounded text-xs font-medium cursor-pointer transition-colors"
                      >
                        {sendingVerification ? "Sending..." : "Verify Email"}
                      </button>
                    )}
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Roles
                  </label>
                  <div className="flex flex-wrap gap-1">
                    {keycloak?.tokenParsed?.realm_access?.roles?.map(
                      (role, index) => (
                        <span
                          key={index}
                          className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                        >
                          {role}
                        </span>
                      )
                    ) || (
                      <span className="text-sm text-gray-500">
                        No roles assigned
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Save Button */}
            {isEditing && (
              <div className="flex justify-end space-x-4 pt-6 border-t">
                <button
                  onClick={() => setIsEditing(false)}
                  className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 cursor-pointer"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  disabled={loading}
                  className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer"
                >
                  {loading ? (
                    <span className="flex items-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      Saving...
                    </span>
                  ) : (
                    "Save Changes"
                  )}
                </button>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
