import { useState } from "react";
import { Grid } from "@mui/material";
import { TagsSelect } from "./tagSelect";
import { Location, ReportProps } from "../../types/types";
import { JfmReportDownloader } from "./downloader";
import { TileGroup } from "../tile/tileGroup";

export {};

export function JfmReport(props: ReportProps) {
  const [tags, setTags] = useState<string[]>([]);

  let repProps = {
    locations: props.locations,
    tags: tags,
    start: props.start,
    end: props.end,
    count: props.count,
    stats: getStats(props.locations),
  } as ReportProps;

  return (
    <Grid xs={12 / props.count}>
      <TagsSelect
        tags={props.settingsData.tags}
        selected={tags}
        setSelected={setTags}
      />

      <TileGroup {...repProps} />

      <JfmReportDownloader {...repProps} />
    </Grid>
  );
}

export default JfmReport;

function getStats(locations: Location[]): string[] {
  let isIntegrated = locations.map((l) => l.isIntegrated).indexOf(true) > -1;
  return [
    "expenses",
    "sales",
    "salesNo",
    "roi",
    "newLeads",
    "bookings",
    "noShows",
    "showNoSale",
    "leadsConv",
    "bookingsConv",
    "showRate",
    isIntegrated ? "zenotiMembersNo" : "membershipConv",
  ] as string[];
}
