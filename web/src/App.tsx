import { useEffect, useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import { api } from "./api";

function App() {
  const [count, setCount] = useState(0);
  const [greeting, setGreeting] = useState<string>("");
  const [error, setError] = useState<string>("");

  useEffect(() => {
    api.Greeting.getGreeting({})
      .then((response) => {
        setGreeting(response.message);
      })
      .catch((error) => {
        setError(error.message);
      });
  }, []);

  return (
    <>
      <p>{greeting}</p>
      <p>{error}</p>
    </>
  );
}

export default App;
