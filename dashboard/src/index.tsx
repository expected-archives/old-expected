import React from "react";
import ReactDOM from "react-dom";
import { createBrowserHistory } from "history";
import { Router } from "react-router";
import App from "./App";
import "./styles/index.scss";

const history = createBrowserHistory();

ReactDOM.render((
    <Router history={history}>
        <App />
    </Router>
), document.getElementById("root"));
