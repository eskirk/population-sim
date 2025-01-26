import "@pixi/events";
import { Sprite } from "@pixi/react";
import theMan from "./assets/theMan.png";
import { SendMessage } from "react-use-websocket";

export function Actor({ name, x, y, sendMessage }: { name: string; x: number; y: number, sendMessage: SendMessage }) {
  const sendGrabbedMessage = () => {
    console.log(`grabbed ${name}`);
    sendMessage(`grabbed ${name}`);
  }

  return (
    <Sprite
      interactive={true}
      key={name}
      image={theMan}
      x={x}
      y={y}
      scale={{ x: 0.1, y: 0.1 }}
      onmousedown={sendGrabbedMessage}
    />
  );
}
