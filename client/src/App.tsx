import { useMemo, useState } from "react";
import { Actor } from "./Actor";
import "./App.css";

import useWebSocket from "react-use-websocket";
import { Stage } from "@pixi/react";

const WS_URL = "http://127.0.0.1:8080";

const App = () => {
  const { sendMessage, lastMessage } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log("WebSocket connection established.");
    },
  });

  sendMessage("start");

  const [windowWidth, setWindowWidth] = useState(window.innerWidth);
  const [windowHeight, setWindowHeight] = useState(window.innerHeight);

  const handleResize = () => {
    setWindowWidth(window.innerWidth);
    setWindowHeight(window.innerHeight);
  };

  window.addEventListener("resize", handleResize);

  const actors = useMemo(
    () => {
      // console.log(lastMessage)
      const data = lastMessage?.data?.split('\n')
      // console.log(data)

      if (!data) {
        return []
      }

      return data.map((d: string) => {
        const details = d && d.split(' ');

        if (!details) return null;
        return <Actor name={details[0]} x={Number(details[1])} y={Number(details[2])} />;
      })
    },
    [lastMessage]
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
