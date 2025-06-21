import { API_CONFIG, ERROR_MESSAGES } from '@/constants';

// OAuth endpoint mapping - maps app IDs to OAuth endpoint paths
const OAUTH_ENDPOINT_MAP: Record<string, string> = {
  'google-workspace': 'google',
  'microsoft-365': 'microsoft',
  'slack': 'slack',
  'salesforce': 'salesforce',
  'jira': 'jira',
  'trello': 'trello',
  'notion': 'notion',
  'github': 'github',
  'dropbox': 'dropbox',
};

// Types for API responses
export interface SaaSApplication {
  id: string;
  name: string;
  icon: string;
  description: string;
  category: string;
  protocol: string;
  status: 'available' | 'connected' | 'pending' | 'error';
  created_at: string;
  updated_at: string;
  connection_details?: {
    user_email?: string;
    user_name?: string;
    connected_at?: string;
    last_used?: string;
  };
}

export interface AppConnectionResponse {
  auth_url: string;
  state: string;
  provider: string;
}

export interface AppLaunchResponse {
  launch_url: string;
  method: 'redirect' | 'popup' | 'iframe';
  token?: string;
  expires_in?: number;
}

export interface AppsResponse {
  apps: SaaSApplication[];
  count: number;
}

// MFA Types
export interface MFASetupResponse {
  secret: string;
  qr_code_url: string;
  qr_code_data_url: string;
  backup_codes: string[];
}

export interface MFAStatusResponse {
  enabled: boolean;
  setup_date?: string;
  backup_codes_remaining: number;
}

export interface MFAVerifyRequest {
  code: string;
}

// Dashboard types
export interface UserProfile {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  profile_picture_url?: string;
  last_login_at?: string;
}

export interface DashboardMetrics {
  total_apps: number;
  connected_apps: number;
  recent_logins: number;
  security_score: number;
  last_activity: string;
}

export interface AppConnection {
  name: string;
  status: 'connected' | 'disconnected';
  icon: string;
  description: string;
  connect_url: string;
  last_used?: string;
}

export interface ActivityItem {
  id: string;
  type: 'login' | 'app_launch' | 'connection' | 'security';
  description: string;
  timestamp: string;
  icon: string;
  severity: 'info' | 'warning' | 'success';
}

export interface FeatureCard {
  id: string;
  title: string;
  description: string;
  icon: string;
  stats: string;
  color: string;
  features: string[];
}

export interface DashboardData {
  user: UserProfile;
  metrics: DashboardMetrics;
  connections: AppConnection[];
  recent_activity: ActivityItem[];
  features: FeatureCard[];
}

export interface DashboardResponse {
  success: boolean;
  data: DashboardData;
}

export interface MetricsResponse {
  success: boolean;
  metrics: DashboardMetrics;
}

// API client class
class ApiClient {
  private baseURL: string;
  private timeout: number;

  constructor() {
    this.baseURL = API_CONFIG.BASE_URL;
    this.timeout = API_CONFIG.TIMEOUT;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    };

    // Add auth header if available
    const token = this.getAuthToken();
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      };
    }

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), this.timeout);

      const response = await fetch(url, {
        ...config,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        if (error.name === 'AbortError') {
          throw new Error(ERROR_MESSAGES.NETWORK_ERROR);
        }
        throw new Error(error.message);
      }
      throw new Error(ERROR_MESSAGES.GENERIC_ERROR);
    }
  }

  private getAuthToken(): string | null {
    // In a real app, this would get the token from your auth provider (Keycloak)
    // For now, we'll use a demo token - this should be replaced with real Keycloak integration
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token') || 'demo-user-token';
    }
    return 'demo-user-token';
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string; service: string }> {
    return this.request<{ status: string; timestamp: string; service: string }>('/health');
  }

  // Apps API methods
  async getApps(): Promise<AppsResponse> {
    return this.request<AppsResponse>('/apps');
  }

  async connectApp(appId: string): Promise<AppConnectionResponse> {
    // Map app ID to OAuth endpoint path
    const oauthEndpoint = OAUTH_ENDPOINT_MAP[appId] || appId;
    return this.request<AppConnectionResponse>(`/oauth/${oauthEndpoint}/connect`);
  }

  async launchApp(appId: string): Promise<AppLaunchResponse> {
    return this.request<AppLaunchResponse>('/apps/launch', {
      method: 'POST',
      body: JSON.stringify({ app_id: appId }),
    });
  }

  // MFA API methods
  async getMFAStatus(): Promise<MFAStatusResponse> {
    return this.request<MFAStatusResponse>('/user/mfa/status');
  }

  async setupMFA(): Promise<MFASetupResponse> {
    return this.request<MFASetupResponse>('/user/mfa/setup', {
      method: 'POST',
    });
  }

  async verifyMFASetup(code: string): Promise<{ message: string; enabled: boolean }> {
    return this.request<{ message: string; enabled: boolean }>('/user/mfa/verify-setup', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
  }

  async verifyMFA(code: string): Promise<{ message: string; verified: boolean }> {
    return this.request<{ message: string; verified: boolean }>('/user/mfa/verify', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
  }

  async disableMFA(code: string): Promise<{ message: string; enabled: boolean }> {
    return this.request<{ message: string; enabled: boolean }>('/user/mfa/disable', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
  }

  async regenerateBackupCodes(code: string): Promise<{ message: string; backup_codes: string[] }> {
    return this.request<{ message: string; backup_codes: string[] }>('/user/mfa/backup-codes/regenerate', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
  }

  // Dashboard methods
  async getDashboardData(): Promise<DashboardResponse> {
    return this.request<DashboardResponse>('/dashboard/data');
  }

  async getDashboardMetrics(): Promise<MetricsResponse> {
    return this.request<MetricsResponse>('/dashboard/metrics');
  }

  // User profile methods
  async getUserProfile(): Promise<UserProfile> {
    return this.request<UserProfile>('/user/profile');
  }

  async updateUserProfile(profile: Partial<UserProfile>): Promise<UserProfile> {
    return this.request<UserProfile>('/user/profile', {
      method: 'PUT',
      body: JSON.stringify(profile),
    });
  }

  // Utility method to check if backend is available
  async isBackendAvailable(): Promise<boolean> {
    try {
      await this.healthCheck();
      return true;
    } catch {
      return false;
    }
  }

  // OAuth Monitoring APIs
  async getConnections(): Promise<{ connections: EnhancedConnection[]; count: number }> {
    return this.request<{ connections: EnhancedConnection[]; count: number }>('/user/monitoring/connections');
  }

  async getConnectionStats(): Promise<ConnectionStats> {
    return this.request<ConnectionStats>('/user/monitoring/connections/stats');
  }

  async testConnection(connectionId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/user/monitoring/connections/${connectionId}/test`, {
      method: 'POST',
    });
  }

  async recordUsage(connectionId: string, dataTransferred?: number): Promise<{ message: string }> {
    return this.request<{ message: string }>('/user/monitoring/connections/usage', {
      method: 'POST',
      body: JSON.stringify({
        connection_id: connectionId,
        data_transferred: dataTransferred || 0,
      }),
    });
  }

  // Security Events APIs
  async getSecurityEvents(limit?: number): Promise<{ events: SecurityEvent[]; count: number }> {
    const params = limit ? `?limit=${limit}` : '';
    return this.request<{ events: SecurityEvent[]; count: number }>(`/user/monitoring/security/events${params}`);
  }

  async createSecurityEvent(event: CreateSecurityEventRequest): Promise<{ message: string }> {
    return this.request<{ message: string }>('/user/monitoring/security/events', {
      method: 'POST',
      body: JSON.stringify(event),
    });
  }

  // Device Management APIs
  async getTrustedDevices(): Promise<{ devices: TrustedDevice[]; count: number }> {
    return this.request<{ devices: TrustedDevice[]; count: number }>('/user/monitoring/devices');
  }

  async registerDevice(device: RegisterDeviceRequest): Promise<{ message: string }> {
    return this.request<{ message: string }>('/user/monitoring/devices', {
      method: 'POST',
      body: JSON.stringify(device),
    });
  }

  async trustDevice(deviceId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/user/monitoring/devices/${deviceId}/trust`, {
      method: 'PUT',
    });
  }

  async revokeDevice(deviceId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/user/monitoring/devices/${deviceId}`, {
      method: 'DELETE',
    });
  }

  // User Settings APIs
  async getUserSettings(): Promise<{ settings: UserSettings }> {
    return this.request<{ settings: UserSettings }>('/user/settings');
  }

  async updateUserSettings(settings: Partial<UserSettings>): Promise<{ settings: UserSettings }> {
    return this.request<{ settings: UserSettings }>('/user/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  async updateSingleSetting(key: string, value: boolean | number | string): Promise<{ settings: UserSettings }> {
    return this.request<{ settings: UserSettings }>('/user/settings/single', {
      method: 'PUT',
      body: JSON.stringify({ key, value }),
    });
  }

  async resetUserSettings(): Promise<UserSettings> {
    return this.request<UserSettings>('/user/settings/reset', { method: 'POST' });
  }

  // Security Monitoring API methods
  async getSecurityAlerts(): Promise<{ alerts: SecurityAlert[] }> {
    return this.request<{ alerts: SecurityAlert[] }>('/api/v1/security/alerts');
  }

  async updateSecurityAlertStatus(alertId: string, status: SecurityAlert['status']): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/security/alerts/${alertId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  }

  // Adaptive Authentication API methods
  async evaluateAuthentication(data: {
    device_fingerprint: string;
    typing_pattern?: {
      avg_keydown_time?: number;
      key_intervals?: number[];
    };
    mouse_pattern?: {
      click_intervals?: number[];
      movement_speed?: number;
    };
  }): Promise<AdaptiveAuthResponse> {
    // Get current user info and IP address
    const userAgent = navigator.userAgent;
    const ipAddress = '127.0.0.1'; // In production, this would be detected server-side
    
    // For demo purposes, use demo user data
    const requestData = {
      user_id: '12345678-1234-1234-1234-123456789012', // Valid UUID for demo user
      email: 'demo@cloudgate.com', // This should come from Keycloak token in production
      ip_address: ipAddress,
      user_agent: userAgent,
      device_fingerprint: data.device_fingerprint,
      location: {
        country: 'US',
        region: 'California',
        city: 'San Francisco',
        latitude: 37.7749,
        longitude: -122.4194,
        isp: 'Demo ISP',
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        vpn_detected: false
      },
      session_info: {
        typing_pattern: data.typing_pattern,
        mouse_pattern: data.mouse_pattern
      },
      request_headers: {
        'User-Agent': userAgent,
        'Accept-Language': navigator.language
      },
      application_id: 'cloudgate-dashboard'
    };

    return this.request<AdaptiveAuthResponse>('/api/v1/adaptive-auth/evaluate', {
      method: 'POST',
      body: JSON.stringify(requestData),
    });
  }

  async getRiskAssessmentHistory(userId: string, limit?: number): Promise<{ assessments: RiskAssessment[] }> {
    const params = limit ? `?limit=${limit}` : '';
    return this.request<{ assessments: RiskAssessment[] }>(`/api/v1/adaptive-auth/history/${userId}${params}`);
  }

  async getLatestRiskAssessment(userId: string): Promise<RiskAssessment> {
    return this.request<RiskAssessment>(`/api/v1/adaptive-auth/latest/${userId}`);
  }

  async updateRiskThresholds(thresholds: RiskThresholds): Promise<{ message: string }> {
    return this.request<{ message: string }>('/api/v1/adaptive-auth/thresholds', {
      method: 'PUT',
      body: JSON.stringify(thresholds),
    });
  }

  async registerDeviceFingerprint(fingerprint: string): Promise<{ message: string }> {
    return this.request<{ message: string }>('/api/v1/adaptive-auth/register-device', {
      method: 'POST',
      body: JSON.stringify({ 
        user_id: '12345678-1234-1234-1234-123456789012', // Valid UUID for demo user
        fingerprint: fingerprint,
        device_name: 'Browser',
        device_type: 'web',
        browser: navigator.userAgent,
        os: navigator.platform
      }),
    });
  }

  async checkDeviceStatus(fingerprint: string): Promise<{ status: string; trusted: boolean }> {
    return this.request<{ status: string; trusted: boolean }>(`/api/v1/adaptive-auth/device-status?fingerprint=${fingerprint}`);
  }

  // WebAuthn API methods
  async webAuthnRegisterBegin(): Promise<unknown> {
    return this.request<unknown>('/webauthn/register/begin', {
      method: 'POST',
    });
  }

  async webAuthnRegisterFinish(credential: unknown): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/webauthn/register/finish', {
      method: 'POST',
      body: JSON.stringify({ credential }),
    });
  }

  async webAuthnAuthenticateBegin(): Promise<unknown> {
    return this.request<unknown>('/webauthn/authenticate/begin', {
      method: 'POST',
    });
  }

  async webAuthnAuthenticateFinish(credential: unknown): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/webauthn/authenticate/finish', {
      method: 'POST',
      body: JSON.stringify({ credential }),
    });
  }

  async getWebAuthnCredentials(): Promise<{ credentials: unknown[] }> {
    return this.request<{ credentials: unknown[] }>('/webauthn/credentials');
  }

  async deleteWebAuthnCredential(credentialId: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/webauthn/credentials/${credentialId}`, {
      method: 'DELETE',
    });
  }
}

// Additional types for OAuth monitoring
export interface EnhancedConnection {
  id: string;
  user_id: string;
  app_id: string;
  app_name: string;
  provider: string;
  status: string;
  user_email?: string;
  user_name?: string;
  connected_at: string;
  last_used?: string;
  health: ConnectionHealth;
  usage_count: number;
  data_transferred: string;
  created_at: string;
  updated_at: string;
}

export interface SecurityAlert {
  id: string;
  user_id: string;
  type: string;
  severity: "low" | "medium" | "high" | "critical";
  description: string;
  status: "new" | "acknowledged" | "resolved";
  source: string;
  created_at: string;
}

export interface ConnectionHealth {
  status: string;
  last_check: string;
  response_time: number;
  uptime: number;
  error_count: number;
}

export interface ConnectionStats {
  total_connections: number;
  active_connections: number;
  failed_connections: number;
  average_response_time: number;
  uptime_percentage: number;
}

export interface SecurityEvent {
  id: string;
  user_id: string;
  connection_id?: string;
  event_type: string;
  description: string;
  severity: string;
  ip_address: string;
  user_agent: string;
  location?: string;
  risk_score: number;
  resolved: boolean;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSecurityEventRequest {
  event_type: string;
  description: string;
  severity: string;
  location?: string;
  risk_score?: number;
  connection_id?: string;
}

export interface TrustedDevice {
  id: string;
  user_id: string;
  device_name: string;
  device_type: string;
  browser: string;
  os: string;
  fingerprint: string;
  ip_address: string;
  location?: string;
  trusted: boolean;
  last_seen: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterDeviceRequest {
  device_name: string;
  device_type: string;
  browser: string;
  os: string;
  fingerprint: string;
  location?: string;
}

export interface UserSettings {
  id: string;
  user_id: string;
  language: string;
  timezone: string;
  theme: string;
  email_notifications: boolean;
  push_notifications: boolean;
  security_alerts: boolean;
  marketing_emails: boolean;
  session_timeout: number;
  auto_logout: boolean;
  two_factor_enabled: boolean;
  backup_codes_remaining: number;
  password_expiry_days: number;
  login_notifications: boolean;
  suspicious_activity_alerts: boolean;
  created_at: string;
  updated_at: string;
}

export interface RiskAssessment {
  user_id: string;
  risk_score: number;
  risk_level: "low" | "medium" | "high" | "critical";
  risk_factors: RiskFactor[];
  location: {
    country: string;
    city: string;
    is_vpn: boolean;
    is_tor: boolean;
  };
  device: {
    fingerprint: string;
    is_known: boolean;
    trust_score: number;
  };
  behavior: {
    typing_speed_deviation: number;
    mouse_pattern_deviation: number;
  };
  timestamp: string;
}

export interface RiskFactor {
  type: string;
  description: string;
  weight: number;
  score: number;
}

export interface AdaptiveAuthResponse {
  decision: "allow" | "challenge" | "deny" | "monitor";
  risk_score: number;
  risk_level: "low" | "medium" | "high" | "critical";
  required_actions: AuthAction[];
  reasoning: string[];
  session_duration_seconds: number;
  restrictions: AuthRestriction[];
  metadata: Record<string, unknown>;
  expires_at: string;
  // Legacy compatibility - map to new format
  risk_assessment?: RiskAssessment;
  session_restrictions?: {
    max_duration_minutes: number;
    require_mfa: boolean;
    allowed_operations: string[];
  };
}

export interface AuthAction {
  type: string;
  required: boolean;
  timeout_seconds: number;
  metadata: Record<string, unknown>;
  description: string;
}

export interface AuthRestriction {
  type: string;
  value: unknown;
  description: string;
  expires_at?: string;
}

export interface RiskThresholds {
  low_threshold: number;
  medium_threshold: number;
  high_threshold: number;
  critical_threshold: number;
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export individual methods for convenience
export const {
  getApps,
  connectApp,
  launchApp,
  healthCheck,
  getDashboardData,
  getDashboardMetrics,
  getMFAStatus,
  setupMFA,
  verifyMFASetup,
  verifyMFA,
  disableMFA,
  regenerateBackupCodes,
  getUserProfile,
  updateUserProfile,
  isBackendAvailable,
} = apiClient; 