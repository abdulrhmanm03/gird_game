import { useState, useEffect } from "react";
import styles from "../styles/game.module.css";

interface Props {
  interval: number;
}

export default function Timer({ interval }: Props) {
  const [seconds, setSeconds] = useState(interval);

  useEffect(() => {
    if (seconds > 0) {
      const timer = setTimeout(() => setSeconds(seconds - 1), 1000);
      return () => clearTimeout(timer);
    }
  }, [seconds]);

  return <div className={styles.timer}>{seconds}</div>;
}
