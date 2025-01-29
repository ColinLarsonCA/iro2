import { useEffect, useState } from "react";
import { api } from "./api/api";
import { AppShell, Burger, Grid, Space, Text, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Collab } from "./pb/collabcafe";
import { CollabCard } from "./components/CollabCard";

function App() {
  const [opened, { toggle }] = useDisclosure();
  const [collabs, setCollabs] = useState<Collab[]>([]);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");

  const listCollabs = () => {
    api.CollabCafe.listCollabs({ language: "en" })
      .then((response) => {
        setCollabs(response.collabs);
      })
      .catch((error) => {
        setError(error.message);
      });
  };
  useEffect(() => {
    listCollabs();
  }, []);

  useEffect(() => {
    if (search === "") {
      listCollabs();
    } else {
      api.CollabCafe.searchCollabs({ language: "en", query: search })
        .then((response) => {
          setCollabs(response.collabs);
        })
        .catch((error) => {
          setError(error.message);
        });
    }
  }, [search]);

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
        <TextInput
          placeholder="Search"
          onBlur={(event) => {
            setSearch(event.currentTarget.value);
          }}
          onKeyUp={(event) => {
            if (event.key === "Enter") {
              setSearch(event.currentTarget.value);
            }
            if (event.key === "Backspace" && event.currentTarget.value === "") {
              setSearch("");
            }
          }}
        />
        {error && <Text c="red">{error}</Text>}
        <Space h="md" />
        <Grid>
          {(collabs || []).map((collab) => (
            <Grid.Col span={{ base: 12, md: 6, lg: 4 }}>
              <CollabCard collab={collab} />
            </Grid.Col>
          ))}
        </Grid>
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
