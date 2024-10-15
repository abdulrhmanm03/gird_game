import React from "react";
import styles from "../styles/game.module.css";
import Mode1Squer from "./Mode1Squer";
import Mode2Squer from "./Mode2Squer";

interface Props {
  socket: WebSocket;
  mode: number;
  contains: number;
  playerScore: number;
  setPlayerScore: React.Dispatch<React.SetStateAction<number>>;
}

export default function Grid({
  socket,
  mode,
  contains,
  setPlayerScore,
  playerScore,
}: Props) {
  const width = 5;
  const height = 5;

  const grid = new Array(width * height).fill(0);

  return (
    <div className={styles.grid}>
      {grid.map((_, i) => {
        if (mode == 1) {
          return <Mode1Squer key={i} socket={socket} index={i} />;
        }
        if (mode == 2) {
          return (
            <Mode2Squer
              key={i}
              socket={socket}
              index={i}
              contains={contains}
              playerScore={playerScore}
              setPlayerScore={setPlayerScore}
            />
          );
        }
      })}
    </div>
  );
}
