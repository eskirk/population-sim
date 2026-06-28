import { Suspense } from "react";
import { Canvas } from "@react-three/fiber";
import { Grid, OrbitControls } from "@react-three/drei";
import type { SendMessage } from "react-use-websocket";
import { Actor } from "./Actor";
import type { Actor as ActorData } from "./useWorldSocket";

// server world dimensions (see server/pkg/environment/environment.go).
const WORLD_WIDTH = 2000;
const WORLD_HEIGHT = 1000;
// world units = server pixels / SCALE.
const SCALE = 50;

const groundWidth = WORLD_WIDTH / SCALE;
const groundDepth = WORLD_HEIGHT / SCALE;

export function Scene({
  actors,
  sendMessage,
}: {
  actors: ActorData[];
  sendMessage: SendMessage;
}) {
  return (
    <Canvas
      style={{ width: "100vw", height: "100vh" }}
      camera={{ position: [0, 18, 26], fov: 50 }}
      onCreated={({ gl }) => {
        // chrome lets a canvas be dragged as an image, which shows a ghost
        // sprite stuck to the cursor; cancel the native drag so click-drag
        // only ever drives OrbitControls.
        gl.domElement.addEventListener("dragstart", (e) => e.preventDefault());
      }}
    >
      <color attach="background" args={["#10131c"]} />
      <fog attach="fog" args={["#10131c", 45, 100]} />

      <ambientLight intensity={0.7} />
      <directionalLight position={[15, 25, 10]} intensity={1.1} />

      {/* ground plane sized to the world */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.01, 0]}>
        <planeGeometry args={[groundWidth, groundDepth]} />
        <meshStandardMaterial color="#c2b280" />
      </mesh>

      <Grid
        args={[groundWidth, groundDepth]}
        cellSize={1}
        cellThickness={0.6}
        cellColor="#9a8f6b"
        sectionSize={5}
        sectionThickness={1}
        sectionColor="#6f7689"
        fadeDistance={120}
        infiniteGrid={false}
      />

      <Suspense fallback={null}>
        {actors.map((actor) => (
          <Actor
            key={actor.name}
            actor={actor}
            worldWidth={WORLD_WIDTH}
            worldHeight={WORLD_HEIGHT}
            scale={SCALE}
            sendMessage={sendMessage}
          />
        ))}
      </Suspense>

      <OrbitControls target={[0, 0, 0]} maxPolarAngle={Math.PI / 2.1} enableDamping />
    </Canvas>
  );
}
