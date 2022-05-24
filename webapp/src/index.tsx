/// <reference types="telegram-web-app" />

import "./style.css";

import * as React from "react";
import * as ReactDOM from "react-dom/client";
import interact from "interactjs";

const content = <div>
    <h1>Welcome to the web app, {window.Telegram.WebApp.initDataUnsafe.user.username}!</h1>
    <h2 className="hint">Have a nice day!</h2>
    <pre>{JSON.stringify(window.Telegram.WebApp, null, 2)}</pre>
</div>;

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(content);
