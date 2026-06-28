import { useFrame } from "@react-three/fiber";
import { useTexture } from "@react-three/drei";
import { useEffect, useRef } from "react";
import * as THREE from "three";
import type { SendMessage } from "react-use-websocket";
import type { Actor as ActorData } from "./useWorldSocket";
import theMan from "./assets/theMan.png";

const ACTOR_HEIGHT = 1.5;
// fraction of the remaining distance covered per frame; lower is smoother.
const LERP_FACTOR = 0.2;
// grabbed actors are enlarged so the grabbed state reads clearly.
const GRABBED_SCALE = 1.35;

export function Actor({
  actor,
  worldWidth,
  worldHeight,
  scale,
  sendMessage,
}: {
  actor: ActorData;
  worldWidth: number;
  worldHeight: number;
  scale: number;
  sendMessage: SendMessage;
}) {
  const group = useRef<THREE.Group>(null);
  const target = useRef(new THREE.Vector3());
  const initialized = useRef(false);
  const texture = useTexture(theMan);

  const aspect =
    texture.image && texture.image.height
      ? texture.image.width / texture.image.height
      : 1;
  const sizeMultiplier = actor.grabbed ? GRABBED_SCALE : 1;
  const spriteHeight = ACTOR_HEIGHT * sizeMultiplier;
  const spriteWidth = spriteHeight * aspect;

  // base position on the ground; server (x, y) maps onto the centered world
  // ground plane as (x, z).
  const groundPosition = (a: ActorData) =>
    target.current.set(
      a.positionX / scale - worldWidth / scale / 2,
      0,
      a.positionY / scale - worldHeight / scale / 2,
    );

  // snap to the first known position so actors do not fly in from the origin.
  useEffect(() => {
    if (group.current && !initialized.current) {
      group.current.position.copy(groundPosition(actor));
      initialized.current = true;
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // the server pushes only ~20 updates/sec, so interpolate toward the latest
  // position every frame for smooth 60fps motion.
  useFrame(() => {
    if (!group.current) return;
    group.current.position.lerp(groundPosition(actor), LERP_FACTOR);
  });

  return (
    <group
      ref={group}
      onClick={(e) => {
        e.stopPropagation();
        console.log(`grabbed ${actor.name}`);
        sendMessage(`grabbed ${actor.name}`);
      }}
    >
      {/* theMan.png is a black silhouette, so a color tint would be invisible
          (black * color = black). The billboard always faces the camera. */}
      <sprite position={[0, spriteHeight / 2, 0]} scale={[spriteWidth, spriteHeight, 1]}>
        <spriteMaterial map={texture} transparent />
      </sprite>

      {/* ground ring highlight shown only while the actor is grabbed */}
      <mesh
        rotation={[-Math.PI / 2, 0, 0]}
        position={[0, 0.02, 0]}
        visible={actor.grabbed}
      >
        <ringGeometry args={[0.45, 0.7, 32]} />
        <meshBasicMaterial
          color="#ff5555"
          transparent
          opacity={0.9}
          side={THREE.DoubleSide}
        />
      </mesh>
    </group>
  );
}
