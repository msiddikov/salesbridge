import { ArrowForwardIcon } from "@chakra-ui/icons";
import {
  Box,
  Flex,
  Input,
  InputGroup,
  InputRightAddon,
  Spinner,
  Text,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { host } from "../consts";

import { ChatInfo, Message } from "../types/types";

function Chat({
  chatInfo,
  managerName,
  setChatInfo,
}: {
  chatInfo: ChatInfo;
  managerName: string;
  setChatInfo: any;
}) {
  const [messages, setMessages] = useState<Message[]>([]);

  const updateMsg = () => {
    if (chatInfo.chatId === undefined) {
      return;
    }
    fetch(host + "/rc/messages/" + chatInfo.chatId)
      .then((res) => res.json())
      .then((res: Message[]) => {
        setMessages(
          res
            .map((el) => {
              el.date = new Date(el.date);
              return el;
            })
            .sort((a, b) => {
              return a.date > b.date ? -1 : 1;
            })
        );
      });
  };
  useEffect(() => {
    updateMsg();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chatInfo]);

  const ChatMessageInput = () => {
    const [message, setMessage] = useState("");
    const [spinning, setSpinning] = useState(false);
    const [updateSpinning, setUpdateSpinning] = useState(false);

    const onMessageChange = (event: React.ChangeEvent) => {
      // @ts-ignore
      setMessage(event.target.value);
    };

    const sendMessage = () => {
      setSpinning(true);
      fetch(host + "/rc/message", {
        method: "POST",
        body: JSON.stringify({
          locationId: chatInfo.locationId,
          contactId: chatInfo.contactId,
          text: message,
          managerName: managerName,
        }),
      })
        .then((res) => {
          if (res.status !== 200) {
            alert("Something went wrong while sending SMS");
          }
          return res.json();
        })
        .then((res: any) => {
          updateMsg();
          setSpinning(false);
          if (chatInfo.chatId === 0) {
            setChatInfo({ ...chatInfo, chatId: res.ID });
          }
        });
    };

    const updateSync = () => {
      setUpdateSpinning(true);
      fetch(host + "/rc/update").then(() => {
        updateMsg();
        setUpdateSpinning(false);
      });
    };

    const onMessageKeyPress = (e: any) => {
      if (e.key === "Enter") {
        sendMessage();
      }
    };
    return (
      <Flex
        p="15px"
        direction="column"
        backgroundColor="gray.50"
        borderTop="1px"
        borderColor="gray.100"
      >
        <InputGroup>
          <Input
            value={message}
            onChange={onMessageChange}
            onKeyDown={onMessageKeyPress}
            type="tel"
            placeholder="Type in your message"
            backgroundColor="white"
          />
          <InputRightAddon
            onClick={spinning ? undefined : sendMessage}
            cursor="pointer"
          >
            {spinning ? <Spinner /> : <ArrowForwardIcon />}
          </InputRightAddon>
        </InputGroup>
        <Text p="4px" noOfLines={1} fontSize={11}>
          {"Sending messages as '" + managerName + "'"}
        </Text>
        <Text noOfLines={1} fontSize={11} cursor="pointer" onClick={updateSync}>
          Update chat{"    "}
          {updateSpinning ? <Spinner size="xs" /> : ""}
        </Text>
      </Flex>
    );
  };

  const ChatMessage = ({ msg }: { msg: Message }) => {
    return (
      <Box
        alignSelf={msg.inbound ? "flex-begin" : "flex-end"}
        // background={"-webkit-linear-gradient("
        //   .concat(msg.inbound ? "right" : "left")
        //   .concat(", #e3fff0, #befada)")}
        backgroundColor="orange.50"
        border="1px"
        borderColor="orange.100"
        marginBottom="4"
        w="70%"
        borderTopRightRadius="15px"
        borderTopLeftRadius="15px"
        borderBottomRightRadius={msg.inbound ? "15px" : "0px"}
        borderBottomLeftRadius={msg.inbound ? "0px" : "15px"}
        p="10px"
      >
        <Box w="full">
          <Text fontSize={14}>{msg.text}</Text>
          <Text
            align={msg.inbound ? "left" : "right"}
            noOfLines={1}
            fontSize={11}
          >
            {msg.managerName + " • " + msg.date.toLocaleString()}
          </Text>
        </Box>
      </Box>
    );
  };
  return (
    <Flex w="full" h="full" direction="column">
      <Flex overflow="hidden" direction="row" p="0" m="0" w="full" h="full">
        <Flex
          direction="column-reverse"
          p="4"
          m="0"
          w="full"
          overflow="hidden"
          paddingLeft="15px"
          paddingRight="15px"
          _hover={{
            overflowY: "scroll",
            paddingRight: "12px",
          }}
        >
          {messages.map((el, k) => (
            <ChatMessage key={k} msg={el}></ChatMessage>
          ))}
        </Flex>
      </Flex>
      <ChatMessageInput></ChatMessageInput>
    </Flex>
  );
}
export default Chat;

// const testMessage: Message[] = [
//   {
//     From: "Boris",
//     Text: "Hi we are so pleased to have you in our company",
//     Date: "2022-07-17 20:15",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes, thank you for invitation",
//     Date: "2022-07-17 21:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Would you like to book an appointment?",
//     Date: "2022-07-18 10:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes please",
//     Date: "2022-07-18 20:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Cool, you are booked to 3PM on July 20th",
//     Date: "2022-07-18 20:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Hello sir, Just a reminder for your meeting",
//     Date: "2022-07-20 10:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Hi we are so pleased to have you in our company",
//     Date: "2022-07-17 20:15",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes, thank you for invitation",
//     Date: "2022-07-17 21:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Would you like to book an appointment?",
//     Date: "2022-07-18 10:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes please",
//     Date: "2022-07-18 20:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Cool, you are booked to 3PM on July 20th",
//     Date: "2022-07-18 20:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Hello sir, Just a reminder for your meeting",
//     Date: "2022-07-20 10:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Hi we are so pleased to have you in our company",
//     Date: "2022-07-17 20:15",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes, thank you for invitation",
//     Date: "2022-07-17 21:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Would you like to book an appointment?",
//     Date: "2022-07-18 10:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Client",
//     Text: "Yes please",
//     Date: "2022-07-18 20:15",
//     Inbound: true,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Cool, you are booked to 3PM on July 20th",
//     Date: "2022-07-18 20:32",
//     Inbound: false,
//     Url: "",
//   },
//   {
//     From: "Boris",
//     Text: "Hello sir, Just a reminder for your meeting",
//     Date: "2022-07-20 10:32",
//     Inbound: false,
//     Url: "",
//   },
// ];
