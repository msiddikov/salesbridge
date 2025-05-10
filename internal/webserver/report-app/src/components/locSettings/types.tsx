export type LocationSettings = {
  location: {
    name: string;
    id: string;
    pipelineId: string;
    bookId: string;
    salesId: string;
    noShowsId: string;
    showNoSaleId: string;
    memberId: string;
    trackNewLeads: boolean;
    zenotiApi: string;
    zenotiUrl: string;
    zenotiCenterId: string;
    zenotiCenterName: string;
    zenotiServiceId: string;
    zenotiServiceName: string;
    zenotiServicePrice: number;
    syncCalendars: boolean;
    syncContacts: boolean;
    autoCreateContacts: boolean;
  };

  pipelines: {
    id: string;
    name: string;
    stages: {
      name: string;
      id: string;
    }[];
  }[];
  workflows: Workflow[];
};

export type OauthLinksRes = {
  url: string;
  update: string;
};

export type Workflow = {
  name: string;
  id: string;
};

export type Survey = {
  id: string;
  label: string;
};

export type ZenotiCenter = {
  id: string;
  name: string;
};

export type ZenotiService = {
  id: string;
  name: string;
};

export type ZenotiCenters = {
  centers: ZenotiCenter[];
  url: string;
};

export type ZenotiServices = {
  data: ZenotiService[];
};
