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
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [isRoomActive, setIsRoomActive] = useState(false);
  const [squerContains, setSquerContains] = useState(0);
  const [playerScore, setPlayerScore] = useState(100);
  const [gameResults, setgameResults] = useState("");
  const [resultNote, setResultNote] = useState("");
  const [bombCount, setBombCount] = useState<number | null>(null);
  const [appleCount, setAppleCount] = useState<number | null>(null);

  const roomTimeInMinutes = 5;
  const roomTime = new Date();
  roomTime.setSeconds(roomTime.getSeconds() + roomTimeInMinutes * 60);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:3000/ws");
    if (ws === null) {
      alert("faild to connect to the socket");
      navigate("/");
    }
    setSocket(ws);

    ws.onopen = () => {
      ws.send(JSON.stringify({ mode }));
    };
    ws.addEventListener("message", (event) => {
      const { room_state, result, note, bomb_count, apple_count } = JSON.parse(
        event.data,
      );
      if (room_state === roomState.active) {
        setIsRoomActive(true);
      }
      if (room_state === roomState.gameOver) {
        setgameResults(result);
        setResultNote(note);
      }
      if (apple_count !== undefined || bomb_count !== undefined) {
        console.log("clicked");
        console.log(bomb_count);

        setBombCount(bomb_count);
        setAppleCount(apple_count);
        setTimeout(() => {
          setBombCount(null);
          setAppleCount(null);
        }, 1000);
      }
    });

    return () => {
      ws.close();
    };
  }, [mode, navigate]);

  useEffect(() => {
    if (socket != null && isRoomActive) {
      socket.addEventListener("message", (event) => {
        const { score } = JSON.parse(event.data);
        if (score) {
          setPlayerScore(score);
        }
      });
    }
  }, [socket, isRoomActive]);

  function changeSquereContains(contains: number) {
    setSquerContains(contains);
  }

  function handelMode1ButtonClick(buttonClicked: number) {
    socket?.send(JSON.stringify({ pos: -1, button_clicked: buttonClicked }));
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
                className={styles.button}
                onClick={() => changeSquereContains(1)}
              >
                bomb
              </button>
              <button
                className={styles.button}
                onClick={() => changeSquereContains(2)}
              >
                apple
              </button>
            </div>
          )}

          {mode === 1 && (
            <>
              <div className={styles.buttonscontainer}>
                <button
                  className={styles.button}
                  onClick={() => handelMode1ButtonClick(1)}
                >
                  apples and bomb count
                </button>
                <button
                  className={styles.button}
                  onClick={() => handelMode1ButtonClick(2)}
                >
                  active cell
                </button>
              </div>
              {bombCount !== null && (
                <div className={styles.appleandbomb}>
                  <span className={styles.appleandbombcount}>
                    <p>{bombCount}</p>
                    <img
                      src="/bomb.svg"
                      className={styles.appleandbombcountimg}
                    />
                  </span>
                  <span className={styles.appleandbombcount}>
                    <p>{appleCount}</p>
                    <img
                      src="/apple.svg"
                      className={styles.appleandbombcountimg}
                    />
                  </span>
                </div>
              )}
            </>
          )}

          {socket != null && (
            <Grid socket={socket} mode={mode} contains={squerContains} />
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
