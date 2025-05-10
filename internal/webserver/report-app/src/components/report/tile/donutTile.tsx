import { Doughnut } from "react-chartjs-2";
import { Chart as ChartJS, Tooltip, Legend, ArcElement } from "chart.js";
import { getColor, getColorByNumber, getRandomColorName } from "./colors";

export type DonutTileProps = {
  data: DonutTileData;
};

export type DonutTileData = {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
  }[];
};

ChartJS.register(ArcElement, Tooltip, Legend);

export function DonutTile(props: DonutTileProps) {
  const options = {
    responsive: true,
    tension: 0.5,
    backgroundColor: "rgba(0, 0, 0, 0.5)",
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: false,
        text: "Chart.js Line Chart",
      },
    },
  };

  let bgColors = [] as string[];
  let bdColors = [] as string[];
  for (let i = 0; i < props.data.labels.length; i++) {
    bgColors.push(getColorByNumber(i, 0.5));
    bdColors.push(getColorByNumber(i));
  }

  const data = {
    labels: props.data.labels,
    datasets: props.data.datasets.map((dataset) => {
      return {
        label: dataset.label,
        data: dataset.data,
        borderColor: bdColors,
        backgroundColor: bgColors,
      };
    }),
  };

  return <Doughnut options={options} data={data} />;
}

export const donutTileDataMock: DonutTileData = {
  labels: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
  datasets: [
    {
      label: "Dataset 1",
      data: [12, 19, 3, 5, 2, 3, 9],
    },
    {
      label: "Dataset 2",
      data: [1, 5, 3, 2, 7, 6, 11],
    },
  ],
};
