import { Action } from "redux";
import { ThunkAction } from "redux-thunk";
import { configureStore } from "@reduxjs/toolkit";
import {
  ActiveUserState,
  activeUserReducer,
} from "../redux/activeUser/activeUserSlice";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { usersSlice } from "../redux/users/usersSlice";
import { eventsSlice } from "../redux/events/eventsSlice";
import { participantsSlice } from "../redux/participants/participantsSlice";

export type AppThunkAction<R = void> = ThunkAction<
  R,
  AppStore,
  undefined,
  Action
>;

export interface AppStore {
  activeUser: ActiveUserState;
  users: ReturnType<typeof usersSlice.reducer>;
  events: ReturnType<typeof eventsSlice.reducer>;
  participants: ReturnType<typeof participantsSlice.reducer>;
}

export const store = configureStore({
  reducer: {
    activeUser: activeUserReducer,
    users: usersSlice.reducer,
    events: eventsSlice.reducer,
    participants: participantsSlice.reducer,
  },
  devTools: true,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
