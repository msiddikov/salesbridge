import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { getColor, getRandomColor, getRandomColorName } from "./colors";

export type LineTileProps = {
  data: LineTileData;
};

export type LineTileData = {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
  }[];
};

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

export function LineTile(props: LineTileProps) {
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

  const data = {
    labels: props.data.labels,
    datasets: props.data.datasets.map((dataset) => {
      let color = getRandomColorName();
      return {
        label: dataset.label,
        data: dataset.data,
        borderColor: getColor(color),
        backgroundColor: getColor(color, 0.5),
      };
    }),
  };

  return <Line options={options} data={data} />;
}

export const lineTileDataMock: LineTileData = {
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
