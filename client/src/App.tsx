import "@pixi/events";
import { useEffect, useMemo, useState } from "react";
import { Actor } from "./Actor";
import "./App.css";

import useWebSocket, { ReadyState } from "react-use-websocket";
import { Stage } from "@pixi/react";

const WS_URL = "http://127.0.0.1:8080";

const App = () => {
  // const [wsOpen, setWsOpen] = useState(false);
  const { sendMessage, lastMessage, readyState } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log("WebSocket connection established.");
    },
    shouldReconnect: () => true,
  });

  // Run when the connection state (readyState) changes
  useEffect(() => {
    console.log("Connection state changed");
    if (readyState === ReadyState.OPEN) {
      sendMessage("start");
    }
  }, [readyState, sendMessage]);

  // let timer: number;
  // document.addEventListener("mousemove", (e) => {
  //   clearTimeout(timer);
  //   timer = setTimeout(() => {
  //     sendMessage(`mouse ${e.pageX} ${e.pageY}`);
  //     console.log(`mouse ${e.pageX} ${e.pageY}`);
  //   }, 1000);
  // });

  const [windowWidth, setWindowWidth] = useState(window.innerWidth);
  const [windowHeight, setWindowHeight] = useState(window.innerHeight);

  const handleResize = () => {
    setWindowWidth(window.innerWidth);
    setWindowHeight(window.innerHeight);
  };

  window.addEventListener("resize", handleResize);

  const actors = useMemo(() => {
    let data;

    try {
      data = JSON.parse(lastMessage?.data);
    } catch (e) {
      console.log(e);
    }

    if (!data) {
      return [];
    }

    return data.map((d: Record<string, string>) => {
      return (
        <Actor
          name={d.name}
          x={Number(d.positionX)}
          y={Number(d.positionY)}
          sendMessage={sendMessage}
        />
      );
    });
  }, [lastMessage, sendMessage]);

  return (
    <Stage
      width={windowWidth}
      height={windowHeight}
      options={{ background: 0xc2b280 }}
    >
      {...actors}
    </Stage>
  );
};

export default App;
