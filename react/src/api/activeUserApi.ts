import api from "./api";
import { LoginDto, User } from "../models/UserModels";
import {
  getCurrentUserThunk,
  signInUserThunk,
} from "../redux/activeUser/activeUserSlice";
import { usersSlice } from "../redux/users/usersSlice";
import { handleApiError } from "./handleApiError";

export const signUpUser = async (
  loginDto: LoginDto,
  thunkAPI
): Promise<void> => {
  try {
    await api.post<User>("/register", loginDto);
    thunkAPI.dispatch(signInUserThunk(loginDto));
    return;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const signInUser = async (
  loginDto: LoginDto,
  thunkAPI
): Promise<void> => {
  try {
    await api.post<User>("/login", loginDto);
    thunkAPI.dispatch(getCurrentUserThunk(undefined));
    return;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const logoutUser = async (_, thunkAPI): Promise<void> => {
  try {
    await api.post<User>("/logout", undefined, {
      withCredentials: true,
    });
    thunkAPI.dispatch(getCurrentUserThunk(undefined));
    return;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const getCurrentUser = async (_, thunkAPI) => {
  const res = await api.get<User>("/user", {
    withCredentials: true,
  });
  thunkAPI.dispatch(usersSlice.actions.upsertOne(res.data));

  return res.data.id;
};
