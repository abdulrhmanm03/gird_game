import { useEffect, useState } from "react";
import Squer from "./Squer";
import styles from "./grid.module.css";

export default function Grid() {
  const width = 5;
  const height = 5;

  const initGrid = new Array(width * height).fill(false);
  const [grid, setGrid] = useState(initGrid);

  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onmessage = (event) => {
      const newGrid = JSON.parse(event.data);
      setGrid(newGrid);
    };

    return () => {
      socket.close();
    };
  }, []);

  const handleClick = (i: number) => {
    if (ws) {
      const newGrid = [...grid];
      newGrid[i] = !newGrid[i];
      console.log(i);
      setGrid(newGrid);
      ws.send(JSON.stringify(newGrid));
    }
  };

  return (
    <div className={styles.grid}>
      {grid.map((value, i) => {
        return (
          <Squer isActive={value} key={i} onClick={() => handleClick(i)} />
        );
      })}
    </div>
  );
}
