import { Box, TextInput, Button } from "@mantine/core";
import { useForm, zodResolver } from "@mantine/form";
import { IStepForm } from "../../utils/models/step";
import zod from "zod";

interface Props {
  initialState?: IStepForm["form"];
  onSubmit: (value: IStepForm["form"]) => void;
}

const scheme = zod.object({
  name: zod.string().min(2, ""),
  phone: zod.string().min(9, ""),
  email: zod.string().email(""),
});

const defaultState: IStepForm["form"] = {
  name: "",
  phone: "",
  email: "",
};

const StepForm = (props: Props) => {
  const { initialState = defaultState, onSubmit } = props;

  const form = useForm({
    initialValues: initialState,
    clearInputErrorOnChange: true,
    validate: zodResolver(scheme),
  });

  const handleSubmit = (values: IStepForm["form"]) => {
    onSubmit(values);
  };

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        width: "517px",
        gap: "20px",
        "@media(max-width:720px)": {
          width: "400px",
        },
        "@media(max-width:540px)": {
          width: "300px",
        },
      }}
      component="form"
      onSubmit={form.onSubmit(handleSubmit)}
    >
      <TextInput
        styles={{ input: { height: "64px", color: "black" } }}
        sx={{ width: "100%" }}
        placeholder="Name"
        {...form.getInputProps("name")}
      />
      <TextInput
        styles={{ input: { height: "64px", color: "black" } }}
        sx={{ width: "100%" }}
        placeholder="Phone"
        {...form.getInputProps("phone")}
      />
      <TextInput
        styles={{ input: { height: "64px", color: "black" } }}
        sx={{ width: "100%" }}
        placeholder="Email"
        {...form.getInputProps("email")}
      />

      <Button
        sx={{
          height: "unset",
          background: "#00D1FA",
          border: "none",
          padding: "10px 80px",
          color: "white",
          fontSize: "2rem",
          textTransform: "uppercase",
          marginTop: "30px",
          cursor: "pointer",
          "&:disabled": {
            filter: "grayscale(0.5)",
          },
          "@media(max-width:720px)": {
            fontSize: "1.5rem",
            letterSpacing: "1px",
          },
          "@media(max-width:540px)": {
            fontSize: "1.1rem",
            letterSpacing: "none",
          },
        }}
        type="submit"
      >
        Submit
      </Button>
    </Box>
  );
};

export default StepForm;
