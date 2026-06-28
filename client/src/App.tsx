import { useEffect } from "react";
import { Scene } from "./Scene";
import { useWorldSocket } from "./useWorldSocket";
import "./App.css";

const App = () => {
  const { actors, sendMessage, isConnected } = useWorldSocket();

  // the server streams state continuously, but we send "start" on connect for
  // parity with the original protocol and to confirm the outbound channel.
  useEffect(() => {
    if (isConnected) {
      sendMessage("start");
    }
  }, [isConnected, sendMessage]);

  return <Scene actors={actors} sendMessage={sendMessage} />;
};

export default App;
