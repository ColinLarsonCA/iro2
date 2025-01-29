import { Container, JsonInput, Text } from "@mantine/core";
import { useParams } from "react-router-dom";
import { api } from "../api/api";
import { useEffect, useState } from "react";
import { Collab } from "../pb/collabcafe";

export function CollabPage() {
  const { id } = useParams();
  const [collab, setCollab] = useState<Collab | undefined>(undefined);
  useEffect(() => {
    api.CollabCafe.getCollab({ id: String(id), language: "en" })
      .then((response) => {
        setCollab(response.collab);
      })
      .catch((error) => {
        console.error(error);
      });
  }, [id]);
  api.CollabCafe.getCollab({ id: String(id), language: "en" });
  return (
    <Container>
      <JsonInput disabled value={JSON.stringify(collab)} formatOnBlur autosize />
    </Container>
  );
}
