import { createSelector, Selector } from "@reduxjs/toolkit";
import { AppStore } from "../../store/store";
import { ActiveUserState } from "../activeUser/activeUserSlice";
import { ErrorMessage, Role, User } from "../../models/UserModels";
import { ApiStatus } from "../../constants";
import { selectUsersEntities } from "../users/usersSlice";
import { EventsState } from "../events/eventsSlice";
import { ParticipantsState } from "../participants/participantsSlice";

export const selectActiveUserState: Selector<
  AppStore,
  ActiveUserState | null
> = (state) => state?.activeUser ?? null;

export const selectActiveUserIsAuthenticated: Selector<AppStore, boolean> =
  createSelector(
    [selectActiveUserState],
    (state) => state?.isAuthenticated ?? false
  );

export const selectActiveUserRegisterStatus: Selector<AppStore, ApiStatus> =
  createSelector(
    [selectActiveUserState],
    (state) => state?.registerUserStatus ?? "IDLE"
  );

export const selectActiveUserId: Selector<AppStore, User["id"] | null> =
  createSelector([selectActiveUserState], (state) => state?.user ?? null);

export const selectActiveUserError: Selector<AppStore, ErrorMessage | null> =
  createSelector([selectActiveUserState], (state) => state?.error ?? null);

export const selectActiveUserLoginStatus: Selector<AppStore, ApiStatus> =
  createSelector(
    [selectActiveUserState],
    (state) => state?.loginUserStatus ?? "IDLE"
  );

export const selectIsAuthenticated: Selector<AppStore, boolean> =
  createSelector(
    [selectActiveUserState],
    (activeUserState) => activeUserState?.isAuthenticated ?? false
  );

export const selectActiveUserEntity: Selector<AppStore, User | null> =
  createSelector(
    [(i) => selectUsersEntities(i.users), selectActiveUserId],
    (users, userId) => {
      if (users && userId) {
        return users[userId];
      }
      return null;
    }
  );

export const selectActiveUserRoles: Selector<AppStore, Role[] | undefined> =
  createSelector([selectActiveUserEntity], (i) => {
    return i?.role;
  });

export const selectActiveUserIsAdmin: Selector<AppStore, boolean> =
  createSelector([selectActiveUserRoles], (i) => {
    if (Array.isArray(i)) {
      return !!i.find((i) => i.role === "admin");
    }
    return false;
  });

export const selectEventsState: Selector<AppStore, EventsState | null> = (
  state
) => state?.events ?? null;

export const selectParticipantsState: Selector<
  AppStore,
  ParticipantsState | null
> = (state) => state?.participants ?? null;

export const selectEventsError: Selector<AppStore, ErrorMessage | null> =
  createSelector([selectEventsState], (state) => state?.error ?? null);

export const selectParticipantsError: Selector<AppStore, ErrorMessage | null> =
  createSelector([selectParticipantsState], (state) => state?.error ?? null);

export const selectPostParticipantStatus: Selector<AppStore, ApiStatus | null> =
  createSelector(
    [selectParticipantsState],
    (state) => state?.postParticipantStatus ?? null
  );

export const selectGetParticipantStatus: Selector<AppStore, ApiStatus | null> =
  createSelector(
    [selectParticipantsState],
    (state) => state?.getParticipantStatus ?? null
  );

export const selectEventsPostStatus: Selector<AppStore, ApiStatus | null> =
  createSelector(
    [selectEventsState],
    (state) => state?.postEventsStatus ?? null
  );

export const selectEventsLimit: Selector<AppStore, number> = createSelector(
  [selectEventsState],
  (state) => state?.limit ?? 10
);

export const selectEventsPage: Selector<AppStore, number> = createSelector(
  [selectEventsState],
  (state) => state?.page ?? 1
);

export const selectEventSortBy: Selector<AppStore, string> = createSelector(
  [selectEventsState],
  (state) => state?.eventSortBy ?? "eventDate"
);

export const selectEventSortOrder: Selector<AppStore, string> = createSelector(
  [selectEventsState],
  (state) => state?.eventSortOrder ?? "ASC"
);

export const selectTotalCount: Selector<AppStore, number> = createSelector(
  [selectEventsState],
  (state) => state?.totalCount ?? 0
);

export const selectEventsTotalPages: Selector<AppStore, number> =
  createSelector([selectEventsState], (state) => state?.totalPages ?? 0);

export const selectInfiniteScroll: Selector<AppStore, boolean> = createSelector(
  [selectEventsState],
  (state) => state?.infiniteScroll ?? false
);
