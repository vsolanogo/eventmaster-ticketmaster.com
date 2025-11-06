import React, { useEffect } from "react";
import { Router } from "./Router";
import { Header } from "./common/Header";
import { useAppDispatch } from "../store/store";
import { getCurrentUserThunk } from "../redux/activeUser/activeUserSlice";
import { Footer } from "./footer/Footer";

export const App = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(getCurrentUserThunk(undefined));
  }, []);

  return (
    <>
      <Header />
      <div className="App mx-auto max-w-6xl text-center my-8">
        <Router />
      </div>
      <Footer />
    </>
  );
};
