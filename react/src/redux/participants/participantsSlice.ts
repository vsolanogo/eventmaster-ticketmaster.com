import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  PayloadAction,
} from "@reduxjs/toolkit";
import { ErrorMessage, Participant } from "../../models/UserModels";
import { ApiStatus } from "../../constants";
import { createParticipant, getParticipants } from "../../api/participantApi";

export type ParticipantsState = {
  postParticipantStatus: ApiStatus;
  getParticipantStatus: ApiStatus;
  error: ErrorMessage | null;
};

export const participantsState: ParticipantsState = {
  postParticipantStatus: "IDLE",
  getParticipantStatus: "IDLE",
  error: null,
};

const participantsAdapter = createEntityAdapter({
  selectId: (i: Participant) => i.id,
});

export const getEventParticipantsThunk = createAsyncThunk(
  "participants/get",
  getParticipants
);

export const createParticipantThunk = createAsyncThunk(
  "participants/create",
  createParticipant
);

export const participantsSlice = createSlice({
  name: "participants",
  initialState: {
    ...participantsState,
    ...participantsAdapter.getInitialState(),
  },
  reducers: {
    addOne: participantsAdapter.addOne,
    upsertOne: participantsAdapter.upsertOne,
    upsertMany: participantsAdapter.upsertMany,
    addMany: participantsAdapter.addMany,
    removeAll: participantsAdapter.removeAll,
    removeOne: participantsAdapter.removeOne,
  },
  extraReducers: (builder) => {
    builder.addCase(
      createParticipantThunk.pending,
      (state, action: PayloadAction<any>) => {
        state.postParticipantStatus = "PENDING";
        state.error = null;
      }
    );
    builder.addCase(
      createParticipantThunk.fulfilled,
      (state, action: PayloadAction<any>) => {
        state.postParticipantStatus = "SUCCESS";
      }
    );
    builder.addCase(
      createParticipantThunk.rejected,
      (state, action: PayloadAction<any>) => {
        state.postParticipantStatus = "ERROR";
        state.error = action.payload;
      }
    );
  },
});

export const {
  selectIds: selectParticipantsIds,
  selectEntities: selectParticipantsEntities,
  selectAll: selectAllParticipants,
  selectTotal: selectTotalParticipants,
  selectById: selectParticipantsById,
} = participantsAdapter.getSelectors();
