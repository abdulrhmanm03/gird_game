import styles from "../styles/game.module.css";

interface Props {
  socket: WebSocket;
  pos: number;
}

export default function Mode1Squer({ socket, pos }: Props) {
  function handleClick(pos: number) {
    socket.send(JSON.stringify({ pos }));
    console.log("mode1 is talking from " + pos);
  }
  return (
    <div className={`${styles.squer} `} onClick={() => handleClick(pos)}></div>
  );
}
