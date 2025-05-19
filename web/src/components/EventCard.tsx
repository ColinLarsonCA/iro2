import { Card, Stack, Text } from "@mantine/core";
import { CollabEvent } from "../pb/collabcafe";
import { displayDate } from "../util/dates";

export interface EventCardProps {
  event: CollabEvent;
}

export function EventCard(props: EventCardProps) {
  const { event } = props;
  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ height: "100%" }}>
      <Stack>
        <Text>Location: {event.location}</Text>
        {event.startDate && <Text>From: {displayDate(event.startDate)}</Text>}
        {event.endDate && <Text>To: {displayDate(event.endDate)}</Text>}
      </Stack>
    </Card>
  );
}
