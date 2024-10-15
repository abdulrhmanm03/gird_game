import { useEffect, useState } from "react";
import styles from "../styles/game.module.css";

interface Props {
  socket: WebSocket;
  index: number;
}

export default function Mode1Squer({ socket, index }: Props) {
  const [isBomb, setIsBomb] = useState(false);
  const [isApple, setIsApple] = useState(false);
  const [isCellActive, setIsCellActive] = useState(false);

  function handleClick(pos: number) {
    socket.send(JSON.stringify({ pos, button_clicked: 0 }));
    console.log("mode 1 is talking from " + pos);
    socket.addEventListener("message", (event) => {
      const { data } = JSON.parse(event.data);
      const { pos, content } = data;

      if (index == pos) {
        switch (content) {
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

  useEffect(() => {
    socket.addEventListener("message", (event) => {
      const { data } = JSON.parse(event.data);
      const { active_cells } = data;
      if (active_cells) {
        if (active_cells.includes(index)) {
          setIsCellActive(true);
          setTimeout(() => setIsCellActive(false), 500);
        }
      }
    });
  }, [socket, index]);

  return (
    <div
      className={`${styles.squer} ${isCellActive ? styles.activesquere : ""}`}
      onClick={() => handleClick(index)}
    >
      {isBomb && (
        <>
          <img className={styles.svg} src="/bomb.svg" />
        </>
      )}
      {isApple && <img className={styles.svg} src="/apple.svg"></img>}
    </div>
  );
}
