import { Box, Title } from "@mantine/core";
import { IconCircleCheck, IconCircleMinus } from "@tabler/icons";

interface Props {
  success: boolean;
}

const StepStatus = (props: Props) => {
  const { success } = props;

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        gap: "20px",
      }}
    >
      {success ? (
        <IconCircleCheck size="100px" color="green" />
      ) : (
        <IconCircleMinus size="100px" color="red" />
      )}
      <Title>
        {success ? "We have received your response" : "Something went wrong"}
      </Title>
    </Box>
  );
};

export default StepStatus;
