import { Sprite} from "@pixi/react";
import theMan from "./assets/theMan.png";

export function Actor({name, x, y}: {name: string, x: number, y: number}) {
  return (
    <>
      <Sprite
        key={name}
        image={theMan}
        x={x}
        y={y}
        scale={{ x: 0.1, y: 0.1 }}
      />
    </>
  );
};
