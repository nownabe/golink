import {
  AlertProps,
  Snackbar as MuiSnackbar,
  Alert as MuiAlert,
} from "@mui/material";
import { forwardRef } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
  severity: "success" | "error";
  children: string;
};

const Alert = forwardRef<HTMLDivElement, AlertProps>(
  function Alert(props, ref) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
  }
);

export default function Snackbar({ open, onClose, severity, children }: Props) {
  return (
    <MuiSnackbar open={open} autoHideDuration={6000} onClose={onClose}>
      <Alert onClose={onClose} severity={severity} sx={{ width: "100%" }}>
        {children}
      </Alert>
    </MuiSnackbar>
  );
}
