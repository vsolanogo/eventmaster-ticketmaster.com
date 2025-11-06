import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  PayloadAction,
  createAction,
} from "@reduxjs/toolkit";
import { Event, ErrorMessage } from "../../models/UserModels";
import { createEvent, getEvents, getEventById } from "../../api/eventApi";
import { ApiStatus } from "../../constants";
import { SortOrderType } from "../../helpers/types/types";

export type EventsState = {
  postEventsStatus: ApiStatus;
  getEventsStatus: ApiStatus;
  error: ErrorMessage | null;
  totalCount: number;
  limit: number;
  page: number;
  totalPages: number;
  eventDateSort: SortOrderType;
  eventSortBy: string;
  eventSortOrder: string;
  infiniteScroll: boolean;
};

export const eventsState: EventsState = {
  postEventsStatus: "IDLE",
  getEventsStatus: "IDLE",
  error: null,
  totalCount: 0,
  limit: 10,
  page: 1,
  totalPages: 0,
  eventDateSort: "ASC",
  eventSortBy: "eventDate",
  eventSortOrder: "ASC",
  infiniteScroll: false,
};

const eventsAdapter = createEntityAdapter({
  selectId: (event: Event) => event.id,
});
export const getEventsThunk = createAsyncThunk("events/list", getEvents);
export const getEventByIdThunk = createAsyncThunk(
  "events/getById",
  getEventById
);
export const postEventThunk = createAsyncThunk("events/create", createEvent);
export const setPage = createAction<number>("events/setPage");
export const setEventsLimit = createAction<number>("events/setEventsLimit");
export const setSortBy = createAction<string>("events/setSortBy");
export const setSortOrder = createAction<string>("events/setSortOrder");
export const setInfiniteScroll = createAction<boolean>(
  "events/setInfiniteScroll"
);

export const eventsSlice = createSlice({
  name: "events",
  initialState: { ...eventsState, ...eventsAdapter.getInitialState() },
  reducers: {
    addOne: eventsAdapter.addOne,
    upsertOne: eventsAdapter.upsertOne,
    upsertMany: eventsAdapter.upsertMany,
    addMany: eventsAdapter.addMany,
    removeAll: eventsAdapter.removeAll,
    removeOne: eventsAdapter.removeOne,
  },
  extraReducers: (builder) => {
    builder.addCase(
      setInfiniteScroll,
      (state, action: PayloadAction<boolean>) => {
        state.infiniteScroll = action.payload;
      }
    );
    builder.addCase(setSortBy, (state, action: PayloadAction<string>) => {
      state.eventSortBy = action.payload;
    });
    builder.addCase(setSortOrder, (state, action: PayloadAction<string>) => {
      state.eventSortOrder = action.payload;
    });

    builder.addCase(setPage, (state, action: PayloadAction<number>) => {
      state.page = action.payload;
    });

    builder.addCase(setEventsLimit, (state, action: PayloadAction<number>) => {
      state.limit = action.payload;
    });

    builder.addCase(
      postEventThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.postEventsStatus = "PENDING";
        state.error = null;
      }
    );
    builder.addCase(
      postEventThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.postEventsStatus = "SUCCESS";
      }
    );
    builder.addCase(
      postEventThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.postEventsStatus = "ERROR";
        state.error = action.payload;
      }
    );

    builder.addCase(
      getEventsThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.getEventsStatus = "PENDING";
        state.error = null;
      }
    );
    builder.addCase(
      getEventsThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.totalCount = action.payload.totalCount;
        const numberOfPages = Math.ceil(
          action.payload.totalCount / state.limit
        );
        state.totalPages = numberOfPages;
        state.getEventsStatus = "SUCCESS";
      }
    );
    builder.addCase(
      getEventsThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.getEventsStatus = "ERROR";
        state.error = action.payload;
      }
    );
  },
});

export const {
  selectIds: selectEventsIds,
  selectEntities: selectEventsEntities,
  selectAll: selectAllEvents,
  selectTotal: selectTotalEvents,
  selectById: selectEventsById,
} = eventsAdapter.getSelectors();
