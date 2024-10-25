import { Sprite, useTick } from "@pixi/react";
import { useState } from "react";
import theMan from "./assets/theMan.png"

export const Actor = () => {
  const [position, setPosition] = useState({
    x: Math.floor(Math.random() * window.innerWidth),
    y: Math.floor(Math.random() * window.innerHeight),
  });

  useTick(() => {
    setPosition({
      x: position.x + Math.floor(Math.random() > 0.5 ? 1 : -1),
      y: position.y + Math.floor(Math.random() > 0.5 ? 1 : -1),
    });
  });

  return (
    <>
      <Sprite image={theMan} x={position.x} y={position.y} scale={{x: 0.1, y: 0.1}} />
    </>
  );
};
