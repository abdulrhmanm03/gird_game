import { useEffect, useState } from "react";
import styles from "../styles/game.module.css";

interface Props {
  socket: WebSocket;
  pos: number;
  contains: number;
  grid: number[];
}

export default function Mode2Squer({ socket, pos, contains, grid }: Props) {
  const [isBomb, setIsBomb] = useState(false);
  const [isApple, setIsApple] = useState(false);
  useEffect(() => {
    switch (grid[pos]) {
      case 1:
        setIsApple(false);
        setIsBomb(true);
        break;
      case 2:
        setIsBomb(false);
        setIsApple(true);
        break;
      default:
        setIsBomb(false);
        setIsApple(false);
        break;
    }
  }, [grid, pos]);

  function handleClick(pos: number) {
    switch (contains) {
      case 1:
        setIsApple(false);
        setIsBomb(true);
        break;
      case 2:
        setIsBomb(false);
        setIsApple(true);
        break;
    }
    socket.send(JSON.stringify({ pos, contains }));
  }

  return (
    <div className={styles.squer} onClick={() => handleClick(pos)}>
      {isBomb && <img className={styles.svg} src="/bomb.svg"></img>}
      {isApple && <img className={styles.svg} src="/apple.svg"></img>}
    </div>
  );
}
