import { Card, Stack, Text } from "@mantine/core";
import { CollabEvent } from "../pb/collabcafe";
import { displayDate } from "../util/dates";

export interface EventCardProps {
  event: CollabEvent;
}

export function EventCard(props: EventCardProps) {
  const { event } = props;
  const from = event.startDate ? displayDate(event.startDate) : "";
  const to = event.endDate ? displayDate(event.endDate) : "";
  const location = event.location || "";
  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder style={{ height: "100%" }}>
      <Stack gap="xs">
        {location && (
          <Text>
            Location: <a href={googleMapsSearchLink(location)} target="_blank" rel="noreferrer">{location}</a>
          </Text>
        )}
        <Text>Period: {event.period}</Text>
        {from && <Text>From: {from}</Text>}
        {to && <Text>To: {to}</Text>}
        {event.mapLink && (
          <Text>
            <a href={event.mapLink} target="_blank" rel="noreferrer">{event.mapLink}</a>
          </Text>
        )}
      </Stack>
    </Card>
  );
}

function googleMapsSearchLink(location: string): string {
  return `https://maps.google.com/?q=${encodeURIComponent(location)}`;
}
