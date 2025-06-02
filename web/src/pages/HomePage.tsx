import { useEffect, useState } from "react";
import { api } from "../api/api";
import { ActionIcon, Container, Grid, Space, Text, TextInput, Tooltip } from "@mantine/core";
import { Collab } from "../pb/collabcafe";
import { CollabCard } from "../components/CollabCard";
import { IconSend2 } from "@tabler/icons-react";
import { useSearchParams } from "react-router-dom";

export function HomePage() {
  const [searchParams] = useSearchParams();
  const searchQuery = searchParams.get("s") || "";
  const [exampleSearch] = useState(getExampleSearch());
  const [collabs, setCollabs] = useState<Collab[]>([]);
  const [error, setError] = useState("");
  const [searchInProgress, setSearchInProgress] = useState("");
  const [search, setSearch] = useState(searchQuery || "");
  useEffect(() => {
    if (search === "") {
    api.CollabCafe.listCollabs({ language: "en" })
      .then((response) => {
        setCollabs(response.collabs);
      })
      .catch((error) => {
        setError(error.message);
      });
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
  const searchButton = (
    <Tooltip label="Search">
      <ActionIcon
        variant="light"
        onClick={() => {
          setSearch(searchInProgress);
        }}
      >
        <IconSend2 />
      </ActionIcon>
    </Tooltip>
  )
  return (
    <Container>
      <TextInput
        style={{ maxWidth: 300 }}
        defaultValue={search}
        placeholder={`Try searching for ${exampleSearch}`}
        onChange={(event) => {
          setSearchInProgress(event.currentTarget.value);
        }}
        rightSection={searchButton}
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
    </Container>
  );
}

function getExampleSearch(): string {
  const examples = ["Kuromi", "Sakamoto Days", "Sanrio", "Lawson", "Miku"];
  return examples[Math.floor(Math.random() * examples.length)];
}
