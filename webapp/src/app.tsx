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
        this.ws.addEventListener("open", this.onOpenConnection.bind(this));
        this.ws.addEventListener("close", this.onCloseConnection.bind(this));
        this.ws.addEventListener("error", this.onError.bind(this));
        this.ws.addEventListener("message", this.onMessage.bind(this));

        this.handlers = {
            moveCard: this.moveCard.bind(this)
        };
    }

    onOpenConnection(e: Event) {
        this.ws.send(JSON.stringify({
            "messageType": "authorize",
            "payload": {
                "gameId": this.props.config.gameId,
                "playerId": this.props.config.playerId
            }
        }));
    }

    onCloseConnection(e: CloseEvent) {
        window.Telegram.WebApp.close();
    }

    onError(e: Event) {
        window.Telegram.WebApp.close();
    }

    onMessage(e: MessageEvent) {
        var message = JSON.parse(e.data);

        if (message.messageType == "joinedAsPlayer") {
            window.Telegram.WebApp.close();
            return
        }   
    }

    moveCard(cardId: CardID, deckId: DeckID) {
        this.ws.send(JSON.stringify({cardId: cardId, deckId: deckId}));
    }

    render(): JSX.Element {
        return <PreferansView gameState={this.state.gameState} handlers={this.handlers}/>;
    }
}

export default App;
