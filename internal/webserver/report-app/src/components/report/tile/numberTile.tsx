import { Container } from "@mui/material";
import CountUp from "react-countup";

export {};

export type TileProps = {
  data: number;
  pre: string;
  after: string;
};

export function NumberTile(props: TileProps) {
  return (
    <Container>
      {props.pre}
      <CountUp duration={0.7} end={props.data}></CountUp>
      {props.after}
    </Container>
  );
}
