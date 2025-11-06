import React, { useState, useEffect } from "react";
import { useAppDispatch } from "../../store/hooks";
import { Spin } from "antd";
import { signInUserThunk } from "../../redux/activeUser/activeUserSlice";
import {
  useActiveUserLoginStatus,
  useActiveUserError,
} from "../../redux/selectors/selectorHooks";
import { LoginDto } from "../../models/UserModels";
import { notification } from "antd";
import { validateEmail } from "../../helpers/validateEmail";
import { pop } from "../howler/pop";
import {
  buttonTw,
  formHeadingTw,
  formLabelTw,
  inputTw,
} from "../../tailwind/tailwindClassNames";

const initialState: LoginDto = {
  email: "",
  password: "",
};

export const Login = () => {
  const dispatch = useAppDispatch();
  const userLoginStatus = useActiveUserLoginStatus();
  const error = useActiveUserError();
  const [form, setForm] = useState(initialState);

  const [api, contextHolder] = notification.useNotification();

  const activate = () => {
    if (!validateEmail(form.email)) {
      api["error"]({
        message: "Invalid Email",
        description: "Please enter a valid email address.",
      });
      pop.play();
      return;
    }
    if (!form.password) {
      api["error"]({
        message: "Wrong Password",
        description:
          "Please enter a password. Password should contain at least 8 characters.",
      });
      pop.play();
      return;
    }
    dispatch(signInUserThunk(form));
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((state) => ({
      ...state,
      [e.target.name]: e.target.value,
    }));
  };

  useEffect(() => {
    if (userLoginStatus === "ERROR" && error) {
      api["error"]({
        message: error?.message,
        description: error?.description,
      });
      pop.play();
    }
  }, [userLoginStatus]);

  return (
    <>
      {contextHolder}
      <div className="flex flex-col">
        <h1 className={formHeadingTw}>Login</h1>
        {userLoginStatus === "PENDING" ? <Spin /> : null}
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

        <label className={formLabelTw} htmlFor="password">
          Password
        </label>
        <input
          className={inputTw}
          type="password"
          name="password"
          value={form.password}
          onChange={onChange}
          required
          placeholder="Enter your password (min. 8 characters)"
        />

        <button
          className={`${buttonTw} mt-6 text-white`}
          onClick={activate}
          disabled={userLoginStatus === "PENDING"}
        >
          {userLoginStatus === "PENDING" ? "Processing..." : "Login"}
        </button>
      </div>
    </>
  );
};
