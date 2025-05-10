import { Center, Flex, ChakraProvider } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import Chat from "../chat/chat";
import ContactInfo from "../contactInfo/contactInfo";
import { ChatInfo } from "../types/types";
import ContactList from "../contactist/contacts";
import "./app.css";

function ChatApp() {
  const [currentChat, setCurrentChat] = useState<ChatInfo>({} as ChatInfo);
  const [managerName, setManagerName] = useState("");

  useEffect(() => {
    const queryParams = new URLSearchParams(window.location.search);
    let id = queryParams.get("id") || "";
    let newChatInfo = { ...currentChat, locationId: id };
    if (newChatInfo !== currentChat) {
      setCurrentChat(newChatInfo);
    }
    let name = queryParams.get("name") || "";
    setManagerName(name);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    // @ts-ignore
    <ChakraProvider>
      <Flex alignItems="stretch">
        <Center w="600px" h="100vh" p="0">
          <ContactList
            chatInfo={currentChat}
            setChat={setCurrentChat}
          ></ContactList>
        </Center>
        <Center
          h="100vh"
          w="full"
          borderRight="1px"
          borderRightColor="gray.100"
          borderLeft="1px"
          borderLeftColor="gray.100"
        >
          <Chat
            chatInfo={currentChat}
            managerName={managerName}
            setChatInfo={setCurrentChat}
          ></Chat>
        </Center>
        <Center h="100vh" w="500px" overflow="hidden">
          <ContactInfo chatInfo={currentChat}></ContactInfo>
        </Center>
      </Flex>
    </ChakraProvider>
  );
}

export default ChatApp;
