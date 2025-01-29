import { useEffect, useState } from "react";
import { api } from "../api/api";
import { Container, Grid, Space, Text, TextInput } from "@mantine/core";
import { Collab } from "../pb/collabcafe";
import { CollabCard } from "../components/CollabCard";

export function HomePage() {
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
    <Container>
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
    </Container>
  );
}
