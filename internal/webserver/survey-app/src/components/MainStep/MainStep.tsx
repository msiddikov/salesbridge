import { useState, useEffect } from "react";
import { Box, Button, Loader } from "@mantine/core";
import StepItem from "../Step/StepItem";
import StepForm from "../Step/StepForm";
import { IStepForm, IForm } from "../../utils/models/step";
import { SetState } from "../../utils/models";
import StepDecoration from "../StepDecoration/StepDecoration";
import { host, defaultForm } from "../../utils/constans/step";
import StepStatus from "../Step/StepStatus";
import MainStepTitle from "./MainStepTitle";
import weightLoss from "../../utils/mock/weightLoss";
import hairRemoval from "../../utils/mock/hairRemoval";

const options = {
  weightLoss: weightLoss,
  hairRemoval: hairRemoval,
};

type Props = {
  currentStep: number;
  form: IForm;
  setForm: SetState<IForm>;
  setStep: SetState<number>;
  onClick: () => void;
};

const MainStep = (props: Props) => {
  const { currentStep, form, setForm, setStep, onClick } = props;
  const [nextNumber, setNextNumber] = useState(0);
  const [spinning, setSpinning] = useState(false);
  const [success, setSuccess] = useState(false);
  const urlParams = new URLSearchParams(window.location.search);
  const surveyId = urlParams.get("surveyId");
  const locationId = urlParams.get("locationId");
  const workflowId = urlParams.get("workflowId");
  const url = urlParams.get("url");

  const StepOptions = options[surveyId as keyof typeof options];

  const handleSubmit = (values: IStepForm["form"]) => {
    const payload = { ...values, answers: form.answers };

    console.log(payload);

    setSpinning(true);
    fetch(host + "/survey/" + locationId + "/" + workflowId, {
      method: "post",
      body: JSON.stringify(payload),
    })
      .then((res: Response) => {
        if (res.ok) {
          setSuccess(true);
          if (url) {
            window.open(url);
          }
        }
      })
      .finally(() => {
        setSpinning(false);
        setStep(currentStep + 1);
      });

    setForm(defaultForm);
  };

  const handleNext = (text: string) => {
    setForm({
      ...form,
      answers: [
        ...form.answers.filter(
          (v) => v.question !== StepOptions.optionsTitle[nextNumber].title
        ),
        {
          question: StepOptions.optionsTitle[nextNumber].title,
          answer: text,
        },
      ],
    });
  };

  useEffect(() => {
    setNextNumber(currentStep);
  }, [currentStep]);

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
        "@media(min-width:720px)": {
          width: "650px",
        },
      }}
    >
      <MainStepTitle titleNumber={nextNumber} StepOptions={StepOptions} />
      {currentStep > StepOptions.options.length ? null : (
        <StepDecoration StepOptions={StepOptions} stepOf={currentStep} />
      )}
      <Box
        sx={{
          gridTemplateRows:
            nextNumber < StepOptions.options.length
              ? `repeat(${Math.ceil(
                  StepOptions.options[nextNumber].length / 3
                )}, 148px)`
              : undefined,
          display: nextNumber < StepOptions.options.length ? "grid" : "flex",
          gridTemplateColumns: "repeat(3, 1fr)",
          gap: "min(10px, 2%)",

          width: "30rem",
          margin: "auto",
          justifyContent: "center",

          "@media (max-width: 30rem)": {
            maxWidth: "30rem",
            width: "100%",
          },
        }}
      >
        {spinning ? (
          <Loader size={100} />
        ) : nextNumber < StepOptions.options.length ? (
          StepOptions.options[nextNumber].map((item) => (
            <StepItem
              key={item.id}
              largeLetter={nextNumber}
              text={item.text}
              imageSrc={item.url}
              onClick={handleNext}
              StepOptions={StepOptions}
              isSelected={
                form.answers.filter(
                  (v) =>
                    v.question === StepOptions.optionsTitle[nextNumber].title
                )[0]?.answer === item.text
              }
            />
          ))
        ) : nextNumber === StepOptions.options.length ? (
          <StepForm onSubmit={handleSubmit} />
        ) : (
          <StepStatus success={success} />
        )}
      </Box>
      <Button
        onClick={onClick}
        className="app__button"
        sx={{
          background: "#00D1FA",
          border: "none",
          padding: "10px 80px",
          color: "white",
          height: "60px",
          fontSize: "2rem",
          textTransform: "uppercase",
          marginTop: "30px",
          cursor: "pointer",

          "&:disabled": {
            filter: "grayscale(0.5)",
          },
          display:
            currentStep >= StepOptions.optionsTitle.length - 2
              ? "none"
              : undefined,
          "@media(max-width:720px)": {
            fontSize: "1.5rem",
            height: "50px",
          },
          "@media(max-width:540px)": {
            fontSize: "1.1rem",
            height: "40px",
          },
        }}
        disabled={
          !form.answers.filter(
            (v) => v.question === StepOptions.optionsTitle[currentStep].title
          )[0]?.answer
        }
      >
        next
      </Button>
    </Box>
  );
};

export default MainStep;
