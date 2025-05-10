import moment from "moment";
import { type } from "os";

export type Range = {
  startDate: Date;
  endDate: Date;
  key: string;
};

export type Location = {
  name: string;
  id: string;
  isIntegrated: boolean;
};

export type Setting = {
  locations: Location[];
  tags: string[];
  range: Range[];
  type: string;
};

export type Settings = {
  items: Setting[];
  agencyMode: boolean;
  type: string;
};

export type SettingsData = {
  locations: Location[];
  tags: string[];
};

export type ReportProps = {
  stats: string[];
  locations: Location[];
  tags: string[];
  start: Date;
  end: Date;
  settingsData: SettingsData;
  count: number;
};

export const reportTypeJFM = "jfm";
export const reportTypeFacebook = "facebook";

export const rangeDef = {
  startDate: moment().subtract("months", 1).toDate(),
  endDate: new Date(),
  key: "selection",
};

export const settingsItemDef = {
  locations: [],
  tags: [],
  range: [rangeDef],
  type: reportTypeJFM,
} as Setting;

export const settingsDef = {
  items: [settingsItemDef],
  agencyMode: false,
  type: reportTypeJFM,
} as Settings;

export const settingsDataDef = {
  locations: [],
  tags: [],
} as SettingsData;
