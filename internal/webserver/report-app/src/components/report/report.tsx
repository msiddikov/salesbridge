import { useEffect, useState } from "react";
import { Grid, Typography } from "@mui/material";
import { Box } from "@mui/system";
import { useFetchServer } from "../server/server";
import {
  Settings,
  SettingsData,
  reportTypeFacebook,
  settingsDataDef,
  settingsDef,
} from "../types/types";
import { Location } from "../types/types";

import { CommonSettings } from "./commonSettings";
import JfmReport from "./jfm_report/jfm_report";
import SimpleReport from "./tile/simpleReport";
import { facebookStats } from "./facebook_report/facebook_report";

function Report() {
  const [settings, setSettings] = useState(settingsDef as Settings);
  const [settingsData, setSettingsData] = useState(settingsDataDef);

  const fetchServer = useFetchServer();

  const updateLocations = () => {
    fetchServer<SettingsData>("/reports/settings").then((res) => {
      setSettingsData(res);
      readLocation(res.locations);
    });
  };

  useEffect(updateLocations, []);

  const readLocation = (locations: Location[]) => {
    const queryParams = new URLSearchParams(window.location.search);
    let id = queryParams.get("location") || "";
    if (id === "jfm") {
      setSettings({ ...settings, agencyMode: true });
      return;
    }
    setSettings({
      ...settings,
      items: settings.items.map((v) => {
        v.locations = locations.filter((l) => l.id === id);
        return v;
      }),
    });
  };

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        overflowY: "scroll",
      }}
    >
      <Heading
        agencyMode={settings.agencyMode}
        locationName={settings.items[0].locations[0]?.name || ""}
      />
      <CommonSettings
        settings={settings}
        setSettings={setSettings}
        settingsData={settingsData}
      />

      <Box
        sx={{
          width: "100%",
          textAlign: "center",
          fontSize: "15pt",
          minHeight: "90%",
        }}
      >
        <Grid container columns={{ xs: 4, sm: 8, md: 12 }}>
          {settings.items.map((v, i) => {
            let reportProps = {
              locations: v.locations,
              start: v.range[0].startDate,
              end: v.range[0].endDate,
              count: settings.items.length,
              settingsData: settingsData,
              stats: [] as string[],
              tags: [] as string[],
            };
            switch (settings.type) {
              case reportTypeFacebook:
                reportProps.stats = facebookStats;
                return <SimpleReport {...reportProps} />;
              default:
                return <JfmReport {...reportProps} />;
            }
          })}
        </Grid>
      </Box>
      <Typography
        sx={{
          fontSize: "10pt",
          textAlign: "center",
        }}
      >
        {" "}
        Jump Forward Media 2022 ©{" "}
      </Typography>
    </Box>
  );
}

export default Report;

function Heading({
  agencyMode,
  locationName,
}: {
  agencyMode: boolean;
  locationName: string;
}) {
  return (
    <Box
      sx={{
        width: "100%",
        textAlign: "center",
        fontSize: "15pt",
      }}
    >
      <Typography
        variant="h2"
        sx={{
          margin: "20px",
        }}
      >
        {agencyMode ? "Jump Forward Media" : locationName}
      </Typography>
    </Box>
  );
}
