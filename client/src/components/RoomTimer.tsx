import { TimerSettings, useTimer } from "react-timer-hook";
import styles from "../styles/game.module.css";

export default function RoomTimer({ expiryTimestamp }: TimerSettings) {
  const { seconds, minutes } = useTimer({
    expiryTimestamp,
    onExpire: () => console.warn("onExpire called"),
  });

  return (
    <div className={styles.roomtimer}>
      <span>{minutes}</span>:<span>{seconds}</span>
    </div>
  );
}
