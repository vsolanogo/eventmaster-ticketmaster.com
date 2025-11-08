import React from "react";
import { Redirect, Route, Switch } from "wouter";
import { Register } from "./register/Register";
import { Home } from "./home/Home";
import { Login } from "./login/Login";
import { CreateEvent } from "./events/CreateEvent";
import { EventDetails } from "./events/EventDetails";
import { RegisterParticipant } from "./events/RegisterParticipant";

export const Router = () => {
  return (
    <>
      <Switch>
        <Route path="/" component={Home} />
        <Route path="/login" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/create-event" component={CreateEvent} />
        <Route path="/events/:id" component={EventDetails} />
        <Route path="/events/register/:id" component={RegisterParticipant} />
        <Route path="*" component={() => <Redirect to="/" />} />
      </Switch>
    </>
  );
};
