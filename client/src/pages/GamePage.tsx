import { useEffect, useState } from "react";
import Waiting from "../components/Waiting";
import Grid from "../components/Grid";
import { useLocation, useNavigate } from "react-router-dom";

export default function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.state;
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRoomActive, setIsRoomActive] = useState(false);
  const [roomId, setRoomId] = useState(-1);
  const [squerContains, setSquerContains] = useState(0);

  useEffect(() => {
    async function getRoomId() {
      try {
        const response = await fetch("http://localhost:3000/getRoomId", {
          method: "POST",
          body: JSON.stringify({ mode }),
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (response.ok) {
          const data = await response.json();
          setRoomId(data.room_id);
        } else {
          alert("something went wrong");
          console.log(response);
          navigate("/");
        }
      } catch {
        alert("faild to connect to server");
        navigate("/");
      }
    }

    getRoomId();
  }, [mode, navigate]);

  useEffect(() => {
    if (roomId > -1) {
      const socket = new WebSocket("ws://localhost:3000/ws");
      setWs(socket);

      if (socket === null) {
        alert("faild to connect to the socket");
        navigate("/");
      }

      socket.onopen = () => {
        socket.send(JSON.stringify({ room_id: roomId, mode }));
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
    }
  }, [mode, roomId, navigate]);

  function changeSquereContains(contains: number) {
    setSquerContains(contains);
  }

  return (
    <>
      {isRoomActive ? (
        <>
          {mode === 2 && (
            <button onClick={() => changeSquereContains(1)}>red</button>
          )}

          {ws != null && (
            <Grid socket={ws} mode={mode} contains={squerContains} />
          )}
        </>
      ) : (
        <Waiting />
      )}
    </>
  );
}
