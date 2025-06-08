// API Configuration
export const API_CONFIG = {
  BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081',
  TIMEOUT: 10000,
  RETRY_ATTEMPTS: 3,
} as const;

// Keycloak Configuration
export const KEYCLOAK_CONFIG = {
  URL: process.env.NEXT_PUBLIC_KEYCLOAK_URL || 'http://localhost:8080',
  REALM: process.env.NEXT_PUBLIC_KEYCLOAK_REALM || 'cloudgate',
  CLIENT_ID: process.env.NEXT_PUBLIC_KEYCLOAK_CLIENT_ID || 'cloudgate-frontend',
} as const;

// Demo User Data
export const DEMO_USER = {
  id: 'demo-user-123',
  email: 'demo@cloudgate.com',
  name: 'Demo User',
  preferred_username: 'demouser',
  given_name: 'Demo',
  family_name: 'User',
  roles: ['user', 'sso-user'],
} as const;

// Fallback SaaS Applications (used when API fails)
export const FALLBACK_SAAS_APPS = [
  {
    id: "google-workspace",
    name: "Google Workspace",
    icon: "📧",
    description: "Email, Drive, Calendar, and productivity tools",
    category: "productivity",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "microsoft-365",
    name: "Microsoft 365",
    icon: "📊",
    description: "Office apps, Teams, SharePoint, and OneDrive",
    category: "productivity",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "slack",
    name: "Slack",
    icon: "💬",
    description: "Team communication and collaboration",
    category: "communication",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "salesforce",
    name: "Salesforce",
    icon: "☁️",
    description: "Customer relationship management platform",
    category: "crm",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "jira",
    name: "Jira",
    icon: "🎯",
    description: "Issue tracking and project management",
    category: "project-management",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "confluence",
    name: "Confluence",
    icon: "📝",
    description: "Team workspace and knowledge management",
    category: "documentation",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "github",
    name: "GitHub",
    icon: "🐙",
    description: "Code repository and collaboration platform",
    category: "development",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "dropbox",
    name: "Dropbox",
    icon: "📦",
    description: "Cloud storage and file synchronization",
    category: "storage",
    protocol: "oauth2",
    status: "available",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
] as const;

// Application Launch URLs
export const APP_LAUNCH_URLS = {
  "google-workspace": "https://workspace.google.com",
  "microsoft-365": "https://office.com",
  "slack": "https://slack.com/signin",
  "salesforce": "https://login.salesforce.com",
  "jira": "https://atlassian.net",
  "confluence": "https://atlassian.net",
  "github": "https://github.com",
  "dropbox": "https://dropbox.com",
} as const;

// Status Colors and Text
export const STATUS_CONFIG = {
  connected: {
    color: "bg-green-100 text-green-800",
    text: "Connected",
  },
  pending: {
    color: "bg-yellow-100 text-yellow-800",
    text: "Pending",
  },
  available: {
    color: "bg-gray-100 text-gray-800",
    text: "Available",
  },
  error: {
    color: "bg-red-100 text-red-800",
    text: "Error",
  },
} as const;

// Security Features for Dashboard
export const SECURITY_FEATURES = [
  {
    id: "mfa",
    title: "MFA Enabled",
    description: "Multi-factor authentication active",
    icon: "check",
    color: "green",
  },
  {
    id: "session",
    title: "Session Secure",
    description: "Encrypted session active",
    icon: "lock",
    color: "blue",
  },
  {
    id: "risk",
    title: "Risk Score: Low",
    description: "Adaptive security monitoring",
    icon: "info",
    color: "yellow",
  },
] as const;

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: "Network error. Please check your connection and try again.",
  AUTH_ERROR: "Authentication failed. Please login again.",
  APP_NOT_FOUND: "Application not found.",
  CONNECTION_FAILED: "Failed to connect to application. Please try again.",
  LAUNCH_FAILED: "Failed to launch application. Please try again.",
  GENERIC_ERROR: "An unexpected error occurred. Please try again.",
} as const;

// Success Messages
export const SUCCESS_MESSAGES = {
  APP_CONNECTED: "Successfully connected to application!",
  APP_LAUNCHED: "Application launched successfully!",
  DATA_REFRESHED: "Data refreshed successfully!",
} as const;

// Loading Messages
export const LOADING_MESSAGES = {
  CONNECTING: "Connecting to application...",
  LAUNCHING: "Launching application...",
  LOADING_APPS: "Loading applications...",
  REFRESHING: "Refreshing data...",
  PROCESSING: "Processing...",
} as const;

// Demo Configuration
export const DEMO_CONFIG = {
  SIMULATE_OAUTH: true,
  SHOW_DEMO_ALERTS: true,
  AUTO_CONNECT_DELAY: 2000,
  LAUNCH_DELAY: 1000,
} as const;

// Application Categories
export const APP_CATEGORIES = {
  productivity: "Productivity",
  communication: "Communication",
  crm: "Customer Relations",
  "project-management": "Project Management",
  documentation: "Documentation",
  development: "Development",
  storage: "Cloud Storage",
  analytics: "Analytics",
  security: "Security",
  finance: "Finance",
} as const;

// OAuth Configuration Templates
export const OAUTH_CONFIGS = {
  "google-workspace": {
    client_id: "your-google-client-id",
    scope: "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
    auth_url: "https://accounts.google.com/o/oauth2/v2/auth",
    token_url: "https://oauth2.googleapis.com/token",
  },
  "microsoft-365": {
    client_id: "your-microsoft-client-id",
    scope: "https://graph.microsoft.com/User.Read",
    auth_url: "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
    token_url: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
  },
  slack: {
    client_id: "your-slack-client-id",
    scope: "identity.basic,identity.email",
    auth_url: "https://slack.com/oauth/v2/authorize",
    token_url: "https://slack.com/api/oauth.v2.access",
  },
} as const;

// Theme Configuration
export const THEME_CONFIG = {
  colors: {
    primary: "blue-600",
    secondary: "gray-600",
    success: "green-600",
    warning: "yellow-600",
    error: "red-600",
  },
  spacing: {
    xs: "0.25rem",
    sm: "0.5rem",
    md: "1rem",
    lg: "1.5rem",
    xl: "2rem",
  },
} as const;

// Profile Configuration
export const PROFILE_CONFIG = {
  MAX_FILE_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_FILE_TYPES: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  AVATAR_SIZE: 40,
  PROFILE_AVATAR_SIZE: 120,
} as const;

// Profile Form Fields
export const PROFILE_FIELDS = [
  {
    id: 'given_name',
    label: 'First Name',
    type: 'text',
    required: true,
    placeholder: 'Enter your first name',
    readonly: false,
  },
  {
    id: 'family_name',
    label: 'Last Name',
    type: 'text',
    required: true,
    placeholder: 'Enter your last name',
    readonly: false,
  },
  {
    id: 'email',
    label: 'Email',
    type: 'email',
    required: true,
    placeholder: 'Enter your email address',
    readonly: true, // Email usually can't be changed in SSO
  },
  {
    id: 'preferred_username',
    label: 'Username',
    type: 'text',
    required: true,
    placeholder: 'Enter your username',
    readonly: true, // Username usually can't be changed in SSO
  },
] as const;

// Profile Messages
export const PROFILE_MESSAGES = {
  SAVE_SUCCESS: "Profile updated successfully!",
  SAVE_ERROR: "Failed to update profile. Please try again.",
  UPLOAD_SUCCESS: "Profile picture updated successfully!",
  UPLOAD_ERROR: "Failed to upload profile picture. Please try again.",
  FILE_TOO_LARGE: "File size must be less than 5MB.",
  INVALID_FILE_TYPE: "Please upload a valid image file (JPEG, PNG, GIF, or WebP).",
  VALIDATION_ERROR: "Please fill in all required fields.",
  EMAIL_VERIFICATION_SENT: "Verification email sent! Please check your inbox.",
  EMAIL_VERIFICATION_ERROR: "Failed to send verification email. Please try again.",
  EMAIL_VERIFIED_SUCCESS: "Email verified successfully!",
  EMAIL_VERIFICATION_INVALID: "Invalid verification link.",
} as const;

// Navigation Items
export const NAV_ITEMS = [
  {
    id: 'dashboard',
    label: 'Dashboard',
    href: '/dashboard',
    icon: '🏠',
  },
  {
    id: 'profile',
    label: 'Profile',
    href: '/profile',
    icon: '��',
  },
] as const; 