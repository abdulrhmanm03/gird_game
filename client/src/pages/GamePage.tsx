import { useEffect, useState } from "react";
import Grid from "../components/Grid";
import { useLocation, useNavigate } from "react-router-dom";
import styles from "../styles/game.module.css";
import GameOver from "../components/GameOver";
import RoomTimer from "../components/RoomTimer";

enum roomState {
  gameOver = 0,
  waiting = 1,
  active = 2,
}

export default function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.state;
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRoomActive, setIsRoomActive] = useState(false);
  const [squerContains, setSquerContains] = useState(0);
  const [playerScore, setPlayerScore] = useState(100);
  const [gameResults, setgameResults] = useState("");
  const [resultNote, setResultNote] = useState("");

  const roomTimeInMinutes = 5;
  const roomTime = new Date();
  roomTime.setSeconds(roomTime.getSeconds() + roomTimeInMinutes * 60);

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
    socket.addEventListener("message", (event) => {
      const { room_state, result, note } = JSON.parse(event.data);
      if (room_state === roomState.active) {
        setIsRoomActive(true);
      }
      if (room_state === roomState.gameOver) {
        setgameResults(result);
        setResultNote(note);
      }
    });

    return () => {
      socket.close();
    };
  }, [mode, navigate]);

  useEffect(() => {
    if (ws != null && isRoomActive) {
      ws.addEventListener("message", (event) => {
        const { score } = JSON.parse(event.data);
        if (score) {
          setPlayerScore(score);
        }
      });
    }
  }, [ws, isRoomActive]);

  function changeSquereContains(contains: number) {
    setSquerContains(contains);
  }

  return (
    <>
      {isRoomActive ? (
        <div className={styles.gamepage}>
          <div className={styles.scoreandtimer}>
            <h1>Score: {playerScore}</h1>
            <RoomTimer expiryTimestamp={roomTime} />
          </div>
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
          {gameResults && <GameOver result={gameResults} note={resultNote} />}
        </div>
      ) : (
        <div className={styles.waitingcontainer}>
          <h1>waiting for another player to join the room...</h1>
        </div>
      )}
    </>
  );
}
