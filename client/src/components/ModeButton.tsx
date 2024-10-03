import { useNavigate } from "react-router-dom";
import styles from "../styles/home.module.css";

interface Props {
  mode: number;
}

export default function ModeButton({ mode }: Props) {
  const navigate = useNavigate();

  return (
    <button
      className={styles.modeButton}
      onClick={() => {
        navigate("/game", { state: mode });
      }}
    >
      mode {mode}
    </button>
  );
}
