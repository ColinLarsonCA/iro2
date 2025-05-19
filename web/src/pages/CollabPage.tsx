import {
  Badge,
  Container,
  Grid,
  Group,
  JsonInput,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useParams } from "react-router-dom";
import { api } from "../api/api";
import { useEffect, useState } from "react";
import { Collab, CollabSchedule } from "../pb/collabcafe";
import dayjs from "dayjs";
import { displayDate } from "../util/dates";
import { EventCard } from "../components/EventCard";

interface DateRange {
  start: dayjs.Dayjs | undefined;
  end: dayjs.Dayjs | undefined;
}

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
  const dateRange = getDateRange(collab?.content?.schedule || { events: [] });
  const dateRangeStr = dateRangeToString(dateRange);
  const officialWebsite = collab?.content?.officialWebsite?.url;
  const events = collab?.content?.schedule?.events || [];
  return (
    <Container>
      <Stack>
        <Text>Posted: {displayDate(collab?.postedDate)}</Text>
        <Title>{collab?.content?.title}</Title>
        <Group>
          <Badge color="blue">Series: {collab?.content?.series}</Badge>
          {dateRangeStr && <Badge color="green">{dateRangeStr}</Badge>}
        </Group>
        {officialWebsite && (
          <Text>
            Official website:{" "}
            <a href={officialWebsite} target="_blank" rel="noreferrer">
              {officialWebsite}
            </a>
          </Text>
        )}
        <Text>{collab?.summary?.description}</Text>
        {events && (
          <Grid>
            {events.filter((e) => e.location).map((event) => (
              <Grid.Col span={{ base: 12, md: 6, lg: 4 }}>
                <EventCard event={event} />
              </Grid.Col>
            ))}
          </Grid>
        )}
      </Stack>
      <JsonInput
        disabled
        value={JSON.stringify(collab)}
        formatOnBlur
        autosize
      />
    </Container>
  );
}

function getDateRange(schedule: CollabSchedule): DateRange {
  const startDates = schedule.events
    .map((event) => dayjs(event.startDate))
    .sort((a, b) => a.diff(b))
    .filter((date) => date.isValid());
  const endDates = schedule.events
    .map((event) => dayjs(event.endDate))
    .sort((a, b) => a.diff(b))
    .filter((date) => date.isValid());
  return {
    start: startDates.length > 0 ? startDates[0] : undefined,
    end: endDates.length > 0 ? endDates[endDates.length - 1] : undefined,
  };
}

function dateRangeToString(dateRange: DateRange): string {
  if (!dateRange.start) {
    return "";
  }
  if (dateRange.start && !dateRange.end) {
    return "Start: " + displayDate(dateRange.start);
  }
  return `${displayDate(dateRange.start)} - ${displayDate(dateRange.end)}`;
}
