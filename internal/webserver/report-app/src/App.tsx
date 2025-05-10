import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import Report from "./components/report/report";
import "./App.css";
import { LockOpen } from "@mui/icons-material";
import LocSettings from "./components/locSettings/locSettings";
import { ToastProvider } from "use-toast-mui";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

function App() {
  const queryParams = new URLSearchParams(window.location.search);
  let oper = queryParams.get("oper") || "";

  let app;
  switch (oper) {
    case "report":
      app = <Report />;
      break;
    case "locSettings":
      app = <LocSettings />;
      break;
  }

  return (
    <ToastProvider>
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />
        {app}
      </ThemeProvider>
    </ToastProvider>
  );
}

export default App;
