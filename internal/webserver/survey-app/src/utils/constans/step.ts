import { IStepForm, IForm } from '../models/step';

export const defaultStepForm: IStepForm = {
  1: "",
  2: "",
  3: "",
  form: {
    email: "",
    phone: "",
    name: "",
  },
};
export const defaultForm: IForm = {
  name: "",
  phone: "",
  email: "",
  answers:[]
};

export const host = "https://hayden.lavina.tech"
