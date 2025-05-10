import { IconCheck } from "@tabler/icons";
import { Box, Text } from "@mantine/core";
import { Image } from "@mantine/core";
import weightLoss from "../../utils/mock/weightLoss";

interface Props {
  imageSrc?: string;
  text: string;
  largeLetter: number;
  StepOptions: any;

  onClick: (text: string) => void;

  isSelected: boolean;
}

const StepItem = (props: Props) => {
  const { StepOptions, imageSrc, text, largeLetter, onClick, isSelected } =
    props;

  const handleColor = () => {
    onClick(text);
  };

  return (
    <Box
      sx={{
        border: isSelected ? "3px solid #00D1FA" : "3px solid lightgrey",
        whiteSpace: StepOptions !== weightLoss ? "nowrap" : undefined,

        "&:focus, &:focus-within": {
          outline: "black dashed 2px",
        },

        cursor: "pointer",
        borderRadius: "7px",
        padding: "25px 30px",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        position: "relative",
        overflow: "hidden",
      }}
      onClick={handleColor}
    >
      <IconCheck
        style={{
          color: "white",
          position: "absolute",
          left: "0",
          top: "0",
          zIndex: "9999",
        }}
      />
      <Image
        src={imageSrc}
        alt=""
        width={64}
        height={64}
        sx={{
          display: imageSrc ? "flex" : "none",
        }}
      />
      <Text
        sx={{
          fontSize:
            StepOptions.options[largeLetter].imageSrc == imageSrc
              ? "1.1rem"
              : "0.75rem",
          textAlign: "center",
          lineHeight: "1.5rem",
          fontFamily: '"Space Grotesk", sans-serif',
          fontWeight: "bold",
          textTransform: "uppercase",

          "@media(max-width: 540px)": {
            fontSize:
              StepOptions.options[largeLetter].imageSrc == imageSrc
                ? "0.75rem"
                : "0.6rem",
          },
        }}
      >
        {text}
      </Text>
      <Box
        style={{
          background: isSelected ? "#00D1FA" : "lightgrey",
          width: "6rem",
          height: "2rem",
          position: "absolute",
          left: "-40px",
          top: "0",
          transform: "rotate(140deg)",
          zIndex: "-9999",
        }}
      ></Box>
    </Box>
  );
};

export default StepItem;
