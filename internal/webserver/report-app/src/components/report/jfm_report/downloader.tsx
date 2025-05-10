import { host } from "../../server/server";
import { ReportProps } from "../../types/types";
import { Button } from "@mui/material";

export function JfmReportDownloader(props: ReportProps) {
  const handleDownload = () => {
    fetch(host + "/reports/getDetails", {
      method: "POST",
      body: JSON.stringify({
        From: props.start,
        To: props.end,
        Locations: props.locations.map((v) => v.id),
        Tags: props.tags,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => response.blob())
      .then((blob) => {
        // 2. Create blob link to download
        const url = window.URL.createObjectURL(new Blob([blob]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", `details.xlsx`);
        // 3. Append to html page
        document.body.appendChild(link);
        // 4. Force download
        link.click();
        // 5. Clean up and remove the link
        if (link.parentNode) {
          link.parentNode.removeChild(link);
        }
      });
  };
  return <Button onClick={handleDownload}>Download details</Button>;
}
