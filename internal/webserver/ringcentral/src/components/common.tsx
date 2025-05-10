import { Center } from "@chakra-ui/react";

export const RoundAvatar = ({
  name,
  width,
}: {
  name: string;
  width: number;
}) => {
  const stringToColor = (str: string) => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    let colour = "#";
    for (let i = 0; i < 3; i++) {
      let value = (hash >> (i * 8)) & 0xff;
      colour += ("00" + value.toString(16)).substr(-2);
    }
    return colour;
  };

  const words = name.split(" ", 2);
  let avatarName = "";

  switch (words.length) {
    case 1:
      avatarName += words[0][0] + words[0][words[0].length - 1];
      break;
    case 2:
      avatarName += words[0][0] + words[1][0];
      break;
  }
  avatarName = avatarName.toUpperCase();

  return (
    <Center
      w={"" + width + "px"}
      h={"" + width + "px"}
      borderRadius={"" + width / 2 + "px"}
      backgroundColor={stringToColor(name)}
      fontSize={"" + width / 2}
      textColor="white"
    >
      {avatarName}
    </Center>
  );
};
