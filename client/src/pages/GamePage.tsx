import { useEffect, useState } from "react";
import Waiting from "../components/Waiting";
import Grid from "../components/Grid";
import { useLocation } from "react-router-dom";

export default function GamePage() {
  const location = useLocation();
  const mode = location.state;
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRoomActive, setIsRoomActive] = useState(false);
  const [roomId, setRoomId] = useState(-1);
  const [squerContains, setSquerContains] = useState(0);

  useEffect(() => {
    async function getRoomId() {
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
        alert("Failed to connect to server");
        console.log(response);
      }
    }

    getRoomId();
  }, [mode]);

  useEffect(() => {
    if (roomId > -1) {
      const socket = new WebSocket("ws://localhost:3000/ws");
      setWs(socket);

      if (socket == null) {
        alert("faild to connect to the socket");
      }

      socket.onopen = () => {
        socket.send(JSON.stringify({ room_id: roomId, mode }));
      };
      socket.onmessage = (event) => {
        const massage = JSON.parse(event.data);
        if (massage.room_state == 1) {
          setIsRoomActive(true);
        }
      };

      return () => {
        socket.close();
      };
    }
  }, [mode, roomId]);

  function changeSquereContains(contains: number) {
    setSquerContains(contains);
  }

  return (
    <>
      {isRoomActive ? (
        <>
          <button onClick={() => changeSquereContains(1)}>red</button>
          {ws != null && (
            <Grid socket={ws} mode={mode} contains={squerContains} />
          )}{" "}
        </>
      ) : (
        <Waiting />
      )}
    </>
  );
}
