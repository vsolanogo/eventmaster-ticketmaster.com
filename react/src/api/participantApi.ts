import api from "./api";
import {
  Event,
  RegisterParticipantDto,
  Participant,
} from "../models/UserModels";
import { handleApiError } from "./handleApiError";
import { participantsSlice } from "../redux/participants/participantsSlice";

export const getParticipants = async (eventId, thunkAPI) => {
  try {
    const res = await api.get<Participant[]>(`/participant/event/${eventId}`);
    thunkAPI.dispatch(participantsSlice.actions.removeAll());
    thunkAPI.dispatch(participantsSlice.actions.upsertMany(res.data));
    return res.data;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};

export const createParticipant = async (
  {
    eventId,
    participant,
  }: { eventId: string; participant: RegisterParticipantDto },
  thunkAPI
) => {
  try {
    const res = await api.post<Event>(
      `/participant/event/${eventId}`,
      participant
    );
    return res.data;
  } catch (e: any) {
    return handleApiError(e, thunkAPI);
  }
};
