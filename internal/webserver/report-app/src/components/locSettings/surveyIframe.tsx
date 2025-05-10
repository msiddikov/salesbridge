import {
  Autocomplete,
  Box,
  Button,
  FormControlLabel,
  InputAdornment,
  Modal,
  Switch,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useToast } from "use-toast-mui";
import { host, useFetchServer } from "../server/server";
import { LocationSettings, Survey, Workflow } from "./types";

export default function SurveyIframeGetter({
  location,
}: {
  location: LocationSettings;
}) {
  const [open, setOpen] = useState(false);
  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);
  const [workflow, setWorkflow] = useState<Workflow | null>(null);
  const [scale, setScale] = useState("100");
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [addToWorkflow, setAddToWorkflow] = useState(false);
  const [afterSubmitURL, setAfterSubmitURL] = useState("");
  const toast = useToast();

  const surveys = [
    {
      id: "hairRemoval",
      label: "Hair removal",
    },
    {
      id: "weightLoss",
      label: "Weight loss",
    },
  ] as Survey[];

  const style = {
    position: "absolute" as "absolute",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 600,
    bgcolor: "background.paper",
    border: "2px solid #000",
    boxShadow: 24,
    p: 4,
    "& > :not(style)": { m: 2 },
  };

  return (
    <div>
      <Button onClick={handleOpen}>Survey app</Button>
      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={style}>
          <Typography id="modal-modal-title" variant="h6" component="h1">
            Configure and publish survey app
          </Typography>
          <TextField
            type="number"
            label="Scale"
            value={scale}
            onChange={(e) => {
              setScale(e.target.value);
            }}
            InputProps={{
              endAdornment: <InputAdornment position="end">%</InputAdornment>,
            }}
          />

          <TextField
            label="URL after survey submission"
            value={afterSubmitURL}
            onChange={(e) => {
              setAfterSubmitURL(e.target.value);
            }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">https://</InputAdornment>
              ),
            }}
          />

          <Autocomplete
            disablePortal
            value={survey}
            onChange={(e, val) => setSurvey(val)}
            options={surveys}
            sx={{ width: 300 }}
            renderInput={(params) => (
              <TextField {...params} label="Select survey" />
            )}
          />
          <FormControlLabel
            control={
              <Switch
                checked={addToWorkflow}
                onChange={(e) => {
                  if (!e.target.checked) {
                    setWorkflow(null);
                  }
                  setAddToWorkflow(e.target.checked);
                }}
              />
            }
            label="Add contact to workflow after submission"
          />

          {addToWorkflow ? (
            <Autocomplete
              disablePortal
              value={workflow}
              onChange={(e, val) => setWorkflow(val)}
              options={location.workflows}
              sx={{ width: 300 }}
              renderInput={(params) => (
                <TextField {...params} label="Select workflow" />
              )}
              getOptionLabel={(w) => w.name}
            />
          ) : null}
          <Button
            onClick={() => {
              if (!survey) {
                toast.error("Please select survey first");
                return;
              }
              if (addToWorkflow && !workflow) {
                toast.error("Please select workflow");
                return;
              }
              navigator.clipboard.writeText(
                getIframeCode(
                  location.location.id,
                  workflow?.id,
                  survey?.id,
                  +scale / 100,
                  "https://" + afterSubmitURL
                )
              );
              toast.success("IFrame code has been copied");
            }}
          >
            Copy iframe code
          </Button>
        </Box>
      </Modal>
    </div>
  );
}

export const getIframeCode = (
  locationId: string,
  workflowId: string | undefined,
  surveyId: string | undefined,
  scale: number,
  url: string
) => {
  let iframeH = 750;

  return (
    `
    <iframe id="scaled-frame" src="` +
    host +
    `/apps/survey/?locationId=` +
    locationId +
    `&workflowId=` +
    workflowId +
    `&surveyId=` +
    surveyId +
    `&url=` +
    url +
    `" scrolling="no"></iframe>
    <style>
    #scaled-frame {
      width: 100%;
      height: ` +
    iframeH +
    `px;
      border: 0px;
    }
    #scaled-frame {
      zoom: ` +
    scale +
    `;
      -moz-transform: scale(` +
    scale +
    `);
      -moz-transform-origin: 0 0;
      -o-transform: scale(` +
    scale +
    `);
      -o-transform-origin: 0 0;
      -webkit-transform: scale(` +
    scale +
    `);
      -webkit-transform-origin: 0 0;
    }
    @media screen and (-webkit-min-device-pixel-ratio:0) {
      #scaled-frame {
        zoom: 1;
      }
    }
    </style>`
  );
};
