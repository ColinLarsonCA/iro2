import dayjs, { Dayjs } from "dayjs";

export function displayDate(date: string | Dayjs | undefined): string {
  if (!date) {
    return "";
  }
  if (typeof date === "string") {
    return displayDateFromString(date);
  }
  return displayDateFromDayjs(date);
}

function displayDateFromString(dateStr: string): string {
  const day = dayjs(dateStr);
  return day.isValid() ? day.format("MMM DD, YYYY") : "";
}

function displayDateFromDayjs(date: Dayjs): string {
  return date.format("MMM DD, YYYY");
}