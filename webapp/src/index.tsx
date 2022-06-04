/// <reference types="telegram-web-app"/>

import "./style.css";

import * as React from "react";
import * as ReactDOM from "react-dom/client";
import {DndContext, DragEndEvent, useDraggable, useDroppable} from "@dnd-kit/core"
import {CSS} from "@dnd-kit/utilities";

type CardID = string | number;
type DeckID = string | number;

function Card({cardId}: {cardId: CardID}) {
    const {attributes, listeners, setNodeRef, transform} = useDraggable({id: cardId});
    const style = {transform: CSS.Translate.toString(transform)};

    return (
        <div 
            className="basis-8 bg-white p-1 rounded-md text-center text-black"
            ref={setNodeRef} style={style} {...listeners} {...attributes}
        >
            <strong>{cardId}</strong>
        </div>
    );
}

function CardDeck({deckId, cardIds}: {deckId: DeckID, cardIds: CardID[]}) {
    const {setNodeRef} = useDroppable({id: deckId});

    return (
        <div className="flex flex-row space-x-1" ref={setNodeRef}>
            {cardIds.map((cardId: CardID) => {
                return <Card key={cardId} cardId={cardId}/>
            })}
        </div>
    );
}

interface PreferansState {
}

interface PreferansViewHandlers {
    moveCard: (cardId: CardID, deckId: DeckID) => void;
}

interface PreferansViewProps {
    gameState: PreferansState;
    handlers: PreferansViewHandlers;
}

class PreferansView extends React.Component<PreferansViewProps, {}> {
    constructor(props: PreferansViewProps) {
        super(props);
        this.onCardDragEnd = this.onCardDragEnd.bind(this);
    }

    onCardDragEnd(e: DragEndEvent) {
        if (e.over == null) {
            return
        }
        this.props.handlers.moveCard(e.active.id, e.over.id);
    }

    render(): JSX.Element {
        return (
            <div className="preferans-view container h-screen p-4 space-y-2 w-screen">
                <DndContext onDragEnd={this.onCardDragEnd}>
                    <CardDeck deckId={1} cardIds={[1, 2]}/>
                    <CardDeck deckId={2} cardIds={[3, 4]}/>
                    <CardDeck deckId={3} cardIds={[5, 6]}/>
                </DndContext>
            </div>
        );
    }
}

interface AppConfig {
    webSocketUrl: string;
    gameId: string;
    playerId: string | undefined | null;
}

interface AppState {
    gameState: PreferansState;
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
