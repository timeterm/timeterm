import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import "@rmwc/drawer/styles";
import "@rmwc/list/styles";
import "@rmwc/button/styles";
import "@rmwc/checkbox/styles";
import "@rmwc/icon-button/styles";
import "@rmwc/icon/styles";
import "@rmwc/elevation/styles";
import "@rmwc/theme/styles";
import "@rmwc/data-table/styles";
import "@rmwc/select/styles";

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
