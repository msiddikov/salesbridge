type Stat = {
  title: string;
  info: string;
  resource: string;
  pre: string;
  after: string;
  type: string;
};

const statTypeNumber = "number";
const statTypeLine = "line";
const statTypeDonut = "donut";

export const stats = {
  // JFM report type stats

  expenses: {
    title: "Total expenses",
    info: "Total expenses for the selected period",
    resource: "Expenses",
    pre: "$",
    after: "",
    type: statTypeNumber,
  },
  sales: {
    title: "Total sales",
    info: "Total sales for the selected period",
    resource: "Sales",
    pre: "$",
    after: "",
    type: statTypeNumber,
  },
  salesNo: {
    title: "Number of clients sold",
    info: "Number of clients sold for the selected period",
    resource: "SalesNo",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  roi: {
    title: "Return on investment",
    info: "ROI for the selected period",
    resource: "ROI",
    pre: "$",
    after: "",
    type: statTypeNumber,
  },
  newLeads: {
    title: "New leads",
    info: "New leads in the selected period",
    resource: "NewLeads",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  bookings: {
    title: "Bookings",
    info: "Bookings for the selected period",
    resource: "Bookings",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  noShows: {
    title: "No Shows",
    info: "No shows for the selected period",
    resource: "NoShows",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  showNoSale: {
    title: "Showed but didn't purchase",
    info: "Showed but didn't purchase for the selected period",
    resource: "ShowNoSale",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  leadsConv: {
    title: "lead to booking conversion rate",
    info: "Leads conversion for the selected period",
    resource: "LeadsConv",
    pre: "",
    after: "%",
    type: statTypeNumber,
  },
  bookingsConv: {
    title: "Booking to sales conversion rate",
    info: "Bookings conversion for the selected period",
    resource: "BookingsConv",
    pre: "",
    after: "%",
    type: statTypeNumber,
  },
  showRate: {
    title: "Show Up Rate",
    info: "Show up rate for the selected period",
    resource: "ShowRate",
    pre: "",
    after: "%",
    type: statTypeNumber,
  },
  membershipConv: {
    title: "Membership conversion rate",
    info: "Membership conversion rate for the selected period",
    resource: "MembershipConv",
    pre: "",
    after: "%",
    type: statTypeNumber,
  },

  rank: {
    title: "Rank",
    info: "Location's rank by sales for the last 30 days",
    resource: "Rank",
    pre: "#",
    after: "",
    type: statTypeNumber,
  },

  zenotiMembersNo: {
    title: "Active members",
    info: "Location's active members number",
    resource: "ZenotiMembersNo",
    pre: "",
    after: "",
    type: statTypeNumber,
  },

  // Facebook report type stats

  impressions: {
    title: "Impressions",
    info: "Total impressions for the selected period",
    resource: "SalesNo",
    pre: "",
    after: "",
    type: statTypeNumber,
  },
  impressionsLine: {
    title: "Impressions Line",
    info: "Total impressions for the selected period",
    resource: "SalesNo",
    pre: "",
    after: "",
    type: statTypeLine,
  },
  impressionsDonut: {
    title: "Impressions Donut",
    info: "Total impressions for the selected period",
    resource: "SalesNo",
    pre: "",
    after: "",
    type: statTypeDonut,
  },
} as {
  [index: string]: Stat;
};
