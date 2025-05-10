import {
  SelectChangeEvent,
  Grid,
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
} from "@mui/material";
import React, { useState } from "react";

export {};
export const TagsSelect = ({
  selected,
  setSelected,
  tags,
}: {
  selected: string[];
  setSelected: (tags: string[]) => void;
  tags: string[];
}) => {
  const [selecteds, setSelecteds] = React.useState<string[]>(selected);
  const [open, setOpen] = useState(false);

  const handleChange = (event: SelectChangeEvent<string[]>) => {
    const {
      target: { value },
    } = event;
    const tags = typeof value === "string" ? value.split(",") : value;
    setSelecteds(tags);
  };

  return (
    <Grid xs={12}>
      <FormControl sx={{ m: 1, width: "80%" }}>
        <InputLabel id="demo-multiple-checkbox-label">Tags</InputLabel>
        <Select
          open={open}
          onOpen={() => {
            setOpen(true);
          }}
          onClose={() => {
            setOpen(false);
          }}
          multiple
          value={selecteds.map((t) => t)}
          onChange={handleChange}
          input={<OutlinedInput label="Tags" />}
          renderValue={(selected) => (
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
              {selected.map((value) => (
                <Chip
                  key={value}
                  label={selecteds.filter((t) => t === value)[0]}
                />
              ))}
            </Box>
          )}
        >
          <Button onClick={() => setSelecteds(tags)}>Select all</Button>
          <Button onClick={() => setSelecteds([] as string[])}>
            Deselect all
          </Button>
          {tags.map((t: string) => (
            <MenuItem key={t} value={t}>
              <Checkbox checked={selecteds.indexOf(t) > -1} />
              <ListItemText primary={t} />
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
