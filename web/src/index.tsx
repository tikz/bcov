import React from "react";
import ReactDOM from "react-dom/client";
import App from "./components/App";

import LogRocket from "logrocket";
if (process.env.NODE_ENV !== "development") {
  LogRocket.init("wuajdc/bcov");
}

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
