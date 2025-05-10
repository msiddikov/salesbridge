export {};

// bright colors for charts
export const colors: { [key: string]: string } = {
  blue: "53, 162, 235",
  red: "255, 99, 132",
  green: "75, 192, 192",
  orange: "255, 159, 64",
  purple: "153, 102, 255",
  grey: "201, 203, 207",
  yellow: "255, 205, 86",
  pink: "255, 102, 255",
  cyan: "102, 255, 255",
  lime: "102, 255, 102",
};

export function getColor(key: string, alpha: number = 1): string {
  return `rgba(${colors[key]}, ${alpha})`;
}

export function getColorByNumber(num: number, alpha: number = 1): string {
  const keys = Object.keys(colors);
  const key = keys[num % keys.length];
  return getColor(key, alpha);
}

export function getRandomColor(alpha: number = 1): string {
  const keys = Object.keys(colors);
  const key = keys[Math.floor(Math.random() * keys.length)];
  return getColor(key, alpha);
}

export function getRandomColorName(): string {
  const keys = Object.keys(colors);
  const key = keys[Math.floor(Math.random() * keys.length)];
  return key;
}
