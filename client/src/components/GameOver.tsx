import { useNavigate } from "react-router-dom";
import styles from "../styles/game.module.css";

interface Props {
  result: string;
}
export default function GameOver({ result }: Props) {
  const navigate = useNavigate();

  return (
    <div
      className={styles.gameovercontainer}
      onClick={() => {
        navigate("/");
      }}
    >
      <div className={styles.gameoverbox}>
        <h1 className={styles.gameover}>You {result}</h1>
      </div>
    </div>
  );
}
