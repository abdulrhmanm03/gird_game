import { Link } from "react-router-dom";

interface Props {
  mode: number;
}

export default function ModeButton({ mode }: Props) {
  return (
    <button>
      <Link to="/game" state={mode}>
        mode {mode}
      </Link>
    </button>
  );
}
