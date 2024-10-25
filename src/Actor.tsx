import { Sprite, useTick } from '@pixi/react';
import { useState } from 'react';

export const Actor = () => {
  const [position, setPosition] = useState({ x: 300, y: 150 });
  const bunnyUrl = 'https://pixijs.io/pixi-react/img/bunny.png';

  useTick(() => {
    setPosition({
      x: position.x + Math.floor(Math.random() > 0.5 ? 1 : -1),
      y: position.y + Math.floor(Math.random() > 0.5 ? 1 : -1),
    });
  });

  return (
    <>
      <Sprite image={bunnyUrl} x={position.x} y={position.y} />
    </>
  );
}