import { Title } from "@mantine/core";
import hairRemoval from "../../utils/mock/hairRemoval";
import weightLoss from "../../utils/mock/weightLoss";

type Props = {
  titleNumber: any;
  StepOptions: typeof hairRemoval | typeof weightLoss;
};

const MainStepTitle = (props: Props) => {
  const { titleNumber, StepOptions } = props;

  return (
    <Title
      sx={{
        textAlign: "center",
        fontSize: "1.8rem",
        fontWeight: "bold",
        fontFamily: "Arial, Helvetica, sans-serif",
        width: "80%",

        height: "5.5rem",

        display: "grid",
        placeItems: "center",

        "@media(max-width: 720px)": {
          fontSize: "1.5rem",
        },

        "@media(max-width: 540px)": {
          fontSize: "1.1rem",
        },
      }}
      order={2}
      className="first__step-title"
    >
      {StepOptions.optionsTitle[titleNumber].title}
    </Title>
  );
};

export default MainStepTitle;
