import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface UserSettings {
  // Appearance
  language: string;
  timezone: string;
  dateFormat: string;

  // Notifications
  emailNotifications: boolean;
  pushNotifications: boolean;
  securityAlerts: boolean;
  appUpdates: boolean;
  weeklyReports: boolean;

  // Dashboard
  defaultView: string;
  itemsPerPage: number;
  autoRefresh: boolean;
  refreshInterval: number;

  // Privacy
  analyticsOptIn: boolean;
  shareUsageData: boolean;
  personalizedAds: boolean;

  // Integration
  apiAccess: boolean;
  webhookUrl: string;
  maxApiCalls: number;
}

interface SettingsState {
  settings: UserSettings;
  loading: boolean;
  error: string | null;
  lastSaved: string | null;
}

const initialState: SettingsState = {
  settings: {
    // Appearance
    language: 'en',
    timezone: 'America/New_York',
    dateFormat: 'MM/DD/YYYY',

    // Notifications
    emailNotifications: true,
    pushNotifications: false,
    securityAlerts: true,
    appUpdates: true,
    weeklyReports: false,

    // Dashboard
    defaultView: 'dashboard',
    itemsPerPage: 10,
    autoRefresh: true,
    refreshInterval: 30,

    // Privacy
    analyticsOptIn: true,
    shareUsageData: false,
    personalizedAds: false,

    // Integration
    apiAccess: false,
    webhookUrl: '',
    maxApiCalls: 1000,
  },
  loading: false,
  error: null,
  lastSaved: null,
};

const settingsSlice = createSlice({
  name: 'settings',
  initialState,
  reducers: {
    setSettings: (state, action: PayloadAction<UserSettings>) => {
      state.settings = action.payload;
      state.error = null;
    },
    updateSetting: (state, action: PayloadAction<{ key: keyof UserSettings; value: unknown }>) => {
      const { key, value } = action.payload;
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (state.settings as any)[key] = value;
      state.lastSaved = null;
      state.error = null;
    },
    updateSettings: (state, action: PayloadAction<Partial<UserSettings>>) => {
      state.settings = { ...state.settings, ...action.payload };
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setSaved: (state) => {
      state.lastSaved = new Date().toISOString();
      state.error = null;
    },
    resetSettings: (state) => {
      state.settings = initialState.settings;
    },
  },
});

export const { 
  setSettings, 
  updateSetting, 
  updateSettings, 
  setLoading, 
  setError, 
  setSaved, 
  resetSettings 
} = settingsSlice.actions;
export default settingsSlice.reducer; 