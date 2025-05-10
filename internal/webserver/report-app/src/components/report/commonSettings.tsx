import {
  Grid,
  Typography,
  Popover,
  SelectChangeEvent,
  FormControl,
  InputLabel,
  Select,
  OutlinedInput,
  Box,
  Chip,
  Button,
  MenuItem,
  Checkbox,
  ListItemText,
  Container,
  ToggleButtonGroup,
  ToggleButton,
} from "@mui/material";
import React, { useState } from "react";
import {
  Location,
  Range,
  Settings,
  SettingsData,
  settingsItemDef,
} from "../types/types";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { DateRangePicker } from "react-date-range";

export function CommonSettings({
  settings,
  setSettings,
  settingsData,
}: {
  settings: Settings;
  setSettings: (settings: Settings) => void;
  settingsData: SettingsData;
}) {
  return (
    <Box
      sx={{
        width: "100%",
        textAlign: "center",
        fontSize: "15pt",
      }}
    >
      <ToggleButtonGroup
        value={settings.type}
        exclusive
        onChange={(event, value) => {
          setSettings({ ...settings, type: value });
        }}
      >
        <ToggleButton value="jfm">JFM Report</ToggleButton>
        <ToggleButton value="facebook">Facebook Report</ToggleButton>
      </ToggleButtonGroup>
      <CompareButton settings={settings} setSettings={setSettings} />
      <Grid container columns={{ xs: 4, sm: 8, md: 12 }}>
        {settings.items.map((v, i) => {
          return (
            <Grid xs={12 / settings.items.length}>
              <LocationsSelect
                selected={v.locations}
                setSelected={(loc: Location[]) => {
                  setSettings({
                    ...settings,
                    items: settings.items.map((e, k) => {
                      return k === i ? { ...e, locations: loc } : e;
                    }),
                  });
                }}
                locations={settingsData.locations}
                agencyMode={settings.agencyMode}
              />
              <DatePicker
                range={v.range}
                setRange={(range: Range[]) => {
                  setSettings({
                    ...settings,
                    items: settings.items.map((e, k) => {
                      return k === i ? { ...e, range: range } : e;
                    }),
                  });
                }}
              />
            </Grid>
          );
        })}
      </Grid>
    </Box>
  );
}

function CompareButton({
  settings,
  setSettings,
}: {
  settings: Settings;
  setSettings: (settings: Settings) => void;
}) {
  return (
    <Box>
      Compare
      <Button
        variant="outlined"
        sx={{
          margin: "10px",
        }}
        onClick={() => {
          if (settings.items.length > 1) {
            setSettings({ ...settings, items: settings.items.slice(0, -1) });
          }
        }}
      >
        -
      </Button>
      {settings.items.length}
      <Button
        variant="outlined"
        sx={{
          margin: "10px",
        }}
        onClick={() => {
          if (settings.items.length < 4) {
            setSettings({
              ...settings,
              items: [...settings.items, settingsItemDef],
            });
          }
        }}
      >
        +
      </Button>
    </Box>
  );
}

const DatePicker = ({
  range,
  setRange,
}: {
  range: Range[];
  setRange: (range: Range[]) => void;
}) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const id = open ? "simple-popover" : undefined;

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <Grid xs={12}>
      <Typography
        variant="body1"
        sx={{
          cursor: "pointer",
          opacity: ".8",
          textAlign: "center",
          marginBottom: "20px",
        }}
        onClick={handleClick}
      >
        Period:{" "}
        {range[0].startDate.toLocaleDateString() +
          " - " +
          range[0].endDate.toLocaleDateString()}
        <KeyboardArrowDownIcon
          sx={{
            marginBottom: "-5px",
          }}
        />
      </Typography>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "left",
        }}
      >
        <DateRangePicker
          onChange={(item: any) => setRange([item.selection])}
          moveRangeOnFirstSelection={false}
          months={2}
          ranges={range}
          direction="vertical"
        />
      </Popover>
    </Grid>
  );
};

const LocationsSelect = ({
  selected,
  locations,
  setSelected,
  agencyMode,
}: {
  selected: Location[];
  locations: Location[];
  setSelected: (range: Location[]) => void;
  agencyMode: boolean;
}) => {
  const [selecteds, setSelecteds] = React.useState<Location[]>(selected);
  const [open, setOpen] = useState(false);

  const handleChange = (event: SelectChangeEvent<string[]>) => {
    const {
      target: { value },
    } = event;
    const ids = typeof value === "string" ? value.split(",") : value;
    setSelecteds(
      locations.filter((l) => {
        return ids.indexOf(l.id) > -1;
      })
    );
  };

  if (!agencyMode) {
    return null;
  }
  return (
    <Grid xs={12}>
      <FormControl sx={{ m: 1, width: "80%" }}>
        <InputLabel id="demo-multiple-checkbox-label">Locations</InputLabel>
        <Select
          multiple
          open={open}
          onOpen={() => {
            setOpen(true);
          }}
          onClose={() => {
            setOpen(false);
          }}
          value={selecteds.map((l) => l.id)}
          onChange={handleChange}
          input={<OutlinedInput label="Locations" />}
          renderValue={(selected) => (
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected.map((value) => (
                <Chip
                  key={value}
                  label={selecteds.filter((l) => l.id === value)[0].name}
                />
              ))}
            </Box>
          )}
        >
          <Button onClick={() => setSelecteds(locations)}>Select all</Button>
          <Button onClick={() => setSelecteds([] as Location[])}>
            Deselect all
          </Button>
          {locations.map((l: Location) => (
            <MenuItem key={l.name} value={l.id}>
              <Checkbox checked={selecteds.indexOf(l) > -1} />
              <ListItemText primary={l.name} />
            </MenuItem>
          ))}
          <Container
            sx={{
              width: "100%",
              alignContent: "center",
            }}
          >
            <Button
              sx={{
                margin: "auto",
              }}
              onClick={() => {
                setSelected(selecteds);
                setOpen(false);
              }}
            >
              OK
            </Button>
          </Container>
        </Select>
      </FormControl>
    </Grid>
  );
};
