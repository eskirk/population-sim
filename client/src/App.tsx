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
    // console.log(lastMessage)
    const data = lastMessage?.data?.split("\n");
    // console.log(data)

    if (!data) {
      return [];
    }

    return data.map((d: string) => {
      const details = d && d.split(" ");

      if (!details) return null;
      return (
        <Actor
          name={details[0]}
          x={Number(details[1])}
          y={Number(details[2])}
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
