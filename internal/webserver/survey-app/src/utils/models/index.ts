import { Dispatch, SetStateAction } from "react";

export type SetState<State> = Dispatch<SetStateAction<State>>;
