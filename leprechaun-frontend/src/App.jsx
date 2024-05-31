import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import NavBar from "./NavBar/NavBar";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="flex flex-row">
      <NavBar />
    </div>
  );
}

export default App;
