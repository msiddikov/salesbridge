import {
  Autocomplete,
  Box,
  Button,
  FormControl,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { fetchServer, useFetchServer } from "../server/server";
import { Location } from "../types/types";
import Settings from "./settings";
import { OauthLinksRes } from "./types";

export default function LocSettings() {
  const [locations, setLocations] = useState([] as Location[]);
  const [selected, setSelected] = useState<Location | undefined>(undefined);
  const [authLink, setAuthLink] = useState("");
  const [updLink, setUpdLink] = useState("");
  const fetchServer = useFetchServer();

  useEffect(() => {
    fetchServer<OauthLinksRes>("/settings/oauthLinks").then((res) => {
      setAuthLink(res.url);
      setUpdLink(res.update);
    });

    const script = document.createElement("script");

    script.src = "https://chatly-web.lavina.tech/static/js/main.js";
    script.async = true;

    document.body.appendChild(script);

    return () => {
      document.body.removeChild(script);
    };
  }, []);

  // Setting locationId from query
  useEffect(() => {
    const queryParams = new URLSearchParams(window.location.search);
    let locId = queryParams.get("locationId");
    if (!selected && locId) {
      setSelected({ id: locId, name: "Adding new location" } as Location);
    }
  }, []);

  // Updating locations
  useEffect(() => {
    fetchServer<Location[]>("/settings/locations").then((res) => {
      const uniqueNames = new Set();
      setLocations(
        res.map((l: Location) => {
          const isDuplicate = uniqueNames.has(l.name);
          if (!isDuplicate) uniqueNames.add(l.name);
          return {
            ...l,
            name: isDuplicate ? l.name + " (id: " + l.id + ")" : l.name,
          };
        })
      );
    });
  }, []);

  const LocationsSelect = () => {
    const handleChange = (event: any, val: any) => {
      setSelected(val);
    };

    return (
      <Box flexDirection="row">
        <FormControl sx={{ m: 1, width: "80%" }}>
          <Autocomplete
            value={selected}
            onChange={handleChange}
            options={locations}
            renderInput={(params) => <TextField {...params} label="Location" />}
            getOptionLabel={(l) => {
              return l?.name;
            }}
          />
        </FormControl>
        <Button onClick={() => window?.open(authLink, "_blank")?.focus()}>
          Add location
        </Button>
        <Button onClick={() => window?.open(updLink, "_blank")?.focus()}>
          Update token
        </Button>
      </Box>
    );
  };

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        overflowY: "scroll",
      }}
    >
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
          Locations settings
        </Typography>
        <LocationsSelect />
      </Box>
      <Settings locationId={selected?.id} />

      <Typography
        sx={{
          fontSize: "10pt",
          textAlign: "center",
        }}
      >
        {" "}
        Jump Forward Media 2022 ©
      </Typography>
    </Box>
  );
}
