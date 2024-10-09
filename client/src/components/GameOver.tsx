import { useNavigate } from "react-router-dom";
import styles from "../styles/game.module.css";

interface Props {
  result: string;
  note: string;
}
export default function GameOver({ result, note }: Props) {
  const navigate = useNavigate();

  return (
    <div
      className={styles.gameovercontainer}
      onClick={() => {
        navigate("/");
      }}
    >
      <div className={styles.gameoverbox}>
        <h1 className={styles.gameresult}>You {result}</h1>
        <p className={styles.resultnote}>{note}</p>
      </div>
    </div>
  );
}
