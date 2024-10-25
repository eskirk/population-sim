import { useMemo, useState } from "react";
import { Actor } from "./Actor";
import "./App.css";

import { Stage } from "@pixi/react";

const App = () => {
  const [windowWidth, setWindowWidth] = useState(window.innerWidth);
  const [windowHeight, setWindowHeight] = useState(window.innerHeight);

  const handleResize = () => {
    setWindowWidth(window.innerWidth);
    setWindowHeight(window.innerHeight);
  };

  window.addEventListener("resize", handleResize);

  const actors = useMemo(
    () =>
      [...Array(Math.floor(Math.random() * 20))].map(() => {
        return <Actor />;
      }),
    []
  );

  return (
    <Stage
      width={windowWidth}
      height={windowHeight}
      options={{ background: 0x1099bb }}
    >
      {...actors}
    </Stage>
  );
};

export default App;
