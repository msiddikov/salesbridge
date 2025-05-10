import { Box, Text, Container, Flex, Spinner } from "@chakra-ui/react";
import { EmailIcon, PhoneIcon } from "@chakra-ui/icons";
import {
  ContactInfoProps,
  ContactInfoRes,
  EmptyContactInfo,
} from "../types/types";
import { host } from "../consts";
import { useEffect, useState } from "react";
import { RoundAvatar } from "../common";

// const testInfo = {
//   Name: "John Borrow",
//   Phone: "+998913588881",
//   Email: "m.siddikov@gmail.com",
//   PhotoUrl:
//     "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Google_Contacts_icon.svg/128px-Google_Contacts_icon.svg.png",
// };

function ContactInfo({ chatInfo }: ContactInfoProps) {
  const [info, setInfo] = useState<ContactInfoRes>(EmptyContactInfo);
  const [spinning, setSpinning] = useState(false);
  useEffect(
    () => {
      if (!chatInfo.contactId) {
        return;
      }
      setSpinning(true);
      fetch(
        host + "/rc/contacts/" + chatInfo.contactId + "/" + chatInfo.locationId
      )
        .then((res) => res.json())
        .then((res: ContactInfoRes) => {
          setInfo(res);
          setSpinning(false);
        });
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [chatInfo]
  );

  return (
    <Flex alignItems="stretch" direction="column" h="full">
      <Container h="25%"></Container>
      {spinning ? (
        <Spinner></Spinner>
      ) : (
        <Container
          h="full"
          centerContent
          display={chatInfo.contactId ? undefined : "none"}
        >
          <RoundAvatar
            name={info.firstName + " " + info.lastName}
            width={100}
          ></RoundAvatar>

          <Text marginTop={5} fontSize="18">
            <b>{info.firstName + " " + info.lastName}</b>
          </Text>

          <Box bg="#34c0eb" w="50%" p="1px"></Box>

          <Container w="100%" marginTop={5}>
            <Text>
              <PhoneIcon /> {info.phone}
            </Text>
            <Text marginTop={3}>
              <EmailIcon /> {info.email}
            </Text>
          </Container>
        </Container>
      )}
      <Container h="150px" centerContent>
        {/* <Box
          as="button"
          borderRadius="md"
          bg="#34c0eb"
          color="white"
          px={4}
          h={8}
        >
          Book Appointment
        </Box> */}
      </Container>
    </Flex>
  );
}

export default ContactInfo;
