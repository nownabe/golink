import { Code, ConnectError } from "@bufbuild/connect";
import { Alert, AlertTitle } from "@mui/material";
import { useRouteError } from "react-router-dom";

export default function ErrorDialog() {
  const e = useRouteError() as Error;
  const error = ConnectError.from(e);

  return (
    <Alert severity="error">
      <AlertTitle>Error: {Code[error.code]}</AlertTitle>
      {error.message}
    </Alert>
  );
}
