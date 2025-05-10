import { Box } from "@mantine/core";
import hairRemoval from "../../utils/mock/hairRemoval";
import weightLoss from "../../utils/mock/weightLoss";

type Props = {
  stepOf: number;
  StepOptions: typeof hairRemoval | typeof weightLoss;
};

const StepDecoration = (props: Props) => {
  const { stepOf, StepOptions } = props;

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
      }}
    >
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          gridTemplateRows: "repeat(4,1fr)",
          gap: "10px",
        }}
      >
        <Box
          sx={{
            background: stepOf >= 0 ? "darkgrey" : "lightgrey",
            width: "113px",
            height: "8px",
            borderRadius: "8px",

            "@media(max-width: 720px)": {
              width: "80px",
            },

            "@media(max-width: 540px)": {
              width: "50px",
            },
          }}
        ></Box>
        <Box
          sx={{
            background: stepOf >= 1 ? "darkgrey" : "lightgrey",
            width: "113px",
            height: "8px",
            borderRadius: "8px",
            "@media(max-width: 720px)": {
              width: "80px",
            },

            "@media(max-width: 540px)": {
              width: "50px",
            },
          }}
        ></Box>
        <Box
          sx={{
            background: stepOf >= 2 ? "darkgrey" : "lightgrey",
            width: "113px",
            height: "8px",
            borderRadius: "8px",
            "@media(max-width: 720px)": {
              width: "80px",
            },

            "@media(max-width: 540px)": {
              width: "50px",
            },
          }}
        ></Box>
        {StepOptions.options.length == 3 ? (
          <Box
            sx={{
              background: stepOf >= 3 ? "darkgrey" : "lightgrey",
              width: "113px",
              height: "8px",
              borderRadius: "8px",
              "@media(max-width: 720px)": {
                width: "80px",
              },

              "@media(max-width: 540px)": {
                width: "50px",
              },
            }}
          ></Box>
        ) : (
          ""
        )}
      </Box>
      <Box>
        <h2>
          Step {stepOf + 1} of {StepOptions.options.length + 1}
        </h2>
      </Box>
    </Box>
  );
};

export default StepDecoration;
