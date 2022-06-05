/// <reference types="telegram-web-app"/>

import * as React from "react";
import * as ReactDOM from "react-dom/client";

import "./style.css";
import App, { AppConfig } from "./app"

declare global {
    interface Window { appConfig: AppConfig; }
}

var content = <div className="p-4">
    <p>Please visit this page via Telegram.</p>
</div>;
if (window.appConfig.playerId != null) {
    content = <App config={window.appConfig}/>;
}

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(content);
