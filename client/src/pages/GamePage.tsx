import { useEffect, useState } from "react";
import Grid from "../components/Grid";
import { useLocation, useNavigate } from "react-router-dom";
import styles from "../styles/game.module.css";

export default function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.state;
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRoomActive, setIsRoomActive] = useState(false);
  const [squerContains, setSquerContains] = useState(0);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:3000/ws");
    if (socket === null) {
      alert("faild to connect to the socket");
      navigate("/");
    }
    setWs(socket);

    socket.onopen = () => {
      socket.send(JSON.stringify({ mode }));
    };
    socket.onmessage = (event) => {
      const { room_state } = JSON.parse(event.data);
      if (room_state === 0) {
        setIsRoomActive(true);
      }
    };

    return () => {
      socket.close();
    };
  }, [mode, navigate]);

  function changeSquereContains(contains: number) {
    setSquerContains(contains);
  }

  return (
    <>
      {isRoomActive ? (
        <>
          {mode === 2 && (
            <div className={styles.buttonscontainer}>
              <button
                className={styles.mode2buttons}
                onClick={() => changeSquereContains(1)}
              >
                bomb
              </button>
              <button
                className={styles.mode2buttons}
                onClick={() => changeSquereContains(2)}
              >
                apple
              </button>
            </div>
          )}

          {ws != null && (
            <Grid socket={ws} mode={mode} contains={squerContains} />
          )}
        </>
      ) : (
        <div className={styles.waitingcontainer}>
          <h1>waiting for another player to join the room...</h1>
        </div>
      )}
    </>
  );
}
