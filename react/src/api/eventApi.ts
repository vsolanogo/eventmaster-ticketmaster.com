import api from "./api";
import { CreateEventDto, Event, EventList } from "../models/UserModels";
import { handleApiError } from "./handleApiError";
import { eventsSlice } from "../redux/events/eventsSlice";
import {
  selectEventSortBy,
  selectEventSortOrder,
  selectEventsLimit,
  selectEventsPage,
} from "../redux/selectors/selectors";

export const getEvents = async (undefined, thunkAPI) => {
  try {
    const limit = selectEventsLimit(thunkAPI.getState());
    const page = selectEventsPage(thunkAPI.getState());
    const sortOrder = selectEventSortOrder(thunkAPI.getState());
    const sortBy = selectEventSortBy(thunkAPI.getState());

    const params = { limit, page, sortBy, sortOrder };
    const urlSearchParams = new URLSearchParams(params as any);

    const res = await api.get<EventList>(`/events?${urlSearchParams}`);
    thunkAPI.dispatch(eventsSlice.actions.removeAll());
    thunkAPI.dispatch(eventsSlice.actions.upsertMany(res.data.events));
    return res.data;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const createEvent = async (event: CreateEventDto, thunkAPI) => {
  try {
    const res = await api.post<Event>("/events", event);
    thunkAPI.dispatch(eventsSlice.actions.addOne(res.data));
    return res.data;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const getEventById = async (id: string, thunkAPI) => {
  try {
    const res = await api.get<Event>(`/events/${id}`);
    thunkAPI.dispatch(eventsSlice.actions.upsertOne(res.data));
    return;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

// export const deleteEvent = (id: string) => {
//   return api.delete(`/events/${id}`);
// };

// export const patchEvents = (user: Partial<Event>) => {
//   return api.patch<Event>(`/event/${user.id}`, user).then((res) => res.data);
// };
