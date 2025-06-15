import { API_CONFIG, ERROR_MESSAGES } from '@/constants';

// Types for API responses
export interface SaaSApplication {
  id: string;
  name: string;
  icon: string;
  description: string;
  category: string;
  protocol: string;
  status: 'available' | 'connected' | 'pending';
  created_at: string;
  updated_at: string;
}

export interface AppConnectionResponse {
  auth_url: string;
  state: string;
  challenge?: string;
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
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
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
    // In a real app, this would get the token from your auth provider
    // For demo purposes, we'll use a dummy token
    return 'demo-token';
  }

  // API methods
  async getApps(): Promise<AppsResponse> {
    return this.request<AppsResponse>('/apps');
  }

  async connectApp(appId: string): Promise<AppConnectionResponse> {
    return this.request<AppConnectionResponse>('/apps/connect', {
      method: 'POST',
      body: JSON.stringify({ app_id: appId }),
    });
  }

  async launchApp(appId: string): Promise<AppLaunchResponse> {
    return this.request<AppLaunchResponse>('/apps/launch', {
      method: 'POST',
      body: JSON.stringify({ app_id: appId }),
    });
  }

  async healthCheck(): Promise<{ status: string; timestamp: string; service: string }> {
    return this.request<{ status: string; timestamp: string; service: string }>('/health');
  }

  // Dashboard methods
  async getDashboardData(): Promise<DashboardResponse> {
    return this.request<DashboardResponse>('/dashboard/data');
  }

  async getDashboardMetrics(): Promise<MetricsResponse> {
    return this.request<MetricsResponse>('/dashboard/metrics');
  }
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
} = apiClient; 