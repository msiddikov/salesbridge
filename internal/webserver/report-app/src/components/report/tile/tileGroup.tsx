import { Grid, Button } from "@mui/material";
import { useState } from "react";
import { Location } from "../../types/types";
import { Tile } from "./tile";

export const TileGroup = ({
  locations,
  tags,
  start,
  end,
  count,
  stats,
}: {
  locations: Location[];
  tags: string[];
  start: Date;
  end: Date;
  count: number;
  stats: string[];
}) => {
  const [version, setVersion] = useState(0);

  const versionUp = () => {
    setVersion(version + 1);
  };
  return (
    <Grid xs={12}>
      <Button onClick={versionUp}>Refresh all</Button>
      <Grid
        container
        columns={{ xs: 4, sm: 8, md: 12 }}
        sx={{
          padding: "35px",
        }}
      >
        {stats.map((stat) => {
          return (
            <Tile
              locations={locations}
              tags={tags}
              stat={stat}
              start={start}
              end={end}
              count={count}
              version={version}
            />
          );
        })}
      </Grid>
    </Grid>
  );
};
