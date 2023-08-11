import { Code, ConnectError } from "@bufbuild/connect";
import { Alert, AlertTitle } from "@mui/material";

type Props = {
  error: ConnectError;
};

export default function Error({ error }: Props) {
  return (
    <Alert severity="error">
      <AlertTitle>{Code[error.code]}</AlertTitle>
      {error.message}
    </Alert>
  );
}
