import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Checkbox,
  FormControl,
  FormControlLabel,
  InputLabel,
  MenuItem,
  Select,
  Typography,
  TextField,
  Input,
} from "@mui/material";
import { useState, useEffect } from "react";
import { useToast } from "use-toast-mui";
import { useFetchServer } from "../server/server";
import SurveyIframeGetter from "./surveyIframe";
import {
  LocationSettings,
  ZenotiCenter,
  ZenotiCenters,
  ZenotiService,
  ZenotiServices,
} from "./types";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { Location } from "../types/types";

export default function Settings({
  locationId,
}: {
  locationId: string | undefined;
}) {
  const [locationSettings, setLocationSettings] = useState<
    LocationSettings | undefined
  >(undefined);

  const [zenotiCenters, setZenotiCenters] = useState<ZenotiCenter[]>(
    [] as ZenotiCenter[]
  );

  const [zenotiServices, setZenotiServices] = useState<ZenotiService[]>(
    [] as ZenotiService[]
  );

  const fetchServer = useFetchServer();
  const toast = useToast();

  const save = () => {
    fetchServer("/settings/locations/" + locationId, {
      method: "POST",
      body: JSON.stringify(locationSettings?.location),
    }).then((res) => {
      toast.info("Settings are saved");
    });
  };

  const fetchZenotiCenters = () => {
    fetchServer<ZenotiCenters>(
      "/settings/zenoti-centers/" + locationSettings?.location.zenotiApi
    ).then((res) => {
      setZenotiCenters(res.centers);
      setLocationSettings(
        locationSettings
          ? {
              ...locationSettings,
              location: {
                ...locationSettings.location,
                zenotiUrl: res.url,
              },
            }
          : undefined
      );
    });
  };

  const fetchZenotiServices = () => {
    fetchServer<ZenotiService[]>(
      "/settings/zenoti-centers/" +
        locationSettings?.location.zenotiApi +
        "/" +
        locationSettings?.location.zenotiCenterId +
        "/services"
    ).then((res) => {
      setZenotiServices(res || ([] as ZenotiService[]));
    });
  };

  useEffect(() => {
    if (!locationId) return;
    fetchServer<LocationSettings>("/settings/locations/" + locationId).then(
      (res) => {
        setLocationSettings(res);
      }
    );
  }, [locationId]);

  if (!locationId || !locationSettings) {
    return (
      <Box
        sx={{
          background: "rgba(0, 0, 0, .1)",
          marginTop: "30px",
          margin: "auto",
          padding: "20px",
          border: "solid 1px",
          borderColor: "gray",
          width: "80%",
          borderRadius: "5px",
        }}
      >
        <Typography
          variant="body1"
          sx={{
            textAlign: "center",
          }}
        >
          "Select a location"
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        background: "rgba(0, 0, 0, .1)",
        marginTop: "30px",
        margin: "auto",
        padding: "20px",
        border: "solid 1px",
        borderColor: "gray",
        width: "80%",
        borderRadius: "5px",
      }}
    >
      <Typography
        variant="body1"
        sx={{
          textAlign: "center",
        }}
      >
        {locationSettings?.location?.name.toUpperCase()}
      </Typography>

      <FormControl fullWidth>
        <InputLabel id="pipelinel">Pipeline</InputLabel>
        <Select
          sx={{
            marginBottom: "20px",
          }}
          labelId="pipelinel"
          id="pipeline"
          value={locationSettings?.location.pipelineId}
          label="Pipeline"
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                pipelineId: e.target.value,
                noShowsId: "",
                bookId: "",
              },
            });
          }}
        >
          {locationSettings.pipelines.map((p) => {
            return <MenuItem value={p.id}>{p.name}</MenuItem>;
          })}
        </Select>

        <MySelect
          label="Bookings stage"
          value={locationSettings?.location.bookId}
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                bookId: e.target.value,
              },
            });
          }}
        >
          {locationSettings.pipelines
            .filter((p) => {
              return p.id === locationSettings.location.pipelineId;
            })[0]
            ?.stages.map((p) => {
              return <MenuItem value={p.id}>{p.name}</MenuItem>;
            })}
        </MySelect>

        <MySelect
          label="Sales stage"
          value={locationSettings?.location.salesId}
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                salesId: e.target.value,
              },
            });
          }}
        >
          {locationSettings.pipelines
            .filter((p) => {
              return p.id === locationSettings.location.pipelineId;
            })[0]
            ?.stages.map((p) => {
              return <MenuItem value={p.id}>{p.name}</MenuItem>;
            })}
        </MySelect>

        <MySelect
          label="NoShows stage"
          value={locationSettings?.location.noShowsId}
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                noShowsId: e.target.value,
              },
            });
          }}
        >
          {locationSettings.pipelines
            .filter((p) => {
              return p.id === locationSettings.location.pipelineId;
            })[0]
            ?.stages.map((p) => {
              return <MenuItem value={p.id}>{p.name}</MenuItem>;
            })}
        </MySelect>

        <MySelect
          label="Showed But Didn't Purchase stage"
          value={locationSettings?.location.showNoSaleId}
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                showNoSaleId: e.target.value,
              },
            });
          }}
        >
          {locationSettings.pipelines
            .filter((p) => {
              return p.id === locationSettings.location.pipelineId;
            })[0]
            ?.stages.map((p) => {
              return <MenuItem value={p.id}>{p.name}</MenuItem>;
            })}
        </MySelect>

        <MySelect
          label="Memberships Converted stage"
          value={locationSettings?.location.memberId}
          onChange={(e) => {
            setLocationSettings({
              ...locationSettings,
              location: {
                ...locationSettings.location,
                memberId: e.target.value,
              },
            });
          }}
        >
          {locationSettings.pipelines
            .filter((p) => {
              return p.id === locationSettings.location.pipelineId;
            })[0]
            ?.stages.map((p) => {
              return <MenuItem value={p.id}>{p.name}</MenuItem>;
            })}
        </MySelect>

        <FormControlLabel
          label="Alert if no leads in 24 hours"
          control={
            <Checkbox
              checked={locationSettings.location.trackNewLeads}
              onChange={(e) => {
                setLocationSettings({
                  ...locationSettings,
                  location: {
                    ...locationSettings.location,
                    trackNewLeads: e.target.checked,
                  },
                });
              }}
            />
          }
        />
      </FormControl>

      <Accordion
        sx={{
          background: "rgba(0, 0, 0, 0)",
          borderRadius: "4px",
          border: "solid 1px",
          borderColor: "rgba(255, 255, 255, .3)",
        }}
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>Zenoti integation</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <TextField
            label="Zenoti API key"
            value={locationSettings.location.zenotiApi}
            onChange={(e) => {
              setLocationSettings({
                ...locationSettings,
                location: {
                  ...locationSettings.location,
                  zenotiApi: e.target.value,
                },
              });
            }}
            sx={{
              width: "100%",
              marginBottom: "20px",
            }}
          />
          <FormControl fullWidth>
            <InputLabel id="zenotiCenter-l">Zenoti center</InputLabel>
            <Select
              sx={{
                marginBottom: "20px",
              }}
              labelId="zenotiCenter-l"
              id="zenotiCenter-i"
              value={locationSettings.location.zenotiCenterId}
              label="Zenoti center"
              onOpen={fetchZenotiCenters}
              onChange={(e) => {
                setLocationSettings({
                  ...locationSettings,
                  location: {
                    ...locationSettings.location,
                    zenotiCenterId: e.target.value,
                    zenotiCenterName: zenotiCenters.filter((c) => {
                      return c.id === e.target.value;
                    })[0].name,
                  },
                });
              }}
            >
              {zenotiCenters.length === 0 ? (
                <MenuItem value={locationSettings.location.zenotiCenterId}>
                  {locationSettings.location.zenotiCenterName}
                </MenuItem>
              ) : (
                zenotiCenters.map((c) => {
                  return <MenuItem value={c.id}>{c.name}</MenuItem>;
                })
              )}
            </Select>
            <InputLabel id="zenotiUrl-l">Zenoti center</InputLabel>
            <TextField
              sx={{
                marginBottom: "20px",
              }}
              value={locationSettings.location.zenotiUrl}
              label="Zenoti URL"
              onChange={(e) => {
                setLocationSettings({
                  ...locationSettings,
                  location: {
                    ...locationSettings.location,
                    zenotiUrl: e.target.value,
                  },
                });
              }}
            />
          </FormControl>
          <FormControlLabel
            label="Sync contacts and stages"
            control={
              <Checkbox
                checked={locationSettings.location.syncContacts}
                onChange={(e) => {
                  setLocationSettings({
                    ...locationSettings,
                    location: {
                      ...locationSettings.location,
                      syncContacts: e.target.checked,
                    },
                  });
                }}
              />
            }
          />
          <FormControlLabel
            label="Automatically create contacts in Zenoti"
            control={
              <Checkbox
                checked={locationSettings.location.autoCreateContacts}
                onChange={(e) => {
                  setLocationSettings({
                    ...locationSettings,
                    location: {
                      ...locationSettings.location,
                      autoCreateContacts: e.target.checked,
                    },
                  });
                }}
              />
            }
          />
          <Box flexDirection={"row"}>
            <FormControlLabel
              label="Sync calendars"
              control={
                <Checkbox
                  checked={locationSettings.location.syncCalendars}
                  onChange={(e) => {
                    setLocationSettings({
                      ...locationSettings,
                      location: {
                        ...locationSettings.location,
                        syncCalendars: e.target.checked,
                      },
                    });
                  }}
                />
              }
            />
            <FormControl
              sx={{
                width: "60%",
                marginLeft: "20px",
              }}
            >
              <InputLabel id="zenotiSerivice-l">Zenoti service</InputLabel>
              <Select
                labelId="zenotiSerivice-l"
                id="zenotiService-i"
                value={locationSettings.location.zenotiServiceId}
                label="Zenoti service"
                onOpen={fetchZenotiServices}
                onChange={(e) => {
                  setLocationSettings({
                    ...locationSettings,
                    location: {
                      ...locationSettings.location,
                      zenotiServiceId: e.target.value,
                      zenotiServiceName: zenotiServices.filter((c) => {
                        return c.id === e.target.value;
                      })[0].name,
                    },
                  });
                }}
              >
                {zenotiServices.length === 0 ? (
                  <MenuItem value={locationSettings.location.zenotiServiceId}>
                    {locationSettings.location.zenotiServiceName}
                  </MenuItem>
                ) : (
                  zenotiServices.map((c) => {
                    return <MenuItem value={c.id}>{c.name}</MenuItem>;
                  })
                )}
              </Select>
            </FormControl>
            <TextField
              label="Price"
              type="number"
              value={locationSettings.location.zenotiServicePrice}
              onChange={(e) => {
                setLocationSettings({
                  ...locationSettings,
                  location: {
                    ...locationSettings.location,
                    zenotiServicePrice: e.target.value
                      ? parseInt(e.target.value)
                      : 0,
                  },
                });
              }}
              sx={{
                width: "10%",
                marginLeft: "20px",
              }}
            />
          </Box>
        </AccordionDetails>
      </Accordion>

      <Button onClick={save}>Save location</Button>
      <SurveyIframeGetter location={locationSettings} />
    </Box>
  );
}

const MySelect = ({
  label,
  value,
  onChange,
  children,
}: {
  label: string;
  value: string;
  onChange: (arg0: any) => any;
  children: JSX.Element | JSX.Element[];
}) => {
  return (
    <FormControl fullWidth>
      <InputLabel id={label + "l"}>{label}</InputLabel>
      <Select
        sx={{
          marginBottom: "20px",
        }}
        labelId={label + "l"}
        id={label + "i"}
        value={value}
        label={label}
        onChange={onChange}
      >
        {children}
      </Select>
    </FormControl>
  );
};
