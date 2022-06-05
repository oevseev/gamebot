import * as React from "react";
import {DndContext, DragEndEvent, useDraggable, useDroppable} from "@dnd-kit/core"
import {CSS} from "@dnd-kit/utilities";

export type CardID = string | number;
export type DeckID = string | number;

const CARD_VALUES = ["", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"];
const CARD_SUITS = ["", "♣️", "♦️", "♥️", "♠️"];

// Spades, diamonds, clubs, hearts (alternate red and black suits for easier distinguishing)
const SUIT_ORDER = [0, 3, 2, 4, 1];

function CardDenomination({cardId}: {cardId: CardID}): JSX.Element {
    if (typeof cardId !== "string") {
        return <p></p>;
    }
    const index = cardId.indexOf(",");
    const rank = CARD_VALUES[parseInt(cardId.substring(0, index))];
    const suit = CARD_SUITS[parseInt(cardId.substring(index + 1))];

    return <p><strong className="px-1">{rank}</strong><br/>{suit}</p>;
}

function compareCards(a: CardID, b: CardID): number {
    if (typeof a === "number" || typeof b === "number") {
        throw new Error("compareCards can compare strings only");
    }

    const indexA = a.indexOf(",");
    const indexB = b.indexOf(",");

    const rankA = parseInt(a.substring(0, indexA));
    const rankB = parseInt(b.substring(0, indexB));

    const suitA = SUIT_ORDER[parseInt(a.substring(indexA + 1))];
    const suitB = SUIT_ORDER[parseInt(b.substring(indexB + 1))];

    if (suitA < suitB) {
        return -1;
    }
    if (suitA > suitB) {
        return 1;
    }
    if (rankA < rankB) {
        return 1;
    }
    if (rankA > rankB) {
        return -1;
    }
    return 0;
}

function Card({cardId}: {cardId: CardID}): JSX.Element {
    const {attributes, isDragging, listeners, setNodeRef, transform} = useDraggable({id: cardId});
    const style = {transform: CSS.Translate.toString(transform)};

    return (
        <div
            className={`card ${isDragging ? "dragged" : ""}`}
            ref={setNodeRef} style={style} {...listeners} {...attributes}
        >
            <CardDenomination cardId={cardId}/>
        </div>
    );
}

function CardDeck({deckId, cardIds}: {deckId: DeckID, cardIds: CardID[]}): JSX.Element {
    const {setNodeRef} = useDroppable({id: deckId});

    return (
        <div className="card-deck" ref={setNodeRef}>
            {cardIds.map((cardId: CardID) => {
                return <Card key={cardId} cardId={cardId}/>
            })}
        </div>
    );
}

export interface PreferansConfig {
}

export interface PreferansState {
    table: {
        owner: {
            [cardId: string]: number;
        }
    }
}

export interface PreferansViewHandlers {
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
        var cards: CardID[] = [];
        if (this.props.gameState != null) {
            for (const [card, owner] of Object.entries(this.props.gameState.table.owner)) {
                if (owner != 1) {
                    continue;
                }
                cards.push(card);
            }
        }
        cards.sort(compareCards);

        return (
            <div className="preferans-view container h-screen p-4 space-y-2 w-screen">
                <DndContext onDragEnd={this.onCardDragEnd}>
                    <CardDeck deckId={1} cardIds={cards}/>
                </DndContext>
            </div>
        );
    }
}

export default PreferansView;
