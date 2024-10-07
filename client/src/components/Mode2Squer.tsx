import { useEffect, useState } from "react";
import styles from "../styles/game.module.css";
import Timer from "./Timer";

interface Props {
  socket: WebSocket;
  index: number;
  contains: number;
}

export default function Mode2Squer({ socket, index, contains }: Props) {
  const [isBomb, setIsBomb] = useState(false);
  const [isApple, setIsApple] = useState(false);

  useEffect(() => {
    socket.addEventListener("message", (event) => {
      const { pos } = JSON.parse(event.data);
      if (pos === index) {
        setIsBomb(false);
        setIsApple(false);
      }
    });
  }, [socket, index]);

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
    <div className={styles.squer} onClick={() => handleClick(index)}>
      {isBomb && (
        <>
          <img className={styles.svg} src="/bomb.svg" />
          <Timer interval={15} />
        </>
      )}
      {isApple && (
        <>
          <img className={styles.svg} src="/apple.svg" />
          <Timer interval={15} />
        </>
      )}
    </div>
  );
}
