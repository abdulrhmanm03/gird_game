import { useEffect, useState } from "react";
import styles from "../styles/game.module.css";
import Timer from "./Timer";

enum cellContent {
  empty = 0,
  bomb = 1,
  apple = 2,
}

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
      const { data } = JSON.parse(event.data);
      const { pos, content } = data;
      if (pos === index) {
        switch (content) {
          case cellContent.empty:
            setIsBomb(false);
            setIsApple(false);
            break;
          case cellContent.bomb:
            setIsApple(false);
            setIsBomb(true);
            break;
          case cellContent.apple:
            setIsBomb(false);
            setIsApple(true);
            break;
        }
      }
    });
  }, [socket, index, playerScore, setPlayerScore]);

  function handleClick(pos: number) {
    if (!isApple && !isBomb) {
      switch (contains) {
        case cellContent.bomb:
          setIsApple(false);
          setIsBomb(true);
          setPlayerScore(playerScore - 5);
          break;
        case cellContent.apple:
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
