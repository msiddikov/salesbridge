import { Grid } from "@mui/material";
import { ReportProps } from "../../types/types";
import { TileGroup } from "./tileGroup";

export {};

export function SimpleReport(props: ReportProps) {
  let repProps = {
    locations: props.locations,
    start: props.start,
    end: props.end,
    count: props.count,
    stats: props.stats,
  } as ReportProps;

  return (
    <Grid xs={12 / props.count}>
      <TileGroup {...repProps} />
    </Grid>
  );
}

export default SimpleReport;
