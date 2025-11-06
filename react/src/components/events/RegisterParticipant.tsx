import React, { useState, useEffect } from "react";
import { useAppDispatch } from "../../store/hooks";
import { DatePicker, Spin } from "antd";
import {
  useParticipantsError,
  usePostParticipantStatus,
} from "../../redux/selectors/selectorHooks";
import { RegisterParticipantDto } from "../../models/UserModels";
import { notification } from "antd";
import { validateEmail } from "../../helpers/validateEmail";
import { pop } from "../howler/pop";
import {
  buttonTw,
  formHeadingTw,
  formLabelTw,
  inputTw,
} from "../../tailwind/tailwindClassNames";
import { createParticipantThunk } from "../../redux/participants/participantsSlice";
import { useLocation, useRoute } from "wouter";

const initialState: RegisterParticipantDto = {
  fullName: "",
  email: "",
  dateOfBirth: null,
  sourceOfEventDiscovery: "Social media",
};

export const RegisterParticipant = () => {
  const dispatch = useAppDispatch();
  const postParticipantStatus = usePostParticipantStatus();
  const error = useParticipantsError();
  const [form, setForm] = useState<RegisterParticipantDto>(initialState);
  const [, params] = useRoute("/events/register/:id");
  const [api, contextHolder] = notification.useNotification();
  const [, navigate] = useLocation();

  const activate = () => {
    if (!validateEmail(form.email)) {
      api["error"]({
        message: "Invalid Email",
        description: "Please enter a valid email address.",
      });
      pop.play();
      return;
    }
    dispatch(
      createParticipantThunk({ eventId: params?.id || "", participant: form })
    )
      .unwrap()
      .then(() => {
        navigate(`/events/${params?.id}`);
      });
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((state) => ({
      ...state,
      [e.target.name]: e.target.value,
    }));
  };

  useEffect(() => {
    if (postParticipantStatus === "ERROR" && error) {
      api["error"]({
        message: error?.message,
        description: error?.description,
      });
      pop.play();
    }
  }, [postParticipantStatus]);

  const onSourceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((state) => ({
      ...state,
      ["sourceOfEventDiscovery"]: e.target.value,
    }));
  };

  const onSelectDate = (date) => {
    setForm((state) => ({
      ...state,
      dateOfBirth: date ? date.toDate() : null,
    }));
  };

  return (
    <>
      {contextHolder}
      <div className="container py-8 mx-auto flex flex-col">
        <h1 className={formHeadingTw}>Register for event</h1>

        {postParticipantStatus === "PENDING" ? <Spin /> : null}

        <label className={formLabelTw} htmlFor="fullName">
          Full Name
        </label>
        <input
          className={inputTw}
          type="text"
          name="fullName"
          value={form.fullName}
          onChange={onChange}
          required
          placeholder="Enter your Full Name"
        />

        <label className={formLabelTw} htmlFor="email">
          Email
        </label>

        <input
          className={inputTw}
          type="text"
          name="email"
          value={form.email}
          onChange={onChange}
          placeholder="Enter your email address"
          required
        />

        <label className={formLabelTw}>Date of birth</label>

        <DatePicker
          className={`${inputTw}`}
          onChange={onSelectDate}
          placeholder="Birth Date"
          disabledDate={(currentDate) => currentDate.isAfter(new Date())}
        />

        <label className={formLabelTw} htmlFor="role">
          Where did you hear about this event?
        </label>
        <div className="mt-2">
          <label className="inline-flex items-center">
            <input
              type="radio"
              className="form-radio text-indigo-600"
              name="role"
              value="Social media"
              checked={form.sourceOfEventDiscovery === "Social media"}
              onChange={onSourceChange}
            />
            <span className="ml-2">Social media</span>
          </label>
          <label className="inline-flex items-center ml-6">
            <input
              type="radio"
              className="form-radio text-indigo-600"
              name="role"
              value="Friends"
              checked={form.sourceOfEventDiscovery === "Friends"}
              onChange={onSourceChange}
            />
            <span className="ml-2">Friends</span>
          </label>
          <label className="inline-flex items-center ml-6">
            <input
              type="radio"
              className="form-radio text-indigo-600"
              name="role"
              value="Found myself"
              checked={form.sourceOfEventDiscovery === "Found myself"}
              onChange={onSourceChange}
            />
            <span className="ml-2">Found myself</span>
          </label>
        </div>

        <button
          className={`${buttonTw} mt-6 text-white`}
          onClick={activate}
          disabled={postParticipantStatus === "PENDING"}
        >
          {postParticipantStatus === "PENDING" ? "Processing..." : "Register"}
        </button>
      </div>
    </>
  );
};
