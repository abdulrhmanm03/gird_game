import { useEffect, useState } from "react";
import styles from "./grid.module.css";
import Mode1Squer from "./Mode1Squer";
import Mode2Squer from "./Mode2Squer";

interface Props {
  socket: WebSocket;
  mode: number;
  contains: number;
}

export default function Grid({ socket, mode, contains }: Props) {
  const width = 5;
  const height = 5;

  const initGrid = new Array(width * height).fill(0);
  const [grid, setGrid] = useState(initGrid);

  useEffect(() => {
    if (mode == 2) {
      socket.onmessage = (event) => {
        const { board } = JSON.parse(event.data);
        setGrid(board);
      };
    }
  }, [mode, socket]);
  return (
    <div className={styles.grid}>
      {grid.map((_, i) => {
        if (mode == 1) {
          return <Mode1Squer key={i} socket={socket} pos={i} />;
        }
        if (mode == 2) {
          return (
            <Mode2Squer
              key={i}
              socket={socket}
              pos={i}
              contains={contains}
              grid={grid}
            />
          );
        }
      })}
    </div>
  );
}
