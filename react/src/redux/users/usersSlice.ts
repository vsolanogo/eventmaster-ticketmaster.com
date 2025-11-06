import { createSlice, createEntityAdapter } from "@reduxjs/toolkit";
import { User } from "../../models/UserModels";

const usersAdapter = createEntityAdapter({
  selectId: (i: User) => i.id,
});

export const usersSlice = createSlice({
  name: "users",
  initialState: usersAdapter.getInitialState(),
  reducers: {
    addOne: usersAdapter.addOne,
    upsertOne: usersAdapter.upsertOne,
    addMany: usersAdapter.addMany,
    removeAll: usersAdapter.removeAll,
    removeOne: usersAdapter.removeOne,
  },
});

export const {
  selectIds: selectUsersIds,
  selectEntities: selectUsersEntities,
  selectAll: selectAllUsers,
  selectTotal: selectTotalUsers,
  selectById: selectUsersById,
} = usersAdapter.getSelectors();
