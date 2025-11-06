import React, { useEffect } from "react";
import {
  useActiveUserIsAdmin,
  useIsAuthenticated,
} from "../../redux/selectors/selectorHooks";
import { buttonTw } from "../../tailwind/tailwindClassNames";
import { Link, useRoute, useLocation } from "wouter";
import { useAppDispatch } from "../../store/store";
import { logoutUserThunk } from "../../redux/activeUser/activeUserSlice";
import { Button, Modal, notification } from "antd";
import { pop } from "../howler/pop";
import eventImage from "../../assets/img/event.png";

const buttonHeaderTw = `${buttonTw} mt-0 ml-2 min-w-40 inline-block text-center text-white`;

export const Header = () => {
  const dispatch = useAppDispatch();
  const isAdmin = useActiveUserIsAdmin();
  const [api, contextHolder] = notification.useNotification();
  const isAuthenticated = useIsAuthenticated();
  const [matchLogin] = useRoute("/login");
  const [matchRegister] = useRoute("/register");
  const [, navigate] = useLocation();

  useEffect(() => {
    if ((matchLogin && isAuthenticated) || (matchRegister && isAuthenticated)) {
      navigate("/");
    }
  }, [isAuthenticated, matchLogin, matchRegister]);

  const signOut = () => {
    dispatch(logoutUserThunk(undefined))
      .unwrap()
      .catch((e) => {
        pop.play();
        api["error"]({
          message: e?.message,
          description: e?.description,
        });
      });
  };

  return (
    <>
      {contextHolder}

      <header className="shadow-md sticky top-0 z-50 bg-[#172234] min-h-20 flex">
        <div className="container mx-auto flex justify-center self-center md:justify-end">
          <Link href="/" className="mr-auto">
            <img src={eventImage} alt="Event" className="h-12 w-auto" />
          </Link>

          {!isAuthenticated ? (
            <>
              <Link
                className={`${buttonTw} mt-0 border border-solid border-2 border-[#B29F7E] ml-2 min-w-40 inline-block text-center bg-transparent text-[#B29F7E]`}
                href="/login"
              >
                Log In
              </Link>
              <Link className={buttonHeaderTw} href="/register">
                Sign Up
              </Link>
            </>
          ) : (
            <>
              {isAdmin && (
                <Link className={buttonHeaderTw} href="/create-event">
                  Create Event
                </Link>
              )}
              <Button
                type="primary"
                className={buttonHeaderTw}
                onClick={() => {
                  Modal.confirm({
                    onOk: () => {
                      signOut();
                    },
                    title: "Confirm Logout",
                    content: "Are you sure you want to logout?",
                    footer: (_, { OkBtn, CancelBtn }) => (
                      <>
                        <CancelBtn />
                        <OkBtn />
                      </>
                    ),
                  });
                }}
              >
                Sign Out
              </Button>
            </>
          )}
        </div>
      </header>

      {/* <div className="flex flex-col h-screen">
        <div className="flex-1 bg-red-500"></div>
        <div className="flex-1 bg-blue-500"></div>
        <div className="flex-1 bg-green-500"></div>
      </div> */}
    </>
  );
};
