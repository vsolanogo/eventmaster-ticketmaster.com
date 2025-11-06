import { useAppSelector } from "../../store/store";
import {
  selectActiveUserIsAuthenticated,
  selectActiveUserRegisterStatus,
  selectActiveUserId,
  selectActiveUserError,
  selectActiveUserLoginStatus,
  selectIsAuthenticated,
  selectActiveUserIsAdmin,
  selectEventsError,
  selectEventsPostStatus,
  selectEventsLimit,
  selectEventsPage,
  selectEventsTotalPages,
  selectTotalCount,
  selectParticipantsError,
  selectPostParticipantStatus,
  selectGetParticipantStatus,
  selectEventSortBy,
  selectEventSortOrder,
  selectInfiniteScroll,
} from "./selectors";
import { ApiStatus } from "../../constants";
import { ErrorMessage, User } from "../../models/UserModels";
import {
  selectAllEvents,
  selectEventsById,
  selectEventsEntities,
  selectEventsIds,
  selectTotalEvents,
} from "../events/eventsSlice";
import {
  selectAllParticipants,
  selectParticipantsById,
  selectParticipantsEntities,
  selectParticipantsIds,
  selectTotalParticipants,
} from "../participants/participantsSlice";

export const useActiveUserIsAuthenticated = (): boolean =>
  useAppSelector<boolean>(selectActiveUserIsAuthenticated);

export const useActiveUserRegisterStatus = (): ApiStatus =>
  useAppSelector<ApiStatus>(selectActiveUserRegisterStatus);

export const useActiveUser = (): User["id"] | null =>
  useAppSelector<User["id"] | null>(selectActiveUserId);

export const useActiveUserError = (): ErrorMessage | null =>
  useAppSelector<ErrorMessage | null>(selectActiveUserError);

export const useActiveUserLoginStatus = (): ApiStatus =>
  useAppSelector<ApiStatus>(selectActiveUserLoginStatus);

export const useIsAuthenticated = (): boolean =>
  useAppSelector<boolean>(selectIsAuthenticated);

export const useActiveUserIsAdmin = (): boolean =>
  useAppSelector<boolean>(selectActiveUserIsAdmin);

export const useEventsError = (): ErrorMessage | null =>
  useAppSelector<ErrorMessage | null>(selectEventsError);

export const useParticipantsError = (): ErrorMessage | null =>
  useAppSelector<ErrorMessage | null>(selectParticipantsError);

export const useEventsPostStatus = (): ApiStatus | null =>
  useAppSelector<ApiStatus | null>(selectEventsPostStatus);

export const useEventsLimit = (): number =>
  useAppSelector<number>(selectEventsLimit);

export const useEventsPage = (): number =>
  useAppSelector<number>(selectEventsPage);

export const useEventsTotalPages = (): ApiStatus | number =>
  useAppSelector<ApiStatus | number>(selectEventsTotalPages);

export const useAllEvents = () =>
  useAppSelector((state) => selectAllEvents(state.events));

export const useEventsIds = () =>
  useAppSelector((state) => selectEventsIds(state.events));

export const useEventsEntities = () =>
  useAppSelector((state) => selectEventsEntities(state.events));

export const useTotalEvents = () =>
  useAppSelector((state) => selectTotalEvents(state.events));

export const useEventsById = (id: string) =>
  useAppSelector((state) => selectEventsById(state.events, id));

export const useTotalCount = () =>
  useAppSelector((state) => selectTotalCount(state));

export const usePostParticipantStatus = (): ApiStatus | null =>
  useAppSelector<ApiStatus | null>(selectPostParticipantStatus);

export const useGetParticipantStatus = (): ApiStatus | null =>
  useAppSelector<ApiStatus | null>(selectGetParticipantStatus);

export const useParticipantsIds = () =>
  useAppSelector((state) => selectParticipantsIds(state.participants));
export const useParticipantsEntities = () =>
  useAppSelector((state) => selectParticipantsEntities(state.participants));
export const useAllParticipants = () =>
  useAppSelector((state) => selectAllParticipants(state.participants));
export const useTotalParticipants = () =>
  useAppSelector((state) => selectTotalParticipants(state.participants));
export const useParticipantsById = (id: string) =>
  useAppSelector((state) => selectParticipantsById(state.participants, id));
export const useEventSortBy = (): string =>
  useAppSelector<string>(selectEventSortBy);
export const useEventSortOrder = (): string =>
  useAppSelector<string>(selectEventSortOrder);
export const useInfiniteScroll = (): boolean =>
  useAppSelector<boolean>(selectInfiniteScroll);
