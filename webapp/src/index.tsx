/// <reference types="telegram-web-app"/>

import "./style.css";

import * as React from "react";
import * as ReactDOM from "react-dom/client";
import interact from "interactjs";
import { InteractEvent } from "@interactjs/types";

interface CardProps {
    value: string;
};

class Card extends React.Component<CardProps, {}> {
    render(): JSX.Element {
        return (
            <li className="card">{this.props.value}</li>
        );
    }
}

interface CardDeckProps {
    deckId: string;
    cards: string[];
}

class CardDeck extends React.Component<CardDeckProps, {}> {
    render(): JSX.Element {
        return (
            <ul id={this.props.deckId} className="card-deck">
                {this.props.cards.map(card => (
                    <Card key={card} value={card}/>
                ))}
            </ul>
        );
    }
}

class PreferansView extends React.Component<{}, {}> {
    render(): JSX.Element {
        return (
            <div>
                <CardDeck deckId="deck1" cards={["card1", "card2"]}/>
                <CardDeck deckId="deck2" cards={["card3", "card4"]}/>
            </div>
        );
    }
}

interact(".card-deck").dropzone({
    listeners: {
        ondragenter (event: InteractEvent) {
            var target = event.target;
            var relatedTarget = event.relatedTarget;
            relatedTarget.setAttribute("drag-target", target.id);
        },
        ondragleave (event: InteractEvent) {
            var relatedTarget = event.relatedTarget;
            relatedTarget.removeAttribute("drag-target");
        },
        ondrop (event: InteractEvent) {
            console.log("ondrop");
        },
        ondropdeactivate (event: InteractEvent) {
            var relatedTarget = event.relatedTarget;
            relatedTarget.removeAttribute("drag-target");
        }
    }
});
  
interact(".card").draggable({
    inertia: true,
    listeners: {
        move (event: InteractEvent) {
            var target = event.target;
            var x = (parseFloat(target.getAttribute("drag-x")) || 0) + event.dx;
            var y = (parseFloat(target.getAttribute("drag-y")) || 0) + event.dy;
            target.style.transform = `translate(${x}px, ${y}px)`;
            target.setAttribute("drag-x", x.toString());
            target.setAttribute("drag-y", y.toString());
        },
        end (event: InteractEvent) {
            var target = event.target;
            target.style.removeProperty("transform");
            target.removeAttribute("drag-x");
            target.removeAttribute("drag-y");
        }
    }
});

const content = <PreferansView/>

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(content);
