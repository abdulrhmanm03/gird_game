import ModeButton from "../components/ModeButton";
import styles from "../styles/home.module.css";

export default function Home() {
  return (
    <div className={styles.home}>
      <h1 className={styles.gamename}>Fruit or Doom</h1>
      <h4 className={styles.userprompet}>pick the mode you want to play</h4>
      <div className={styles.buttonscontainer}>
        <ModeButton mode={1} />
        <ModeButton mode={2} />
      </div>
    </div>
  );
}
