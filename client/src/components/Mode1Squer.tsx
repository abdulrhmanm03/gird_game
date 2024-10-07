import { useState } from "react";
import styles from "../styles/game.module.css";

interface Props {
  socket: WebSocket;
  pos: number;
}

export default function Mode1Squer({ socket, pos }: Props) {
  const [isBomb, setIsBomb] = useState(false);
  const [isApple, setIsApple] = useState(false);

  function handleClick(pos: number) {
    socket.send(JSON.stringify({ pos }));
    console.log("mode1 is talking from " + pos);
    socket.addEventListener("message", (event) => {
      const { squere_index, squere_content } = JSON.parse(event.data);

      if (squere_index == pos) {
        switch (squere_content) {
          case 1:
            setIsBomb(true);
            setTimeout(() => {
              setIsBomb(false);
            }, 500);
            break;

          case 2:
            setIsApple(true);
            setTimeout(() => {
              setIsApple(false);
            }, 500);
            break;

          default:
            break;
        }
      }
    });
  }

  return (
    <div className={`${styles.squer} `} onClick={() => handleClick(pos)}>
      {isBomb && (
        <>
          <img className={styles.svg} src="/bomb.svg" />
        </>
      )}
      {isApple && <img className={styles.svg} src="/apple.svg"></img>}
    </div>
  );
}
