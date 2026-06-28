import { useMemo } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

export type Actor = {
  name: string;
  positionX: number;
  positionY: number;
  grabbed: boolean;
};

// resolve the websocket url. prefer the build-time VITE_WS_URL (for example
// wss://sim.elliotkirk.com/ws in production) and fall back to the local go
// server. note the ws:// scheme - the previous http:// value was incorrect.
function resolveWsUrl(): string {
  const fromEnv = import.meta.env.VITE_WS_URL;
  if (fromEnv && fromEnv.length > 0) {
    return fromEnv;
  }
  return "ws://localhost:8080";
}

// parse the server payload defensively: the first frames can arrive before any
// data, and we never want a malformed message to crash the render loop.
function parseActors(data: unknown): Actor[] {
  if (typeof data !== "string" || data.length === 0) {
    return [];
  }

  try {
    const parsed = JSON.parse(data);
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed.map((d: Record<string, unknown>) => ({
      name: String(d.name),
      positionX: Number(d.positionX),
      positionY: Number(d.positionY),
      grabbed: Boolean(d.grabbed),
    }));
  } catch (e) {
    console.warn("failed to parse world state", e);
    return [];
  }
}

export function useWorldSocket() {
  const { sendMessage, lastMessage, readyState } = useWebSocket(resolveWsUrl(), {
    onOpen: () => console.log("WebSocket connection established."),
    shouldReconnect: () => true,
    reconnectAttempts: Infinity,
    reconnectInterval: 1000,
  });

  const actors = useMemo(() => parseActors(lastMessage?.data), [lastMessage]);

  return {
    actors,
    sendMessage,
    readyState,
    isConnected: readyState === ReadyState.OPEN,
  };
}
