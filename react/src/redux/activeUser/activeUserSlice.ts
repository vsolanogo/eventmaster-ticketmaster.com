import {
  signUpUser,
  signInUser,
  logoutUser,
  getCurrentUser,
} from "../../api/activeUserApi";
import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { ErrorMessage, User } from "../../models/UserModels";
import { ApiStatus } from "../../constants";

export type ActiveUserState = {
  loginUserStatus: ApiStatus;
  registerUserStatus: ApiStatus;
  logoutUserStatus: ApiStatus;
  user: User["id"] | null;
  isAuthenticated: boolean;
  error: ErrorMessage | null;
};

export const activeUserState: ActiveUserState = {
  loginUserStatus: "IDLE",
  registerUserStatus: "IDLE",
  logoutUserStatus: "IDLE",
  user: null,
  isAuthenticated: false,
  error: null,
};

export const signInUserThunk = createAsyncThunk("activeUser/login", signInUser);

export const signUpUserThunk = createAsyncThunk(
  "activeUser/register",
  signUpUser
);

export const logoutUserThunk = createAsyncThunk(
  "activeUser/logout",
  logoutUser
);

export const getCurrentUserThunk = createAsyncThunk(
  "activeUser/getCurrent",
  getCurrentUser
);

export const activeUserSlice = createSlice({
  name: "activeUser",
  initialState: activeUserState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      signUpUserThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.registerUserStatus = "PENDING";
        state.error = null;
        state.isAuthenticated = false;
        state.user = null;
      }
    );
    builder.addCase(
      signUpUserThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.registerUserStatus = "SUCCESS";
      }
    );
    builder.addCase(
      signUpUserThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.registerUserStatus = "ERROR";
        state.error = action.payload;
      }
    );

    builder.addCase(
      signInUserThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "PENDING";
        state.error = null;
        state.isAuthenticated = false;
        state.user = null;
      }
    );
    builder.addCase(
      signInUserThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "SUCCESS";
      }
    );
    builder.addCase(
      signInUserThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "ERROR";
        state.error = action.payload;
      }
    );

    builder.addCase(
      getCurrentUserThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "PENDING";
        state.error = null;
        state.isAuthenticated = false;
        state.user = null;
      }
    );
    builder.addCase(
      getCurrentUserThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "SUCCESS";
        state.isAuthenticated = true;
        state.user = action.payload;
      }
    );
    builder.addCase(
      getCurrentUserThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.loginUserStatus = "ERROR";
        state.error = action.payload;
        state.isAuthenticated = false;
        state.user = null;
      }
    );

    builder.addCase(
      logoutUserThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.logoutUserStatus = "PENDING";
        state.error = null;
      }
    );
    builder.addCase(
      logoutUserThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.logoutUserStatus = "SUCCESS";
        state.isAuthenticated = false;
        state.user = null;
      }
    );
    builder.addCase(
      logoutUserThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.logoutUserStatus = "ERROR";
      }
    );
  },
});

export const activeUserReducer = activeUserSlice.reducer;
