/// <reference types="vite/client" />

interface ImportMetaEnv {
  // websocket endpoint for the go server, e.g. wss://sim.elliotkirk.com/ws.
  readonly VITE_WS_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
