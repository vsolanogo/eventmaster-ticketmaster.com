import React, { useState, useEffect } from "react";
import { useAppDispatch } from "../../store/hooks";
import { Spin } from "antd";
import { signUpUserThunk } from "../../redux/activeUser/activeUserSlice";
import {
  useActiveUserRegisterStatus,
  useActiveUserError,
} from "../../redux/selectors/selectorHooks";
import { RegisterDto } from "../../models/UserModels";
import { notification } from "antd";
import { validateEmail } from "../../helpers/validateEmail";
import { pop } from "../howler/pop";
import {
  buttonTw,
  formHeadingTw,
  formLabelTw,
  inputTw,
} from "../../tailwind/tailwindClassNames";

const initialState: RegisterDto = {
  email: "",
  password: "",
  confirmPassword: "",
};

export const Register = () => {
  const dispatch = useAppDispatch();
  const userRegisterStatus = useActiveUserRegisterStatus();
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
    if (form.password !== form.confirmPassword || form.password === "") {
      api["error"]({
        message: "Password Mismatch",
        description:
          "The passwords you entered do not match. Please make sure to enter the same password in both fields.",
      });
      pop.play();
      return;
    }
    if (form.password !== form.confirmPassword) {
      api["error"]({
        message: "Password Mismatch",
        description:
          "The passwords you entered do not match. Please make sure to enter the same password in both fields.",
      });
      pop.play();
      return;
    }
    dispatch(signUpUserThunk(form));
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((state) => ({
      ...state,
      [e.target.name]: e.target.value,
    }));
  };

  useEffect(() => {
    if (userRegisterStatus === "ERROR" && error) {
      api["error"]({
        message: error?.message,
        description: error?.description,
      });
      pop.play();
    }
  }, [userRegisterStatus]);

  return (
    <>
      {contextHolder}
      <div className="container py-8 mx-auto flex flex-col">
        <h1 className={formHeadingTw}>Register</h1>

        {userRegisterStatus === "PENDING" ? <Spin /> : null}
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

        <label className={formLabelTw} htmlFor="password">
          Confirm Password
        </label>

        <input
          className={inputTw}
          type="password"
          name="confirmPassword"
          value={form.confirmPassword}
          onChange={onChange}
          required
          placeholder="Confirm your password"
        />

        <button
          className={`${buttonTw} mt-6 text-white`}
          onClick={activate}
          disabled={userRegisterStatus === "PENDING"}
        >
          {userRegisterStatus === "PENDING" ? "Processing..." : "Register"}
        </button>
      </div>
    </>
  );
};
