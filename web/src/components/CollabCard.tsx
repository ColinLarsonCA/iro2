import { Card, Image, Group, Text, Badge, Button, ScrollArea } from "@mantine/core";
import { Collab } from "../pb/collabcafe";
import dayjs from "dayjs";
import { displayDate } from "../util/dates";

export interface CollabCardProps {
  collab: Collab;
}

export function CollabCard(props: CollabCardProps) {
  const { collab } = props;
  console.log(collab);
  const postedRecently = dayjs().diff(dayjs(collab.postedDate), 'day') <= 3;
  const thumbnailUrl = collab.summary?.thumbnail || "https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/images/bg-8.png";
  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ height: "100%" }}>
      <Card.Section style={{ position: 'relative' }}>
        <Image
          src={thumbnailUrl}
          height={160}
          alt="Norway"
        />
        {postedRecently && (
          <Badge
            gradient={{ from: 'yellow', to: 'pink', deg: 90 }}
            variant="gradient"
            mt="sm"
            style={{ position: 'absolute', bottom: 10, right: 10 }}
          >
            New!
          </Badge>
        )}
      </Card.Section>

      <Group justify="space-between" mt="md" mb="xs">
        <Text>{collab.summary?.title}</Text>
        <Text c="dimmed">{displayDate(collab.postedDate)}</Text>
      </Group>

      <ScrollArea h={200}>
        <Text size="sm">
          {collab.summary?.description}
        </Text>
      </ScrollArea>

      <div style={{ flexGrow: 1 }} />
      <Button variant="light" component="a" color="blue" fullWidth mt="md" radius="md" href={`/collab/${collab.id}`}>
        Learn more
      </Button>
    </Card>
  );
}
