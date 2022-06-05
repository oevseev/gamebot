import * as React from "react";
import PreferansView, { CardID, DeckID, PreferansState, PreferansViewHandlers } from "./games/preferans";

interface AppState {
    gameState: PreferansState;
}

export interface AppConfig {
    webSocketUrl: string;
    gameId: string;
    playerId: string | undefined | null;
}

class App extends React.Component<{config: AppConfig}, AppState> {
    ws: WebSocket;
    handlers: PreferansViewHandlers;

    constructor(props: {config: AppConfig}) {
        super(props);

        this.state = {
            gameState: {}
        };

        this.ws = new WebSocket(this.props.config.webSocketUrl);
        this.ws.addEventListener("open", this.openConnection.bind(this));

        this.handlers = {
            moveCard: this.moveCard.bind(this)
        };
    }

    openConnection(e: Event) {
        this.ws.send(JSON.stringify(this.props.config));
    }

    moveCard(cardId: CardID, deckId: DeckID) {
        this.ws.send(JSON.stringify({cardId: cardId, deckId: deckId}));
    }

    render(): JSX.Element {
        return <PreferansView gameState={this.state.gameState} handlers={this.handlers}/>;
    }
}

export default App;
