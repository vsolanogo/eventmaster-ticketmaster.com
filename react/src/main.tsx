import React from "react";
import { createRoot } from "react-dom/client";
import "./App.css";
import { Provider } from "react-redux";
import { App } from "./components/App";
import { store } from "./store/store";

const root = document.getElementById("root");

const app = (
  <Provider store={store}>
    <App />
  </Provider>
);

if (root) {
  const rootElement = createRoot(root);
  rootElement.render(app);
} else {
  console.error("Root element not found in the document.");
}
