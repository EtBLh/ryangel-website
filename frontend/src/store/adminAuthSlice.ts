import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

interface AdminAuthState {
  token: string | null;
  admin: any | null;
}

const initialState: AdminAuthState = {
  token: null,
  admin: null,
};

const adminAuthSlice = createSlice({
  name: 'adminAuth',
  initialState,
  reducers: {
    setAdminCredentials: (
      state,
      action: PayloadAction<{ token: string; admin: any }>
    ) => {
      state.token = action.payload.token;
      state.admin = action.payload.admin;
    },
    adminLogout: (state) => {
      state.token = null;
      state.admin = null;
    },
  },
});

export const { setAdminCredentials, adminLogout } = adminAuthSlice.actions;
export default adminAuthSlice.reducer;
