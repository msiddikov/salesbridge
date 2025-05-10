import stomach from "../../assets/images/weightLoss/options/stomachreg.png";
import loveHandles from "../../assets/images/weightLoss/options/lovehandles.png";
import underNeck from "../../assets/images/weightLoss/options/underneck.png";
import back from "../../assets/images/weightLoss/options/back.png";
import innerThigh from "../../assets/images/weightLoss/options/innerthigh.png";
import outterThigh from "../../assets/images/weightLoss/options/outterthigh.png";
import butt from "../../assets/images/weightLoss/options/butt.png";
import upperArm from "../../assets/images/weightLoss/options/upperarms.png";
import other from "../../assets/images/weightLoss/options/other.png";

const weightLoss = {
  options: [
    [
      {
        url: stomach,
        text: "stomach",
        id: 1,
      },
      {
        url: loveHandles,
        text: "love handles",
        id: 2,
      },
      {
        url: underNeck,
        text: "under neck",
        id: 3,
      },
      {
        url: back,
        text: "back",
        id: 4,
      },
      {
        url: innerThigh,
        text: "inner thigh",
        id: 5,
      },
      {
        url: outterThigh,
        text: "outer thigh",
        id: 6,
      },
      {
        url: butt,
        text: "butt",
        id: 7,
      },
      {
        url: upperArm,
        text: "upper arm",
        id: 8,
      },
      {
        url: other,
        text: "other",
        id: 9,
      },
    ],
    [
      {
        text: "yes, both",
        id: 10,
      },
      {
        text: "neither",
        id: 11,
      },
      {
        text: "exersize only",
        id: 12,
      },
      {
        text: "healthy diet",
        id: 13,
      },
      {
        text: "not telling",
        id: 14,
      },
    ],
    [
      {
        text: "already there",
        id: 15,
      },
      {
        text: "almost there",
        id: 16,
      },
      {
        text: "less than 20 lbs",
        id: 17,
      },
      {
        text: "more than 20 lbs",
        id: 18,
      },
      {
        text: "not telling",
        id: 19,
      },
    ],
  ] as { text: string; url?: string; id: number }[][],
  optionsTitle: [
      {
        id: 1,
        title: "What area would you like to treat with body contouring?",
      },
      {
        id: 2,
        title: "Do you eat a healthy diet & exercise?",
      },
      {
        id: 3,
        title: "How close are you to your goal weight?",
      },

      {
        id: 4,
        title: "A Little More About You",
      },
      {
        id: 5,
        title: "",
      },
  ],
};

export default weightLoss;
