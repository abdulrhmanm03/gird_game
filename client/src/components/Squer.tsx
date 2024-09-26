import styles from "./grid.module.css";

interface squerProps {
  isActive: boolean;
  onClick?: () => void;
}

export default function Squer({ isActive, onClick }: squerProps) {
  return (
    <div
      className={`${styles.squer} ${isActive ? styles.activesquer : ""}`}
      onClick={onClick}
    ></div>
  );
}
