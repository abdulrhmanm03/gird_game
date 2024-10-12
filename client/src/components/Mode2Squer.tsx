import { useEffect, useState } from "react";
import styles from "../styles/game.module.css";
import Timer from "./Timer";

interface Props {
  socket: WebSocket;
  index: number;
  contains: number;
  playerScore: number;
  setPlayerScore: React.Dispatch<React.SetStateAction<number>>;
}

export default function Mode2Squer({
  socket,
  index,
  contains,
  playerScore,
  setPlayerScore,
}: Props) {
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
    if (!isApple && !isBomb) {
      switch (contains) {
        case 1:
          setIsApple(false);
          setIsBomb(true);
          setPlayerScore(playerScore - 5);
          break;
        case 2:
          setIsBomb(false);
          setIsApple(true);
          setPlayerScore(playerScore - 5);
          break;
      }
      socket.send(JSON.stringify({ pos, contains }));
    }
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
