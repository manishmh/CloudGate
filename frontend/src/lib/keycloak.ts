import { KeycloakConfig } from '@/types/auth';
import Keycloak from 'keycloak-js';

// Keycloak configuration
const keycloakConfig: KeycloakConfig = {
  url: process.env.NEXT_PUBLIC_KEYCLOAK_URL || 'http://localhost:8080',
  realm: process.env.NEXT_PUBLIC_KEYCLOAK_REALM || 'cloudgate',
  clientId: process.env.NEXT_PUBLIC_KEYCLOAK_CLIENT_ID || 'cloudgate-frontend',
};

// Initialize Keycloak instance with explicit configuration
const keycloak = new Keycloak({
  url: keycloakConfig.url,
  realm: keycloakConfig.realm,
  clientId: keycloakConfig.clientId,
});

// Keycloak initialization options
export const keycloakInitOptions = {
  onLoad: 'check-sso' as const,
  checkLoginIframe: false,
  pkceMethod: 'S256' as const,
  enableLogging: true,
  flow: 'standard' as const,
  responseMode: 'fragment' as const,
  messageReceiveTimeout: 10000,
  silentCheckSsoRedirectUri: typeof window !== 'undefined' ? `${window.location.origin}/silent-check-sso.html` : undefined,
};

export { keycloak, keycloakConfig };

// Helper functions
export const getKeycloakToken = (): string | null => {
  return keycloak.token || null;
};

export const isTokenExpired = (): boolean => {
  return keycloak.isTokenExpired();
};

export const refreshToken = async (): Promise<boolean> => {
  try {
    const refreshed = await keycloak.updateToken(30);
    return refreshed;
  } catch (error) {
    console.error('Failed to refresh token:', error);
    return false;
  }
};

export const logout = (): void => {
  keycloak.logout({
    redirectUri: typeof window !== 'undefined' ? `${window.location.origin}/login` : undefined,
  });
};

export const login = (): void => {
  keycloak.login({
    redirectUri: typeof window !== 'undefined' ? `${window.location.origin}/dashboard` : undefined,
  });
}; 