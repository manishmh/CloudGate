import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface SidebarState {
  isOpen: boolean;
  isPinned: boolean;
}

const initialState: SidebarState = {
  isOpen: false,
  isPinned: false,
};

const sidebarSlice = createSlice({
  name: 'sidebar',
  initialState,
  reducers: {
    toggleSidebar: (state) => {
      state.isOpen = !state.isOpen;
    },
    setSidebarOpen: (state, action: PayloadAction<boolean>) => {
      state.isOpen = action.payload;
    },
    togglePinned: (state) => {
      state.isPinned = !state.isPinned;
    },
    setPinned: (state, action: PayloadAction<boolean>) => {
      state.isPinned = action.payload;
    },
  },
});

export const { toggleSidebar, setSidebarOpen, togglePinned, setPinned } = sidebarSlice.actions;
export default sidebarSlice.reducer; 