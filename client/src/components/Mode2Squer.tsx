import { useEffect, useState } from "react";
import styles from "./grid.module.css";

interface Props {
  socket: WebSocket;
  pos: number;
  contains: number;
  grid: number[];
}

export default function Mode2Squer({ socket, pos, contains, grid }: Props) {
  const [isActive, setIsActive] = useState(false);
  useEffect(() => {
    if (grid[pos] == 0) {
      setIsActive(false);
    }
    if (grid[pos] == 1) {
      setIsActive(true);
    }
  }, [grid, pos]);
  function handleClick(pos: number) {
    if (contains == 1) {
      setIsActive(true);
    }
    socket.send(JSON.stringify({ pos, contains }));
    console.log("mode 2 is talking from " + pos);
  }
  return (
    <div
      className={`${styles.squer} ${isActive ? styles.activesquer : ""} `}
      onClick={() => handleClick(pos)}
    ></div>
  );
}
