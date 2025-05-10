import {
  Box,
  Button,
  CircularProgress,
  FormControl,
  InputAdornment,
  InputLabel,
  Modal,
  OutlinedInput,
  Popover,
  Typography,
} from "@mui/material";
import { useState } from "react";
import { Location } from "../../types/types";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { DateRangePicker } from "react-date-range";
import { fetchServer } from "../../server/server";
const moment = require("moment");

export const ExpensesAdder = ({ locations }: { locations: Location[] }) => {
  const style = {
    position: "absolute" as "absolute",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 400,
    bgcolor: "background.paper",
    border: "2px solid #000",
    borderRadius: "15px",
    boxShadow: 24,
    p: 4,
  };
  const [amount, setAmount] = useState("");
  const [spinning, setSpinning] = useState(false);
  const [openM, setOpenM] = useState(false);
  const handleOpenM = () => setOpenM(true);
  const handleCloseM = () => setOpenM(false);
  const [range, setRange] = useState([
    {
      startDate: moment().startOf("month").toDate(),
      endDate: moment().endOf("month").toDate(),
      key: "selection",
    },
  ]);
  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? "simple-popover" : undefined;

  const sendData = () => {
    setSpinning(true);
    fetchServer("/reports/setExpense", {
      method: "POST",
      body: JSON.stringify({
        From: range[0].startDate,
        To: range[0].endDate,
        Locations: locations.map((v) => v.id),
        Total: +amount,
      }),
    })
      .then((res) => {
        if (res.status !== 200) {
          //alert("Something went wrong while fetching");
        }
        return res.json();
      })
      .finally(() => {
        setSpinning(false);
      });
  };

  return (
    <div>
      <Button onClick={handleOpenM}>+ Add expenses</Button>
      <Modal
        open={openM}
        onClose={handleCloseM}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={style}>
          <Typography id="modal-modal-title" variant="h6" component="h2">
            Add expenses for selected locations
          </Typography>
          <Typography
            variant="body1"
            sx={{
              cursor: "pointer",
              opacity: ".5",
              textAlign: "center",
              marginBottom: "20px",
            }}
            onClick={handleClick}
          >
            {range[0].startDate.toLocaleDateString() +
              " - " +
              range[0].endDate.toLocaleDateString()}
            <KeyboardArrowDownIcon
              sx={{
                marginBottom: "-5px",
              }}
            />
          </Typography>
          <Popover
            id={id}
            open={open}
            anchorEl={anchorEl}
            onClose={handleClose}
            anchorOrigin={{
              vertical: "bottom",
              horizontal: "left",
            }}
          >
            <DateRangePicker
              onChange={(item: any) => setRange([item.selection])}
              moveRangeOnFirstSelection={false}
              months={2}
              ranges={range}
              direction="vertical"
            />
          </Popover>
          <FormControl fullWidth sx={{ m: 1 }}>
            <InputLabel htmlFor="outlined-adornment-amount">Amount</InputLabel>
            <OutlinedInput
              id="outlined-adornment-amount"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              startAdornment={
                <InputAdornment position="start">$</InputAdornment>
              }
              label="Amount"
            />
          </FormControl>
          <Button onClick={sendData} disabled={spinning}>
            {" "}
            {spinning ? <CircularProgress /> : "Add expeses"}
          </Button>
          <Typography id="modal-modal-description" sx={{ mt: 2 }}>
            This will equally distribute entered amount between these locations:
            {locations.map((v) => v.name + ", ")}
          </Typography>
        </Box>
      </Modal>
    </div>
  );
};
