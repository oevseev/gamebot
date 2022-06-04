/// <reference types="telegram-web-app"/>

import "./style.css";

import * as React from "react";
import * as ReactDOM from "react-dom/client";

interface PreferansState {
}

interface PreferansViewHandlers {
}

interface PreferansViewProps {
    gameState: PreferansState;
    handlers: PreferansViewHandlers;
}

class PreferansView extends React.Component<PreferansViewProps, {}> {
    render(): JSX.Element {
        return <div>{JSON.stringify(this.props.gameState)}</div>;
    }
}

interface AppConfig {
}

interface AppState {
    gameState: PreferansState;
}

class App extends React.Component<{config: AppConfig}, AppState> {
    handlers: PreferansViewHandlers;

    constructor(props: {config: AppConfig}) {
        super(props);
        this.state = {
            gameState: {}
        };
        this.handlers = {};
    }

    render(): JSX.Element {
        return <PreferansView gameState={this.state.gameState} handlers={this.handlers}/>;
    }
}

const appConfig = {};

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<App config={appConfig}/>);
