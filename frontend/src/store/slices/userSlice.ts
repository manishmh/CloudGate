import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface UserProfile {
  id?: string;
  keycloakId?: string;
  email: string;
  emailVerified: boolean;
  username: string;
  firstName: string;
  lastName: string;
  profilePictureUrl?: string;
  lastLoginAt?: string;
  isActive: boolean;
}

interface UserState {
  profile: UserProfile | null;
  loading: boolean;
  error: string | null;
}

const initialState: UserState = {
  profile: null,
  loading: false,
  error: null,
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<UserProfile>) => {
      state.profile = action.payload;
      state.error = null;
    },
    updateUser: (state, action: PayloadAction<Partial<UserProfile>>) => {
      if (state.profile) {
        state.profile = { ...state.profile, ...action.payload };
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    clearUser: (state) => {
      state.profile = null;
      state.error = null;
    },
  },
});

export const { setUser, updateUser, setLoading, setError, clearUser } = userSlice.actions;
export default userSlice.reducer; 