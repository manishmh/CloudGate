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