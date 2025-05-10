import {
  Grid,
  Box,
  Typography,
  Tooltip,
  CircularProgress,
  Button,
} from "@mui/material";
import { useState, useEffect } from "react";
import { fetchServer } from "../../server/server";
import { stats } from "./stats";
import { Location } from "../../types/types";
import InfoIcon from "@mui/icons-material/Info";
import { NumberTile } from "./numberTile";
import { LineTile, lineTileDataMock } from "./lineTile";
import { DonutTile, donutTileDataMock } from "./donutTile";

export const Tile = ({
  locations,
  tags,
  stat,
  start,
  end,
  count,
  version,
}: {
  locations: Location[];
  tags: string[];
  stat: string;
  start: Date;
  end: Date;
  count: number;
  version: number;
}) => {
  const [data, setData] = useState(0);
  const [spinning, setSpinning] = useState(false);

  const details = stats[stat];

  const fetchData = () => {
    setSpinning(true);
    fetchServer("/reports/stats/" + details.resource, {
      method: "POST",
      body: JSON.stringify({
        From: start,
        To: end,
        Locations: locations.map((v) => v.id),
        Tags: tags,
      }),
    })
      .then((res) => {
        if (res.status !== 200) {
          //alert("Something went wrong while fetching");
        }
        return res.json();
      })
      .then((res: any) => {
        setData(res.Data);
      })
      .finally(() => {
        setSpinning(false);
      });
  };

  useEffect(fetchData, [
    locations,
    tags,
    start,
    end,
    details.resource,
    version,
  ]);

  let tileComponent = <></>;

  switch (details.type) {
    case "number":
      tileComponent = (
        <NumberTile data={data} pre={details.pre} after={details.after} />
      );
      break;
    case "line":
      tileComponent = <LineTile data={lineTileDataMock} />;
      break;
    case "donut":
      tileComponent = <DonutTile data={donutTileDataMock} />;
      break;
    default:
      break;
  }

  return (
    <Grid xs={count === 1 ? 3 : 12}>
      <Box
        sx={{
          background: "rgba(0, 0, 0, .1)",
          margin: "3px",
          padding: "5px",
          paddingTop: "40px",
          paddingBottom: "40px",
          "&:hover": {
            border: "solid 1px",
            borderColor: "rgba(157, 61, 252, .5)",
            margin: "2px",
          },
        }}
      >
        <Typography
          variant="body1"
          sx={{
            textAlign: "center",
            height: "60px",
          }}
        >
          {details.title.toUpperCase()}

          <Tooltip title={details.info}>
            <InfoIcon
              sx={{
                fontSize: 14,
                opacity: "0.5",
                marginLeft: "10px",
              }}
            />
          </Tooltip>
        </Typography>
        <Typography
          variant="h2"
          sx={{
            textAlign: "center",
          }}
        >
          {spinning ? <CircularProgress /> : tileComponent}{" "}
        </Typography>
        <Button onClick={fetchData}>Refresh</Button>
      </Box>
    </Grid>
  );
};
