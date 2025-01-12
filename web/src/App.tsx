import { useEffect, useState } from "react";
import { api } from "./api";
import { AppShell, Burger, Text } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";

function App() {
  const [opened, { toggle }] = useDisclosure();
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
    <AppShell
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: "sm",
        collapsed: { mobile: !opened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
        <div>Logo</div>
      </AppShell.Header>
      <AppShell.Navbar p="md">Navbar</AppShell.Navbar>
      <AppShell.Main>
        <Text>{greeting}</Text>
        <Text>{error}</Text>
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
