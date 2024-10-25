import { Sprite, useTick } from "@pixi/react";
import { useState } from "react";
import theMan from "./assets/theMan.png";

export const Actor = () => {
  const [position, setPosition] = useState({
    x: Math.floor(Math.random() * window.innerWidth),
    y: Math.floor(Math.random() * window.innerHeight),
  });

  useTick(() => {
    const setXY = (dimension: "x" | "y"): number => {
      const distance = Math.floor(Math.random() * 5);
      const direction = Math.floor(Math.random() > 0.5 ? 1 : -1);
      const movement = distance * direction;

      if (
        position[dimension] + movement >
          window[dimension == "x" ? "innerWidth" : "innerHeight"] ||
        position[dimension] + movement < 0
      ) {
        return position[dimension] + -movement;
      }

      return position[dimension] + movement;
    };

    setPosition({
      x: setXY("x"),
      y: setXY("y"),
    });
  });

  return (
    <>
      <Sprite
        image={theMan}
        x={position.x}
        y={position.y}
        scale={{ x: 0.1, y: 0.1 }}
      />
    </>
  );
};
