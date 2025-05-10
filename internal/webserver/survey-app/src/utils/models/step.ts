export interface IStepForm {
  1: string;
  2: string;
  3: string;
  form: {
    name: string;
    phone: string;
    email: string;
  };
}

export interface IForm {
  name: string;
  phone: string;
  email: string;
  answers: {
    question: string;
    answer: string;
  }[];
}

export type TOptionStep = {
  url?: string;
  text?: string;
  id: number;
};

export type TSteps = {
  weightLoss: TOptionStep[][];
  hairRemoval: TOptionStep[][];
};
