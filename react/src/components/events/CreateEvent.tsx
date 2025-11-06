import React, { useEffect, useState } from "react";
import {
  useActiveUser,
  useActiveUserIsAdmin,
  useEventsError,
  useEventsPostStatus,
} from "../../redux/selectors/selectorHooks";
import { pop } from "../howler/pop";
import { useLocation } from "wouter";
import { Upload, Button, notification, DatePicker } from "antd";
import { UploadOutlined } from "@ant-design/icons";
import type { UploadFile } from "antd";
import { axiosParams } from "../../api/api";
import { CreateEventDto } from "../../models/UserModels";
import {
  buttonTw,
  formLabelTw,
  inputTw,
} from "../../tailwind/tailwindClassNames";
import { useAppDispatch } from "../../store/store";
import { postEventThunk } from "../../redux/events/eventsSlice";

const eventState: CreateEventDto = {
  title: "",
  description: "",
  latitude: 0.0,
  longitude: 0.0,
  images: [],
  eventDate: null,
};

const floats = ["latitude", "longitude"];

export const CreateEvent = () => {
  const dispatch = useAppDispatch();
  const isAdmin = useActiveUserIsAdmin();
  const userId = useActiveUser();
  const [, navigate] = useLocation();
  const [api, contextHolder] = notification.useNotification();
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [form, setForm] = useState<CreateEventDto>(eventState);
  const error = useEventsError();
  const eventsPostStatus = useEventsPostStatus();

  const handleImgChange = (info: { fileList: UploadFile[] }) => {
    setFileList(info.fileList);
  };

  useEffect(() => {
    if (!isAdmin && userId) {
      pop.play();
      api["error"]({
        message: "Access Denied",
        description: "You do not have permission to access this page.",
      });
    }
  }, [isAdmin, userId]);

  const inputOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = floats.includes(e.target.name)
      ? parseFloat(e.target.value)
      : e.target.value;

    setForm((state) => ({
      ...state,
      [e.target.name]: value,
    }));
  };

  useEffect(() => {
    if (eventsPostStatus === "ERROR" && error) {
      api["error"]({
        message: error?.message,
        description: error?.description,
      });
      pop.play();
    }
  }, [eventsPostStatus]);

  const activate = () => {
    const prep = {
      ...form,
      images: fileList.map((i) => i?.response?.id),
    };

    dispatch(postEventThunk(prep))
      .unwrap()
      .then((i) => {
        navigate(`/events/${i?.id}`);
      });
  };

  const onSelectDate = (date) => {
    setForm((state) => ({
      ...state,
      eventDate: date ? date.toDate() : null,
    }));
  };

  return (
    <>
      {contextHolder}

      {isAdmin && (
        <>
          <div className="container py-8 mx-auto flex flex-col">
            <label className={formLabelTw} htmlFor="title">
              Title
            </label>

            <input
              value={form.title}
              onChange={inputOnChange}
              placeholder="Title"
              name="title"
              className={inputTw}
              type="text"
            />

            <label className={formLabelTw} htmlFor="description">
              Description
            </label>

            <input
              value={form.description}
              onChange={inputOnChange}
              placeholder="Description"
              name="description"
              className={inputTw}
              type="text"
            />

            <label className={formLabelTw} htmlFor="latitude">
              Latitude
            </label>

            <input
              value={form.latitude}
              onChange={inputOnChange}
              placeholder="Latitude"
              name="latitude"
              className={inputTw}
              type="number"
              min={-90}
              max={90}
              pattern="[0-9]*"
            />

            <label className={formLabelTw} htmlFor="longitude">
              Longitude
            </label>

            <input
              value={form.longitude}
              onChange={inputOnChange}
              placeholder="Longitude"
              name="longitude"
              className={inputTw}
              type="number"
              min={-180}
              max={180}
              pattern="[0-9]*"
            />

            <label className={formLabelTw}>Event Date</label>

            <DatePicker
              className={`${inputTw}`}
              onChange={onSelectDate}
              placeholder="Event Date"
            />

            <label className={formLabelTw}>Images</label>

            <Upload
              action={`${axiosParams.baseURL}/image`}
              listType="picture"
              onChange={handleImgChange}
              accept=".png,.jpg,.jpeg"
            >
              <Button icon={<UploadOutlined />}>Upload</Button>
            </Upload>

            <button
              className={`${buttonTw} mt-6 text-white`}
              onClick={activate}
              // disabled={userRegisterStatus === "PENDING"}
            >
              Create
            </button>
          </div>
        </>
      )}
    </>
  );
};
