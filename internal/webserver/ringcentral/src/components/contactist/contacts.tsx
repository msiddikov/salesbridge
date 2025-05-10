import {
  Center,
  Container,
  Flex,
  Icon,
  Input,
  InputGroup,
  InputLeftElement,
  InputRightElement,
  Text,
} from "@chakra-ui/react";
import { ChatInfo, ContactListProps, ContactsLookupRes } from "../types/types";
import React, { useEffect, useState } from "react";
import { host } from "../consts";
import { useElementSize } from "usehooks-ts";
import { IoCloseCircleSharp, IoSearchOutline } from "react-icons/io5";
import { RoundAvatar } from "../common";

function ContactList({ chatInfo, setChat }: ContactListProps) {
  const [predictions, setPredictions] = useState<ChatInfo[]>([]);
  const [inputVal, setInput] = useState("");
  const [inputRef] = useElementSize();
  const [chats, setChats] = useState<ChatInfo[]>([]);

  useEffect(() => {
    if (chatInfo.locationId === undefined) {
      return;
    }
    fetch(host + "/rc/chats/" + chatInfo.locationId)
      .then((res) => res.json())
      .then((res: ChatInfo[]) => {
        if (
          !res.find((el) => el.contactId === chatInfo.contactId) &&
          chatInfo.contactId !== undefined
        ) {
          res = [chatInfo, ...res];
        }
        setChats(res.sort((a, b) => (a.date > b.date ? -1 : 1)));
      });
  }, [chatInfo]);

  const updatePredictions = (event: React.ChangeEvent) => {
    // @ts-ignore
    setInput(event.target.value);
    fetch(
      host +
        "/rc/contacts?location=" +
        chatInfo.locationId +
        "&query=" +
        // @ts-ignore
        event.target.value
    )
      .then((res) => res.json())
      .then((res: ContactsLookupRes[]) =>
        setPredictions(
          res.map((v) => {
            const ch: ChatInfo = {
              contactName: v.firstName + " " + v.lastName,
              locationId: chatInfo.locationId,
              contactId: v.id,
            } as ChatInfo;
            return ch;
          })
        )
      );
  };

  const onPredictionSelect = (event: React.MouseEvent) => {
    setInput("");
    let chat = chats.find((el) => {
      return el.contactId === event.currentTarget.id;
    });
    if (!chat) {
      const pred = predictions.find((el) => {
        return el.contactId === event.currentTarget.id;
      });

      setChat({
        locationId: chatInfo.locationId,
        contactId: pred?.contactId,
        contactName: pred?.contactName,
        chatId: 0,
      } as ChatInfo);
    } else {
      setChat(chat);
    }
    setPredictions([]);
  };

  const onChatSelect = (event: React.MouseEvent) => {
    let chatId = +event.currentTarget.id;

    let chat = chats.find((el) => {
      return el.chatId === chatId;
    });
    if (!chat) return;
    setChat(chat);
  };

  const Contact = (props: {
    item: ChatInfo;
    selected: boolean;
    isPrediction: boolean;
  }) => {
    const name = props.item.contactName;

    return (
      <Container
        border="1px"
        borderColor="gray.200"
        backgroundColor={props.selected ? "green.50" : "white"}
        w="100%"
        p="0"
        cursor="pointer"
        id={props.isPrediction ? props.item.contactId : "" + props.item.chatId}
        onClick={props.isPrediction ? onPredictionSelect : onChatSelect}
        marginTop="10px"
        borderRadius="10px"
      >
        <Flex
          alignItems="stretch"
          direction="row"
          w="full"
          h="70px"
          p="5px"
          id={name}
        >
          <Center w="45px" p="0" m="5px">
            <RoundAvatar name={name} width={40} />
          </Center>

          <Container
            w="full"
            overflow="hidden"
            p="0"
            textOverflow="ellipsis"
            paddingTop="8px"
          >
            <Text noOfLines={1}>
              <b>{name}</b>
            </Text>
            <Text noOfLines={1} fontSize={11}>
              {props.item.lastMessage}
            </Text>
          </Container>
        </Flex>
      </Container>
    );
  };
  return (
    <Flex alignItems="stretch" w="100%" h="full" direction="column">
      <Flex p="5px" backgroundColor="gray.50" direction="column">
        <InputGroup ref={inputRef}>
          <InputLeftElement
            pointerEvents="none"
            children={<Icon as={IoSearchOutline} />}
          />
          <Input
            onChange={updatePredictions}
            value={inputVal}
            type="tel"
            placeholder="Contact name or email"
            backgroundColor="white"
            borderLeft="none"
          />
          <InputRightElement
            cursor="pointer"
            children={
              <Icon
                onClick={() => {
                  setInput("");
                }}
                as={IoCloseCircleSharp}
              />
            }
          />
        </InputGroup>
      </Flex>
      <Flex
        alignItems="stretch"
        w="100%"
        h="full"
        direction="column"
        overflow="hidden"
        backgroundColor="gray.50"
        margin="auto"
        className="contacts-scroll"
        paddingLeft="8px"
        paddingRight="8px"
        _hover={{
          overflowY: "scroll",
          paddingRight: "5px",
        }}
      >
        {(inputVal === "" ? chats : predictions).map((el) => (
          <Contact
            key={el.contactId}
            selected={el.chatId === chatInfo.chatId}
            item={el}
            isPrediction={inputVal !== ""}
          ></Contact>
        ))}
      </Flex>
    </Flex>
  );
}

export default ContactList;

// const testList: contact[] = [
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked you some awesome features...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked you a new lkjsdfoiu...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
//   {
//     Name: "Carla Geleer",
//     LastMsg: "Hi, we have booked yo...",
//     Type: "mail",
//     Url: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
//   },
// ];
