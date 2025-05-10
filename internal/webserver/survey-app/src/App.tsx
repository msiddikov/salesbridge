import { useState } from "react";
import MainStep from "./components/MainStep/MainStep";
import { defaultForm } from "./utils/constans/step";
import { Box } from "@mantine/core";

function App() {
  const [form, setForm] = useState(defaultForm);

  const [currentStep, setCurrentStep] = useState(0);

  const handleCurrentStepInrease = () => {
    setCurrentStep((currentStep) => currentStep + 1);
  };

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
      }}
      className="App"
    >
      <MainStep
        form={form}
        setForm={setForm}
        currentStep={currentStep}
        setStep={setCurrentStep}
        onClick={handleCurrentStepInrease}
      />
    </Box>
  );
}

export default App;
