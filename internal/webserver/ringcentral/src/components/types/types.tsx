export type ChatInfo = {
  locationId: string;
  chatId: number;
  contactName: string;
  contactId: string;
  lastMessage: string;
  date: string;
  url: string;
};

export const EmptyChatInfo = {
  locationId: "",
  contactId: "",
  name: "",
};

export type ContactListProps = {
  chatInfo: ChatInfo;
  setChat: React.Dispatch<React.SetStateAction<ChatInfo>>;
};

export type ContactInfoProps = {
  chatInfo: ChatInfo;
};

export type ChatMessageProps = {
  msg: Message;
  marginTop: boolean;
};

export type ContactsLookupRes = {
  id: string;
  email: string;
  phone: string;
  firstName: string;
  lastName: string;
};
export type ContactInfoRes = {
  email: string;
  firstName: string;
  lastName: string;
  phone: string;
  Url: string;
};

export const EmptyContactInfo = {
  email: "",
  firstName: "",
  lastName: "",
  phone: "",
  Url: "",
};

export type Message = {
  date: Date;
  messageId: string;
  text: string;
  managerName: string;
  inbound: boolean;
  Url: string;
};
